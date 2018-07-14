package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/textileio/textile-go/crypto"
	"github.com/tyler-smith/go-bip39"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/config"
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

// IdentityConfig initializes a new identity.
func IdentityConfig(sk libp2pc.PrivKey) (config.Identity, error) {
	log.Infof("generating Ed25519 keypair for peer identity...")

	ident := config.Identity{}
	sk, pk, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return ident, err
	}

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := sk.Bytes()
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	log.Infof("new peer identity: %s\n", ident.PeerID)
	return ident, nil
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
