package mempool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/crypto_implementation"
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

// Validates All Signatures check if each signature is matching with the pubKey
// If a invalidSignatures return an error;
func (memTx *mempoolTx) ValidateSignaturesAndGetVotingPower(validators *types.ValidatorSet) (int64, error) {
	var accumulatedVotingPower int64
	var invalidSignatures []string

	memTx.signatureMutex.Lock()
	defer memTx.signatureMutex.Unlock()

	// Iterate over each signature using the map iteration logic.
	for pubKeyStr, signature := range memTx.signatures {
		// Convert the stored string back to a PubKey.
		pubKey, err := crypto_implementation.PubKeyFromBytes([]byte(pubKeyStr))
		if err != nil {
			invalidSignatures = append(invalidSignatures, pubKeyStr)
			continue
		}

		// Look up the validator by the public key's address.
		_, validator := validators.GetByAddress(pubKey.Address())
		if validator == nil {
			fmt.Printf("Invalid signature from unknown validator: %v\n", pubKey.Address())
			invalidSignatures = append(invalidSignatures, pubKeyStr)
			continue
		}

		// Verify the signature against the transaction.
		if !pubKey.VerifySignature(memTx.tx, signature) {
			fmt.Printf("Invalid signature from validator: %v\n", validator.Address.String())
			invalidSignatures = append(invalidSignatures, pubKeyStr)
			continue
		}

		// Accumulate the voting power of the valid validator.
		accumulatedVotingPower += validator.VotingPower
	}

	if len(invalidSignatures) > 0 {
		return accumulatedVotingPower, fmt.Errorf("invalid signatures found for public keys: %v", invalidSignatures)
	}

	return accumulatedVotingPower, nil
}


// func (memTx *mempoolTx) ValidateSignatures() error {
// 	var invalidSignatures []string

// 	memTx.signatureMutex.Lock()
// 	defer memTx.signatureMutex.Unlock()

// 	for pubKeyStr, signature := range memTx.signatures {
// 		pubKey, err := crypto_implementation.PubKeyFromBytes([]byte(pubKeyStr))
// 		if err != nil {
// 			return fmt.Errorf("invalid publicKey from bytes: %s", pubKeyStr)
// 		}

// 		if !pubKey.VerifySignature(memTx.tx, signature) {
// 			invalidSignatures = append(invalidSignatures, pubKeyStr)
// 		}
// 	}

// 	if len(invalidSignatures) > 0 {
// 		return fmt.Errorf("invalid signatures found for public keys: %v", invalidSignatures)
// 	}

// 	return nil
// }

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

// SetSignatures safely sets the signatures map.
func (memTx *mempoolTx) SetSignatures(signatures map[string][]byte) {
	memTx.signatureMutex.Lock()
	defer memTx.signatureMutex.Unlock()

	memTx.signatures = signatures
}

// AddSignature safely adds a signature to the map and increments the signature count only if the key is new.
func (memTx *mempoolTx) AddSignature(pubKey crypto.PubKey, signature []byte) {
    pubKeyStr := string(pubKey.Bytes())

    memTx.signatureMutex.Lock()
    defer memTx.signatureMutex.Unlock()

    if memTx.signatures == nil {
        memTx.signatures = make(map[string][]byte)
    }

    // Only increment the count if the pubKey is not already present.
    if _, exists := memTx.signatures[pubKeyStr]; !exists {
        atomic.AddInt32(&memTx.signatureCount, 1)
    }
    memTx.signatures[pubKeyStr] = signature
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
