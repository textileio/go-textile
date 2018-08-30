package crypto

import (
	"github.com/dgrijalva/jwt-go"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// Implements the Ed25519 signing method
// Expects *crypto.Ed25519PublicKey for signing and *crypto.Ed25519PublicKey for validation
type SigningMethodEd25519 struct {
	Name string
}

// Specific instance for Ed25519
var SigningMethodEd25519i *SigningMethodEd25519

func init() {
	SigningMethodEd25519i = &SigningMethodEd25519{"Ed25519"}
	jwt.RegisterSigningMethod(SigningMethodEd25519i.Alg(), func() jwt.SigningMethod {
		return SigningMethodEd25519i
	})
}

func (m *SigningMethodEd25519) Alg() string {
	return m.Name
}

// Implements the Verify method from SigningMethod
// For this signing method, must be a *crypto.Ed25519PublicKey structure.
func (m *SigningMethodEd25519) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	var ed25519Key *libp2pc.Ed25519PublicKey
	var ok bool

	if ed25519Key, ok = key.(*libp2pc.Ed25519PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	// Verify the signature
	valid, err := ed25519Key.Verify([]byte(signingString), sig)
	if err != nil {
		return err
	}
	if !valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

// Implements the Sign method from SigningMethod
// For this signing method, must be a *crypto.Ed25519PublicKey structure.
func (m *SigningMethodEd25519) Sign(signingString string, key interface{}) (string, error) {
	var ed25519Key *libp2pc.Ed25519PrivateKey
	var ok bool

	// Validate type of key
	if ed25519Key, ok = key.(*libp2pc.Ed25519PrivateKey); !ok {
		return "", jwt.ErrInvalidKey
	}

	// Sign the string and return the encoded bytes
	sigBytes, err := ed25519Key.Sign([]byte(signingString))
	if err != nil {
		return "", err
	}
	return jwt.EncodeSegment(sigBytes), nil
}
