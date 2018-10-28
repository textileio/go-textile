package crypto

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/nacl/box"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

const (
	// Length of nacl nonce
	NonceBytes = 24

	// Length of nacl ephemeral public key
	EphemeralPublicKeyBytes = 32
)

var (
	// Nacl box decryption failed
	BoxDecryptionError = errors.New("failed to decrypt curve25519")
)

func Encrypt(pubKey libp2pc.PubKey, bytes []byte) ([]byte, error) {
	ed25519Pubkey, ok := pubKey.(*libp2pc.Ed25519PublicKey)
	if ok {
		return encryptCurve25519(ed25519Pubkey, bytes)
	}
	return nil, errors.New("could not determine key type")
}

func encryptCurve25519(pubKey *libp2pc.Ed25519PublicKey, bytes []byte) ([]byte, error) {
	// Generated ephemeral key pair
	ephemPub, ephemPriv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	// Convert recipient's key into curve25519
	pk, err := pubKey.ToCurve25519()
	if err != nil {
		return nil, err
	}

	// Encrypt with nacl
	var ciphertext []byte
	var nonce [24]byte
	n := make([]byte, 24)
	_, err = rand.Read(n)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 24; i++ {
		nonce[i] = n[i]
	}
	ciphertext = box.Seal(ciphertext, bytes, &nonce, pk, ephemPriv)

	// Prepend the ephemeral public key
	ciphertext = append(ephemPub[:], ciphertext...)

	// Prepend nonce
	ciphertext = append(nonce[:], ciphertext...)
	return ciphertext, nil
}

func Decrypt(privKey libp2pc.PrivKey, ciphertext []byte) ([]byte, error) {
	ed25519Privkey, ok := privKey.(*libp2pc.Ed25519PrivateKey)
	if ok {
		return decryptCurve25519(ed25519Privkey, ciphertext)
	}
	return nil, errors.New("could not determine key type")
}

func decryptCurve25519(privKey *libp2pc.Ed25519PrivateKey, ciphertext []byte) ([]byte, error) {
	curve25519Privkey := privKey.ToCurve25519()
	var plaintext []byte

	n := ciphertext[:NonceBytes]
	ephemPubkeyBytes := ciphertext[NonceBytes : NonceBytes+EphemeralPublicKeyBytes]
	ct := ciphertext[NonceBytes+EphemeralPublicKeyBytes:]

	var ephemPubkey [32]byte
	for i := 0; i < 32; i++ {
		ephemPubkey[i] = ephemPubkeyBytes[i]
	}

	var nonce [24]byte
	for i := 0; i < 24; i++ {
		nonce[i] = n[i]
	}

	plaintext, success := box.Open(plaintext, ct, &nonce, &ephemPubkey, curve25519Privkey)
	if !success {
		return nil, BoxDecryptionError
	}
	return plaintext, nil
}

// verify verifies a signature
func Verify(pk libp2pc.PubKey, data []byte, sig []byte) error {
	good, err := pk.Verify(data, sig)
	if err != nil || !good {
		return errors.New("bad signature")
	}
	return nil
}
