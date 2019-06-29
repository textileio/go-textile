package crypto

import (
	"crypto/rand"
	"fmt"

	extra "github.com/agl/ed25519/extra25519"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	"golang.org/x/crypto/nacl/box"
)

const (
	// Length of nacl nonce
	NonceBytes = 24

	// Length of nacl ephemeral public key
	EphemeralPublicKeyBytes = 32
)

var (
	// Nacl box decryption failed
	BoxDecryptionError = fmt.Errorf("failed to decrypt curve25519")
)

func Encrypt(pubKey libp2pc.PubKey, bytes []byte) ([]byte, error) {
	ed25519Pubkey, ok := pubKey.(*libp2pc.Ed25519PublicKey)
	if ok {
		return encryptCurve25519(ed25519Pubkey, bytes)
	}
	return nil, fmt.Errorf("could not determine key type")
}

func encryptCurve25519(pubKey *libp2pc.Ed25519PublicKey, bytes []byte) ([]byte, error) {
	// generated ephemeral key pair
	ephemPub, ephemPriv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// convert recipient's key into curve25519
	pk, err := publicToCurve25519(pubKey)
	if err != nil {
		return nil, err
	}

	// encrypt with nacl
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

	// prepend the ephemeral public key
	ciphertext = append(ephemPub[:], ciphertext...)

	// prepend nonce
	ciphertext = append(nonce[:], ciphertext...)
	return ciphertext, nil
}

func publicToCurve25519(k *libp2pc.Ed25519PublicKey) (*[32]byte, error) {
	var cp [32]byte
	var pk [32]byte
	r, err := k.Raw()
	if err != nil {
		return nil, err
	}
	copy(pk[:], r)
	success := extra.PublicKeyToCurve25519(&cp, &pk)
	if !success {
		return nil, fmt.Errorf("error converting ed25519 pubkey to curve25519 pubkey")
	}
	return &cp, nil
}

func Decrypt(privKey libp2pc.PrivKey, ciphertext []byte) ([]byte, error) {
	ed25519Privkey, ok := privKey.(*libp2pc.Ed25519PrivateKey)
	if ok {
		return decryptCurve25519(ed25519Privkey, ciphertext)
	}
	return nil, fmt.Errorf("could not determine key type")
}

func decryptCurve25519(privKey *libp2pc.Ed25519PrivateKey, ciphertext []byte) ([]byte, error) {
	curve25519Privkey, err := privateToCurve25519(privKey)
	if err != nil {
		return nil, err
	}

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

func privateToCurve25519(k *libp2pc.Ed25519PrivateKey) (*[32]byte, error) {
	var cs [32]byte
	r, err := k.Raw()
	if err != nil {
		return nil, err
	}
	var sk [64]byte
	copy(sk[:], r)
	extra.PrivateKeyToCurve25519(&cs, &sk)
	return &cs, nil
}

func Verify(pk libp2pc.PubKey, data []byte, sig []byte) error {
	good, err := pk.Verify(data, sig)
	if err != nil || !good {
		return fmt.Errorf("bad signature")
	}
	return nil
}
