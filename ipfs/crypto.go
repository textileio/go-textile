package ipfs

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/textileio/textile-go/wallet"
	"github.com/tyler-smith/go-bip39"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/config"
)

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
	pkbytes, err := pk.Bytes()
	if err != nil {
		return ident, err
	}
	pks := base64.StdEncoding.EncodeToString(pkbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	log.Infof("new peer identity: id: %s, pk: %s", ident.PeerID, pks)
	return ident, nil
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

// IdFromEncodedPublicKey return the underlying id from an encoded public key
func IdFromEncodedPublicKey(key string) (peer.ID, error) {
	pk, err := UnmarshalPublicKeyFromString(key)
	if err != nil {
		return "", err
	}
	return peer.IDFromPublicKey(pk)
}

// EncodeKey returns a base64 encoded key
func EncodeKey(key libp2pc.Key) (string, error) {
	keyb, err := key.Bytes()
	if err != nil {
		return "", err
	}
	return libp2pc.ConfigEncodeKey(keyb), nil
}

// DecodePrivKey returns a private key from a base64 encoded string
func DecodePrivKey(key string) (libp2pc.PrivKey, error) {
	keyb, err := libp2pc.ConfigDecodeKey(key)
	if err != nil {
		return nil, err
	}
	return libp2pc.UnmarshalPrivateKey(keyb)
}

// DecodePubKey returns a public key from a base64 encoded string
func DecodePubKey(key string) (libp2pc.PubKey, error) {
	keyb, err := libp2pc.ConfigDecodeKey(key)
	if err != nil {
		return nil, err
	}
	return libp2pc.UnmarshalPublicKey(keyb)
}

// CreateMnemonic creates a new mnemonic phrase with given bit size
func CreateMnemonic(count wallet.WordCount) (string, error) {
	entropy, err := bip39.NewEntropy(int(count))
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// PrivKeyFromMnemonic creates a private key form a mnemonic phrase
func PrivKeyFromMnemonic(mnemonic *string) (libp2pc.PrivKey, string, error) {
	if mnemonic == nil {
		mnemonics, err := CreateMnemonic(wallet.TwentyFourWords)
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
