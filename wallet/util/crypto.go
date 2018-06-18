package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/textileio/textile-go/crypto"
	"github.com/tyler-smith/go-bip39"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

// PrivKeyFromMnemonic creates a private key form a mnemonic phrase
func PrivKeyFromMnemonic(mnemonic *string) (libp2pc.PrivKey, string, error) {
	if mnemonic == nil {
		mnemonics, err := createMnemonic(bip39.NewEntropy, bip39.NewMnemonic)
		if err != nil {
			return nil, "", err
		}
		mnemonic = &mnemonics
	}

	// create the bip39 seed from the phrase
	// TODO: allow password?
	seed := bip39.NewSeed(*mnemonic, "")
	key, err := identityKeyFromSeed(seed)
	if err != nil {
		return nil, "", err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(key)
	if err != nil {
		return nil, "", err
	}
	return sk, *mnemonic, nil
}

// IDAndSecretFromMnemonic create a secret key from a mnemonic,
// returns the mnemonic (it may have been generated),
// the secret's public key (base64 string) as an id,
// and the raw secret bytes.
// Used for generating a new master identity.
func IDAndSecretFromMnemonic(mnemonic *string) (mnem string, id string, secret []byte, err error) {
	sk, mnem, err := PrivKeyFromMnemonic(mnemonic)
	if err != nil {
		return "", "", nil, err
	}
	skb, err := sk.Bytes()
	if err != nil {
		return "", "", nil, err
	}
	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		return "", "", nil, err
	}
	id = libp2pc.ConfigEncodeKey(pkb)
	return mnem, id, skb, nil
}

// GetEncryptedReaderBytes reads reader bytes and returns the encrypted result
func GetEncryptedReaderBytes(reader io.Reader, key []byte) ([]byte, error) {
	bts, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return crypto.EncryptAES(bts, key)
}

// GetNowBytes returns the current unix time as a byte string
func GetNowBytes() []byte {
	return []byte(strconv.Itoa(int(time.Now().Unix())))
}

// UnmarshalPrivateKeyFromString attempts to create a private key from a base64 encoded string
func UnmarshalPrivateKeyFromString(key string) (libp2pc.PrivKey, error) {
	keyb, err := libp2pc.ConfigDecodeKey(key)
	if err != nil {
		return nil, err
	}
	return libp2pc.UnmarshalPrivateKey(keyb)
}

// UnmarshalPublicKeyFromString attempts to create a public key from a base64 encoded string
func UnmarshalPublicKeyFromString(key string) (libp2pc.PubKey, error) {
	keyb, err := libp2pc.ConfigDecodeKey(key)
	if err != nil {
		return nil, err
	}
	return libp2pc.UnmarshalPublicKey(keyb)
}

// createMnemonic creates a new mnemonic phrase with given entropy
func createMnemonic(newEntropy func(int) ([]byte, error), newMnemonic func([]byte) (string, error)) (string, error) {
	entropy, err := newEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := newMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// identityKeyFromSeed returns a new key identity from a seed
func identityKeyFromSeed(seed []byte) ([]byte, error) {
	hm := hmac.New(sha256.New, []byte("scythian horde"))
	hm.Write(seed)
	reader := bytes.NewReader(hm.Sum(nil))
	// bits are not meaningful w/ this method in ed25519, so specify whatever
	sk, _, err := libp2pc.GenerateKeyPairWithReader(libp2pc.Ed25519, 2048, reader)
	if err != nil {
		return nil, err
	}
	encodedKey, err := sk.Bytes()
	if err != nil {
		return nil, err
	}
	return encodedKey, nil
}
