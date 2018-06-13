package crypto

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ed25519"
)

// Implements the Ed25519 signing method
// Expects *ed25519.PrivateKey for signing and *ed25519.PublicKey for validation
type signingMethodEd25519 struct {
	Name string
}

// Specific instance for Ed25519
var SigningMethodEd25519 *signingMethodEd25519

func init() {
	SigningMethodEd25519 = &signingMethodEd25519{"Ed25519"}
	jwt.RegisterSigningMethod(SigningMethodEd25519.Alg(), func() jwt.SigningMethod {
		return SigningMethodEd25519
	})
}

func (m *signingMethodEd25519) Alg() string {
	return m.Name
}

// Implements the Verify method from SigningMethod
// For this signing method, must be an *ed25519.PublicKey structure.
func (m *signingMethodEd25519) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	var ed25519Key ed25519.PublicKey
	var ok bool

	if ed25519Key, ok = key.(ed25519.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	// Verify the signature
	if !ed25519.Verify(ed25519Key, []byte(signingString), sig) {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

// Implements the Sign method from SigningMethod
// For this signing method, must be an *ed25519.PrivateKey structure.
func (m *signingMethodEd25519) Sign(signingString string, key interface{}) (string, error) {
	var ed25519Key ed25519.PrivateKey
	var ok bool

	// Validate type of key
	if ed25519Key, ok = key.(ed25519.PrivateKey); !ok {
		return "", jwt.ErrInvalidKey
	}

	// Sign the string and return the encoded bytes
	sigBytes := ed25519.Sign(ed25519Key, []byte(signingString))
	return jwt.EncodeSegment(sigBytes), nil
}
