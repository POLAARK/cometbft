package mempool

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"

	abcicli "github.com/cometbft/cometbft/abci/client"
	protomem "github.com/cometbft/cometbft/api/cometbft/mempool/v1"
	cfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/crypto_implementation"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/p2p"
	tcpconn "github.com/cometbft/cometbft/p2p/transport/tcp/conn"
	"github.com/cometbft/cometbft/types"
)

// Reactor handles mempool tx broadcasting amongst peers.
// It maintains a map from peer ID to counter, to prevent gossiping txs to the
// peers you received it from.
type Reactor struct {
	p2p.BaseReactor
	config  *cfg.MempoolConfig
	mempool *CListMempool

	waitSync   atomic.Bool
	waitSyncCh chan struct{} // for signaling when to start receiving and sending txs

	// Semaphores to keep track of how many connections to peers are active for broadcasting
	// transactions. Each semaphore has a capacity that puts an upper bound on the number of
	// connections for different groups of peers.
	activePersistentPeersSemaphore    *semaphore.Weighted
	activeNonPersistentPeersSemaphore *semaphore.Weighted

	privVal *types.PrivValidator

	txBroadcastThreshold int32
	thresholdPercent int32
}


// NewReactor returns a new Reactor with the given config and mempool.
//TODO : check where the reactor is created, we will need the privValidator to sign transactions there.
func NewReactor(config *cfg.MempoolConfig, privVal *types.PrivValidator, mempool *CListMempool, waitSync bool, thresholdPercent int32) *Reactor {

	memR := &Reactor{
		config:   config,
		mempool:  mempool,
		privVal : privVal,
		waitSync: atomic.Bool{},
		txBroadcastThreshold: thresholdPercent,
		thresholdPercent: thresholdPercent,
	}
	memR.BaseReactor = *p2p.NewBaseReactor("Mempool", memR)
	if waitSync {
		memR.waitSync.Store(true)
		memR.waitSyncCh = make(chan struct{})
	}
	memR.activePersistentPeersSemaphore = semaphore.NewWeighted(int64(memR.config.ExperimentalMaxGossipConnectionsToPersistentPeers))
	memR.activeNonPersistentPeersSemaphore = semaphore.NewWeighted(int64(memR.config.ExperimentalMaxGossipConnectionsToNonPersistentPeers))

	return memR
}

func (memR *Reactor) calculateTxBroadcastThreshold() {
    peerCount := memR.Switch.Peers().Size() // Dynamically get the number of peers
    newThreshold := (memR.thresholdPercent * int32(peerCount)) / 100

    atomic.StoreInt32(&memR.txBroadcastThreshold, newThreshold) // Atomic store

    memR.Logger.Debug("Calculated txBroadcastThreshold",
        "thresholdPercent", memR.thresholdPercent,
        "peerCount", peerCount,
        "txBroadcastThreshold", newThreshold)
}

// SetLogger sets the Logger on the reactor and the underlying mempool.
func (memR *Reactor) SetLogger(l log.Logger) {
	memR.Logger = l
	memR.mempool.SetLogger(l)
}

// OnStart implements p2p.BaseReactor.
func (memR *Reactor) OnStart() error {
	if memR.WaitSync() {
		memR.Logger.Info("Starting reactor in sync mode: tx propagation will start once sync completes")
	}
	if !memR.config.Broadcast {
		memR.Logger.Info("Tx broadcasting is disabled")
	}
	return nil
}

// StreamDescriptors implements Reactor by returning the list of channels for this
// reactor.
func (memR *Reactor) StreamDescriptors() []p2p.StreamDescriptor {
	// Create a largestTx byte slice to simulate the maximum transaction size.
	largestTx := make([]byte, memR.config.MaxTxBytes)

	// TODOPB should be use an initiated tx map ?
	// Construct a Transaction with the largestTx and an empty signatures map.
	largestTransaction := &protomem.Transaction{
		TransactionBytes:         largestTx,
		Signatures: map[string][]byte{}, // Empty signatures map for the example.
	}

	// Create a batch message with the new structure of Txs containing Transactions.
	batchMsg := protomem.Message{
		Sum: &protomem.Message_Txs{
			Txs: &protomem.Txs{Txs: []*protomem.Transaction{largestTransaction}},
		},
	}

	return []p2p.StreamDescriptor{
		&tcpconn.ChannelDescriptor{
			ID:                  MempoolChannel,
			Priority:            5,
			RecvMessageCapacity: batchMsg.Size(),
			MessageTypeI:        &protomem.Message{},
		},
	}
}

// AddPeer implements Reactor.
// It starts a broadcast routine ensuring all txs are forwarded to the given peer.
func (memR *Reactor) AddPeer(peer p2p.Peer) {
	if memR.config.Broadcast && peer.HasChannel(MempoolChannel) {
		go func() {
			// Always forward transactions to unconditional peers.
			if !memR.Switch.IsPeerUnconditional(peer.ID()) {
				// Depending on the type of peer, we choose a semaphore to limit the gossiping peers.
				var peerSemaphore *semaphore.Weighted
				if peer.IsPersistent() && memR.config.ExperimentalMaxGossipConnectionsToPersistentPeers > 0 {
					peerSemaphore = memR.activePersistentPeersSemaphore
				} else if !peer.IsPersistent() && memR.config.ExperimentalMaxGossipConnectionsToNonPersistentPeers > 0 {
					peerSemaphore = memR.activeNonPersistentPeersSemaphore
				}

				if peerSemaphore != nil {
					for peer.IsRunning() {
						// Block on the semaphore until a slot is available to start gossiping with this peer.
						// Do not block indefinitely, in case the peer is disconnected before gossiping starts.
						ctxTimeout, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
						// Block sending transactions to peer until one of the connections become
						// available in the semaphore.
						err := peerSemaphore.Acquire(ctxTimeout, 1)
						cancel()

						if err != nil {
							continue
						}

						// Release semaphore to allow other peer to start sending transactions.
						defer peerSemaphore.Release(1)
						break
					}
				}
			}

			memR.mempool.metrics.ActiveOutboundConnections.Add(1)
			defer memR.mempool.metrics.ActiveOutboundConnections.Add(-1)
			memR.calculateTxBroadcastThreshold()
			memR.broadcastTxRoutine(peer)
		}()
	}
}

// Receive implements Reactor.
// It adds any received transactions to the mempool.
func (memR *Reactor) Receive(e p2p.Envelope) {
	memR.Logger.Debug("Receive", "src", e.Src, "chId", e.ChannelID, "msg", e.Message)

	switch msg := e.Message.(type) {
	case *protomem.Txs:
		if memR.WaitSync() {
			memR.Logger.Debug("Ignored message received while syncing", "msg", msg)
			return
		}

		protoTxs := msg.GetTxs()
		if len(protoTxs) == 0 {
			memR.Logger.Error("Received empty Txs message from peer", "src", e.Src)
			return
		}

		for _, protoTransaction := range protoTxs {
			tx := types.Tx(protoTransaction.TransactionBytes)
			_, err := memR.TryAddTx(tx, e.Src)
			if err != nil {
				memR.Logger.Error("Failed to add transaction", "err", err)
				continue
			}

			// ADD SIGNATURES AFTER !
			signatures := protoTransaction.Signatures
			convertedSignatures, err := crypto_implementation.ConvertSignatures(signatures)
			if err != nil {
				fmt.Printf("Error converting signatures: %v\n", err)
				return
			}
			err = memR.mempool.AddSignatures(tx.Key(), convertedSignatures)
			if err != nil {
				memR.Logger.Error("Failed to add signatures to transaction", "err", err)
			}
		}

	default:
		memR.Logger.Error("Unknown message type", "src", e.Src, "chId", e.ChannelID, "msg", e.Message)
		memR.Switch.StopPeerForError(e.Src, fmt.Errorf("mempool cannot handle message of type: %T", e.Message))
		return
	}

	// Broadcasting happens from go routines per peer
}

// TryTx attempts to add an incoming transaction to the mempool.
// When the sender is nil, it means the transaction comes from an RPC endpoint.
func (memR *Reactor) TryAddTx(tx types.Tx, sender p2p.Peer) (*abcicli.ReqRes, error) {
	senderID := noSender
	if sender != nil {
		senderID = sender.ID()
	}

	reqRes, err := memR.mempool.CheckTx(tx, senderID)
	if err != nil {
		switch {
		case errors.Is(err, ErrTxInCache):
			memR.Logger.Debug("Tx already exists in cache", "tx", log.NewLazySprintf("%X", tx.Hash()), "sender", senderID)
		case errors.As(err, &ErrMempoolIsFull{}):
			// using debug level to avoid flooding when traffic is high
			memR.Logger.Debug(err.Error())
		default:
			memR.Logger.Info("Could not check tx", "tx", log.NewLazySprintf("%X", tx.Hash()), "sender", senderID, "err", err)
		}
		return nil, err
	}

	return reqRes, nil
}

func (memR *Reactor) EnableInOutTxs() {
	memR.Logger.Info("Enabling inbound and outbound transactions")
	if !memR.waitSync.CompareAndSwap(true, false) {
		return
	}

	// Releases all the blocked broadcastTxRoutine instances.
	if memR.config.Broadcast {
		close(memR.waitSyncCh)
	}
}

func (memR *Reactor) WaitSync() bool {
	return memR.waitSync.Load()
}

// PeerState describes the state of a peer.
type PeerState interface {
	GetHeight() int64
}

// Send new mempool txs to peer.
func (memR *Reactor) broadcastTxRoutine(peer p2p.Peer) {
	// If the node is catching up, don't start this routine immediately.
	if memR.WaitSync() {
		select {
		case <-memR.waitSyncCh:
			// EnableInOutTxs() has set WaitSync() to false.
		case <-memR.Quit():
			return
		}
	}


	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-peer.Quit():
			cancel()
		case <-memR.Quit():
			cancel()
		}
	}()

	iter := NewBlockingIterator(ctx, memR.mempool, string(peer.ID()))
	for {
		// In case of both next.NextWaitChan() and peer.Quit() are variable at the same time
		if !memR.IsRunning() || !peer.IsRunning() {
			return
		}

		entry := <-iter.WaitNextCh()
		// If the entry we were looking at got garbage collected (removed), try again.
		if entry == nil {
			continue
		}

		// If we suspect that the peer is lagging behind, at least by more than
		// one block, we don't send the transaction immediately. This code
		// reduces the mempool size and the recheck-tx rate of the receiving
		// node. See [RFC 103] for an analysis on this optimization.
		//
		// [RFC 103]: https://github.com/CometBFT/cometbft/blob/main/docs/references/rfc/rfc-103-incoming-txs-when-catching-up.md
		for {
			// Make sure the peer's state is up to date. The peer may not have a
			// state yet. We set it in the consensus reactor, but when we add
			// peer in Switch, the order we call reactors#AddPeer is different
			// every time due to us using a map. Sometimes other reactors will
			// be initialized before the consensus reactor. We should wait a few
			// milliseconds and retry.
			peerState, ok := peer.Get(types.PeerStateKey).(PeerState)
			if ok && peerState.GetHeight()+1 >= entry.Height() {
				break
			}
			select {
			case <-time.After(PeerCatchupSleepIntervalMS * time.Millisecond):
			case <-peer.Quit():
				return
			case <-memR.Quit():
				return
			}
		}

		// TODOPB : is this enough ?
		// TODOPB : will the network sync or should we add a mechanism to push to consensus
		if entry.SignatureCount() >= int(memR.txBroadcastThreshold) {
            memR.Logger.Debug("Transaction reached threshold, stopping broadcast",
                "tx", log.NewLazySprintf("%X", entry.Tx().Hash()), "threshold", memR.txBroadcastThreshold)
            continue
        }


		// NOTE: Transaction batching was disabled due to
		// https://github.com/tendermint/tendermint/issues/5796

		// We are paying the cost of computing the transaction hash in
		// any case, even when logger level > debug. So it only once.
		// See: https://github.com/cometbft/cometbft/issues/4167
		txHash := entry.Tx().Hash()

		// Do not send this transaction if we receive it from peer.
		if entry.IsSender(peer.ID()) {
			memR.Logger.Debug("Skipping transaction, peer is sender",
				"tx", log.NewLazySprintf("%X", txHash), "peer", peer.ID())
			continue
		}

		for {
			memR.Logger.Debug("Sending transaction to peer",
			"tx", log.NewLazySprintf("%X", txHash), "peer", peer.ID())

			if err := memR.signAndValidate(entry.(*mempoolTx)); err != nil {
				memR.Logger.Error("Failed to sign and validate transaction", "error", err)
				continue
			}
			// The entry may have been removed from the mempool since it was
			// chosen at the beginning of the loop. Skip it if that's the case.
			if !memR.mempool.Contains(entry.Tx().Key()) {
				break
			}

			memR.Logger.Debug("Sending transaction to peer",
				"tx", log.NewLazySprintf("%X", txHash), "peer", peer.ID())

			signaturesMap := entry.Signatures()

			success := peer.Send(p2p.Envelope{
				ChannelID: MempoolChannel,
				Message: &protomem.Txs{
					Txs: []*protomem.Transaction{
						{
							TransactionBytes: entry.Tx(),
							Signatures:       signaturesMap,
						},
					},
				},
			})

			if success {
				break
			}

			memR.Logger.Debug("Failed sending transaction to peer",
				"tx", log.NewLazySprintf("%X", txHash), "peer", peer.ID())

			select {
			case <-time.After(PeerCatchupSleepIntervalMS * time.Millisecond):
			case <-peer.Quit():
				return
			case <-memR.Quit():
				return
			}
		}
	}
}

func (memR *Reactor) signAndValidate(entry Entry) error {
	tx := entry.Tx()

	// Validate existing signatures
	if err := entry.ValidateSignatures(); err != nil {
		return fmt.Errorf("signature validation failed: %w", err)
	}

	// Sign transaction
	signature, err := (*memR.privVal).SignBytes(tx) // TODOPB: SHOULD WE ONLY SIGN THE HASH OR THE WHOLE TX
	if err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}

	// Get public key
	pubKey, err := (*memR.privVal).GetPubKey()
	if err != nil {
		return fmt.Errorf("public key retrieval failed: %w", err)
	}

	// Add the new signature
	entry.AddSignature(pubKey, signature)
	return nil
}