package crypto_implementation

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
)

// PubKeyFromBytes creates a crypto.PubKey from its serialized byte representation.
func PubKeyFromBytes(pubKeyBytes []byte) (crypto.PubKey, error) {
	if len(pubKeyBytes) != ed25519.PubKeySize {
		return nil, fmt.Errorf("invalid public key size: got %d, expected %d", len(pubKeyBytes), ed25519.PubKeySize)
	}
	return ed25519.PubKey(pubKeyBytes), nil
}

func ConvertSignatures(signatures map[string][]byte) (map[crypto.PubKey][]byte, error) {
	converted := make(map[crypto.PubKey][]byte, len(signatures))

	for key, value := range signatures {
		pubKey, err := PubKeyFromBytes([]byte(key))
		if err != nil {
			return nil, fmt.Errorf("failed to convert public key: %w", err)
		}
		converted[pubKey] = value
	}

	return converted, nil
}
