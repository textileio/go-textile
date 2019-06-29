package keypair

import (
	"crypto/rand"
	"fmt"
	"io"

	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/textileio/go-textile/strkey"
)

var (
	// ErrInvalidKey will be returned by operations when the keypair being used
	// could not be decoded.
	ErrInvalidKey = fmt.Errorf("invalid key")

	// ErrInvalidSignature is returned when the signature is invalid, either
	// through malformation or if it does not verify the message against the
	// provided public key
	ErrInvalidSignature = fmt.Errorf("signature verification failed")

	// ErrCannotSign is returned when attempting to sign a message when
	// the keypair does not have the secret key available
	ErrCannotSign = fmt.Errorf("cannot sign")

	// ErrCannotDecrypt is returned when attempting to decrypt a message when
	// the keypair does not have the secret key available
	ErrCannotDecrypt = fmt.Errorf("cannot decrypt")
)

// KeyPair is the main interface for this package
type KeyPair interface {
	Address() string
	Hint() [4]byte
	Id() (peer.ID, error)
	LibP2PPrivKey() (*libp2pc.Ed25519PrivateKey, error)
	LibP2PPubKey() (*libp2pc.Ed25519PublicKey, error)
	Verify(input []byte, signature []byte) error
	Sign(input []byte) ([]byte, error)
	Encrypt(input []byte) ([]byte, error)
	Decrypt(input []byte) ([]byte, error)
}

// Random creates a random full keypair
func Random() *Full {
	var rawSeed [32]byte
	_, err := io.ReadFull(rand.Reader, rawSeed[:])
	if err != nil {
		panic(err)
	}

	kp, err := FromRawSeed(rawSeed)
	if err != nil {
		panic(err)
	}

	return kp
}

// Parse constructs a new KeyPair from the provided string, which should be either
// an address, or a seed. If the provided input is a seed, the resulting KeyPair
// will have signing capabilities.
func Parse(addressOrSeed string) (KeyPair, error) {
	_, err := strkey.Decode(strkey.VersionByteAccountID, addressOrSeed)
	if err == nil {
		return &FromAddress{addressOrSeed}, nil
	}

	if err != strkey.ErrInvalidVersionByte {
		return nil, err
	}

	_, err = strkey.Decode(strkey.VersionByteSeed, addressOrSeed)
	if err == nil {
		return &Full{addressOrSeed}, nil
	}

	return nil, err
}

// MustParse is the panic-on-fail version of Parse
func MustParse(addressOrSeed string) KeyPair {
	kp, err := Parse(addressOrSeed)
	if err != nil {
		panic(err)
	}

	return kp
}

// FromRawSeed creates a new keypair from the provided raw ED25519 seed
func FromRawSeed(rawSeed [32]byte) (*Full, error) {
	seed, err := strkey.Encode(strkey.VersionByteSeed, rawSeed[:])
	if err != nil {
		return nil, err
	}

	return &Full{seed}, nil
}
