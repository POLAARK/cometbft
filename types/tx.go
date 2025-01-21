package types

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"
)

// TxKeySize is the size of the transaction key index.
const TxKeySize = sha256.Size

type (
	// Tx is an arbitrary byte array.
	// NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
	// Might we want types here ?
	Tx []byte

	// TxKey is the fixed length array key used as an index.
	TxKey [TxKeySize]byte
)

// Hash computes the TMHASH hash of the wire encoded transaction.
func (tx Tx) Hash() cmtbytes.HexBytes {
	return tmhash.Sum(tx)
}

func (tx Tx) Key() TxKey {
	return sha256.Sum256(tx)
}

// String returns the hex-encoded transaction as a string.
func (tx Tx) String() string {
	return fmt.Sprintf("Tx{%X}", []byte(tx))
}

// Txs is a slice of Tx.
type Txs []Tx

// Hash returns the Merkle root hash of the transaction hashes.
// i.e. the leaves of the tree are the hashes of the txs.
func (txs Txs) Hash() []byte {
	hl := txs.hashList()
	return merkle.HashFromByteSlices(hl)
}

// Index returns the index of this transaction in the list, or -1 if not found.
func (txs Txs) Index(tx Tx) int {
	for i := range txs {
		if bytes.Equal(txs[i], tx) {
			return i
		}
	}
	return -1
}

// IndexByHash returns the index of this transaction hash in the list, or -1 if not found.
func (txs Txs) IndexByHash(hash []byte) int {
	for i := range txs {
		if bytes.Equal(txs[i].Hash(), hash) {
			return i
		}
	}
	return -1
}

func (txs Txs) Proof(i int) TxProof {
	hl := txs.hashList()
	root, proofs := merkle.ProofsFromByteSlices(hl)

	return TxProof{
		RootHash: root,
		Data:     txs[i],
		Proof:    *proofs[i],
	}
}

func (txs Txs) hashList() [][]byte {
	hl := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		hl[i] = txs[i].Hash()
	}
	return hl
}

// Txs is a slice of transactions. Sorting a Txs value orders the transactions
// lexicographically.

// Deprecated: Do not use.
func (txs Txs) Len() int { return len(txs) }

// Deprecated: Do not use.
func (txs Txs) Swap(i, j int) { txs[i], txs[j] = txs[j], txs[i] }

// Deprecated: Do not use.
func (txs Txs) Less(i, j int) bool {
	return bytes.Compare(txs[i], txs[j]) == -1
}

func ToTxs(txl [][]byte) Txs {
	txs := make([]Tx, 0, len(txl))
	for _, tx := range txl {
		txs = append(txs, tx)
	}
	return txs
}

func (txs Txs) Validate(maxSizeBytes int64) error {
	var size int64
	for _, tx := range txs {
		size += ComputeProtoSizeForTxs([]Tx{tx})
		if size > maxSizeBytes {
			return fmt.Errorf("transaction data size exceeds maximum %d", maxSizeBytes)
		}
	}
	return nil
}

// ToSliceOfBytes converts a Txs to slice of byte slices.
func (txs Txs) ToSliceOfBytes() [][]byte {
	txBzs := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		txBzs[i] = txs[i]
	}
	return txBzs
}

// TxProof represents a Merkle proof of the presence of a transaction in the Merkle tree.
type TxProof struct {
	RootHash cmtbytes.HexBytes `json:"root_hash"`
	Data     Tx                `json:"data"`
	Proof    merkle.Proof      `json:"proof"`
}

// Leaf returns the hash(tx), which is the leaf in the merkle tree which this proof refers to.
func (tp TxProof) Leaf() []byte {
	return tp.Data.Hash()
}

// Validate verifies the proof. It returns nil if the RootHash matches the dataHash argument,
// and if the proof is internally consistent. Otherwise, it returns a sensible error.
func (tp TxProof) Validate(dataHash []byte) error {
	if !bytes.Equal(dataHash, tp.RootHash) {
		return errors.New("proof matches different data hash")
	}
	if tp.Proof.Index < 0 {
		return errors.New("proof index cannot be negative")
	}
	if tp.Proof.Total <= 0 {
		return errors.New("proof total must be positive")
	}
	valid := tp.Proof.Verify(tp.RootHash, tp.Leaf())
	if valid != nil {
		return errors.New("proof is not internally consistent")
	}
	return nil
}

func (tp TxProof) ToProto() cmtproto.TxProof {
	pbProof := tp.Proof.ToProto()

	pbtp := cmtproto.TxProof{
		RootHash: tp.RootHash,
		Data:     tp.Data,
		Proof:    pbProof,
	}

	return pbtp
}

func TxProofFromProto(pb cmtproto.TxProof) (TxProof, error) {
	pbProof, err := merkle.ProofFromProto(pb.Proof)
	if err != nil {
		return TxProof{}, err
	}

	pbtp := TxProof{
		RootHash: pb.RootHash,
		Data:     pb.Data,
		Proof:    *pbProof,
	}

	return pbtp, nil
}


// // ToProto converts a types.Tx to a Protobuf Transaction.
// func (tx Tx) ToProto(signatures map[crypto.PubKey][]byte) *mempoolv1.Transaction {
// 	// Convert signatures map from crypto.PubKey to string
// 	signatureMap := make(map[string][]byte)
// 	for pubKey, signature := range signatures {
// 		// Serialize public key (e.g., Base64 encoding)
// 		pubKeyStr := base64.StdEncoding.EncodeToString(pubKey.Bytes())
// 		signatureMap[pubKeyStr] = signature
// 	}

// 	return &mempoolv1.Transaction{
// 		TransactionBytes:         transactionBytes,
// 		Signatures: signatureMap,
// 	}
// }

// func (transactionBytes *Tx) ToProtobufBytes(signatures map[crypto.PubKey][]byte) ([]byte, error) {
// 	// Convert signatures to Protobuf-compatible format
// 	signatureMap := make(map[string][]byte)
// 	for pubKey, signature := range signatures {
// 		pubKeyStr := base64.StdEncoding.EncodeToString(pubKey.Bytes())
// 		signatureMap[pubKeyStr] = signature
// 	}

// 	protoTx := &mempoolv1.Transaction{
// 		TransactionBytes:         *transactionBytes,
// 		Signatures: signatureMap,
// 	}
// 	return protoTx.Marshal()
// }

// // Unwrap Protobuf Transaction from Tx (byte array)
// func (tx Tx) FromProtobuf() (*mempoolv1.Transaction, error) {
// 	// Create an empty Protobuf Transaction object
// 	protoTx := &mempoolv1.Transaction{}

// 	// Use proto.Unmarshal to populate the Protobuf object
// 	if err := protoTx.Unmarshal([]byte(tx)); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
// 	}

// 	return protoTx, nil
// }

// ComputeProtoSizeForTxs wraps the transactions in cmtproto.Data{} and calculates the size.
// https://developers.google.com/protocol-buffers/docs/encoding
func ComputeProtoSizeForTxs(txs []Tx) int64 {
	data := Data{Txs: txs}
	pdData := data.ToProto()
	return int64(pdData.Size())
}
