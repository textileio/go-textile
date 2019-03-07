package ipfs

import (
	"crypto/rand"
	"encoding/base64"

	libp2pc "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	config "gx/ipfs/QmUAuYuiafnJRZxDDX7MuruMNsicYNuyub5vUeAcupUBNs/go-ipfs-config"
	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
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
