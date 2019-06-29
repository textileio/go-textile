package ipfs

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/golang/protobuf/proto"
	config "github.com/ipfs/go-ipfs-config"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/crypto/ed25519"
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

// UnmarshalPrivateKey converts a protobuf serialized private key into its
// representative object
func UnmarshalPrivateKey(data []byte) (libp2pc.PrivKey, error) {
	pmes := new(pb.PrivateKey)
	err := proto.Unmarshal(data, pmes)
	if err != nil {
		return nil, err
	}

	um, ok := libp2pc.PrivKeyUnmarshallers[pmes.GetType()]
	if !ok {
		return nil, libp2pc.ErrBadKeyType
	}

	// Manually shorten key length becuase libp2p backwards compat test will not catch our keys
	// since they do not have the redundant public key, just empty bytes.
	pd := pmes.GetData()
	if len(pd) == ed25519.PrivateKeySize+ed25519.PublicKeySize {
		k := make([]byte, ed25519.PrivateKeySize)
		copy(k, pd[:ed25519.PrivateKeySize])
		pd = k
	}

	return um(pd)
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
