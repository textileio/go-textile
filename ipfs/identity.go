package ipfs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	peer "gx/ipfs/QmQnuSxgSFubscHgkgSeayLxKmVcmNhFUaZw4gHtV3tJ15/go-libp2p-peer"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
	libp2p "gx/ipfs/Qmb9aAJwV1mDc5iPNtVuzVvsNiKA6kkDpZspMUgVfXPVc8/go-libp2p-crypto"
)

func IdentityFromKey(privkey []byte) (config.Identity, error) {

	ident := config.Identity{}
	sk, err := libp2p.UnmarshalPrivateKey(privkey)
	if err != nil {
		return ident, err
	}
	skbytes, err := sk.Bytes()
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(sk.GetPublic())
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	return ident, nil
}

func IdentityKeyFromSeed(seed []byte, bits int) ([]byte, error) {
	hmac := hmac.New(sha256.New, []byte("OpenBazaar seed"))
	hmac.Write(seed)
	reader := bytes.NewReader(hmac.Sum(nil))
	sk, _, err := libp2p.GenerateKeyPairWithReader(libp2p.Ed25519, bits, reader)
	if err != nil {
		return nil, err
	}
	encodedKey, err := sk.Bytes()
	if err != nil {
		return nil, err
	}
	return encodedKey, nil
}
