package mempool

import (
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cometbft/cometbft/crypto"
	bytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/cometbft/cometbft/p2p/nodekey"
	"github.com/cometbft/cometbft/types"
)

// mempoolTx is an entry in the mempool.
type mempoolTx struct {
	height    int64    // height that this tx had been validated in
	gasWanted int64    // amount of gas this tx states it will require
	tx        types.Tx // validated by the application
	lane      LaneID
	seq       int64
	timestamp time.Time // time when entry was created

	// signatures of peers who've sent us this tx (as a map for quick lookups).
	// Map keys are string representations of PubKey.
	signatures map[string][]byte
	signatureMutex sync.Mutex // Protects access to the signatures map

	// Number of valid signatures
	signatureCount int32 // Atomic counter for signatures
	// ids of peers who've sent us this tx (as a map for quick lookups).
	// senders: PeerID -> struct{}
	senders sync.Map

	// Is the threshold reached on this transaction
	isTresholdReached bool
}

func (memTx *mempoolTx) Tx() types.Tx {
	return memTx.tx
}

func (memTx *mempoolTx) Height() int64 {
	return atomic.LoadInt64(&memTx.height)
}

func (memTx *mempoolTx) GasWanted() int64 {
	return memTx.gasWanted
}

func (memTx *mempoolTx) IsSender(peerID nodekey.ID) bool {
	_, ok := memTx.senders.Load(peerID)
	return ok
}

func (memTx *mempoolTx) SetThresholdReached(reached bool)  {
	memTx.isTresholdReached = reached
}

func (memTx *mempoolTx) GetThresholdReached() bool {
	return memTx.isTresholdReached
}

func (memTx *mempoolTx) ValidateSignatures(validators *types.ValidatorSet) (int64, error) {
    var accumulatedVotingPower int64
    var invalidSignatures []string

    memTx.signatureMutex.Lock()
    defer memTx.signatureMutex.Unlock()

    // Log the transaction hash (and optionally some other tx details)
    // info-level logging may be passed down as a logger parameter if needed.
    // For example: logger.Info("Validating signatures", "txHash", fmt.Sprintf("%X", txHash))

    for pubKeyStr, signature := range memTx.signatures {
        // Convert the stored string back to a PubKey.
		pubKeyBytes, err := hex.DecodeString(strings.ToLower(pubKeyStr))
		if err != nil {
			//
		}
		pubKeyAdrr := bytes.HexBytes(pubKeyBytes)
        if err != nil {
            invalidSignatures = append(invalidSignatures, pubKeyStr)
            continue
        }

        // Look up the validator by the public key's address.
        _, validator := validators.GetByAddress(pubKeyAdrr)
		pubKey := validator.PubKey

        if validator == nil {
            // Using info-level log; ensure the output is consistent.
            print("MEMPOOLTX INFO: Invalid signature from unknown validator")
            invalidSignatures = append(invalidSignatures, pubKeyStr)
            continue
        }
        // Log that we’re about to verify.
        // print("MEMPOOLTX INFO: Verifying signature for validator " + pubKeyStr);

        // Verify the signature against the transaction.
        if !pubKey.VerifySignature(memTx.Tx().Hash(), signature) {
            invalidSignatures = append(invalidSignatures, pubKeyStr)
            continue
        }

        accumulatedVotingPower += validator.VotingPower
    }

    if len(invalidSignatures) > 0 {
        return accumulatedVotingPower, fmt.Errorf("invalid signatures found for public keys: %v", invalidSignatures)
    }

    return accumulatedVotingPower, nil
}

// GetSignatures returns the signatures map, initializing it if necessary.
func (memTx *mempoolTx) GetSignatures() map[string][]byte {
    memTx.signatureMutex.Lock()
    defer memTx.signatureMutex.Unlock()

    copy := make(map[string][]byte)
    for k, v := range memTx.signatures {
        copy[k] = v
    }
    return copy
}

// SetSignatures safely updates the signatures map, merging new signatures with existing ones.
func (memTx *mempoolTx) SetSignatures(signatures map[string][]byte) {
	memTx.signatureMutex.Lock()
	defer memTx.signatureMutex.Unlock()

	// If no signatures exist yet, initialize the map
	if memTx.signatures == nil {
		memTx.signatures = make(map[string][]byte)
	}

	// Merge new signatures into existing map
	for pubKey, signature := range signatures {
		memTx.signatures[pubKey] = signature
	}
}

// AddSignature safely adds a signature to the map and increments the signature count.
func (memTx *mempoolTx) AddSignature(pubKey crypto.PubKey, signature []byte) {

	memTx.signatureMutex.Lock()
	defer memTx.signatureMutex.Unlock()

	if memTx.signatures == nil {
		memTx.signatures = make(map[string][]byte)
	}

	memTx.signatures[pubKey.Address().String()] = signature
	atomic.AddInt32(&memTx.signatureCount, 1)
}

// SignatureCount returns the number of valid signatures.
func (memTx *mempoolTx) SignatureCount() int {
	return int(atomic.LoadInt32(&memTx.signatureCount))
}

// AddSender adds the peer ID to the list of senders. Returns true if it already existed.
func (memTx *mempoolTx) addSender(peerID nodekey.ID) bool {
	if len(peerID) == 0 {
		return false
	}
	if _, loaded := memTx.senders.LoadOrStore(peerID, struct{}{}); loaded {
		return true
	}
	return false
}
