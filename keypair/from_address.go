package keypair

import (
	libp2pc "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	pb "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto/pb"
	"gx/ipfs/QmW7VUmSvhvSGbYbdsh7uRjhGmsYkc9fL8aJ5CorxxrU5N/go-crypto/ed25519"
	"gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"

	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/strkey"
)

// FromAddress represents a keypair to which only the address is know.  This KeyPair
// can verify signatures, but cannot sign them.
//
// NOTE: ensure the address provided is a valid strkey encoded textile address.
// Some operations will panic otherwise. It's recommended that you create these
// structs through the Parse() method.
type FromAddress struct {
	address string
}

func (kp *FromAddress) Address() string {
	return kp.address
}

func (kp *FromAddress) Hint() (r [4]byte) {
	copy(r[:], kp.publicKey()[28:])
	return
}

func (kp *FromAddress) Id() (peer.ID, error) {
	pub, err := kp.LibP2PPubKey()
	if err != nil {
		return "", nil
	}
	return peer.IDFromPublicKey(pub)
}

func (kp *FromAddress) LibP2PPrivKey() (*libp2pc.Ed25519PrivateKey, error) {
	return nil, ErrCannotSign
}

func (kp *FromAddress) LibP2PPubKey() (*libp2pc.Ed25519PublicKey, error) {
	pmes := new(pb.PublicKey)
	pmes.Data = kp.publicKey()[:]
	pk, err := libp2pc.UnmarshalEd25519PublicKey(pmes.GetData())
	if err != nil {
		return nil, err
	}
	epk, ok := pk.(*libp2pc.Ed25519PublicKey)
	if !ok {
		return nil, nil
	}
	return epk, nil
}

func (kp *FromAddress) Verify(input []byte, sig []byte) error {
	if len(sig) != ed25519.PrivateKeySize {
		return ErrInvalidSignature
	}
	var asig [ed25519.PrivateKeySize]byte
	copy(asig[:], sig[:])

	if !ed25519.Verify(kp.publicKey(), input, asig[:]) {
		return ErrInvalidSignature
	}
	return nil
}

func (kp *FromAddress) Sign(input []byte) ([]byte, error) {
	return nil, ErrCannotSign
}

func (kp *FromAddress) Encrypt(input []byte) ([]byte, error) {
	pub, err := kp.LibP2PPubKey()
	if err != nil {
		return nil, err
	}
	return crypto.Encrypt(pub, input)
}

func (kp *FromAddress) Decrypt(input []byte) ([]byte, error) {
	return nil, ErrCannotDecrypt
}

func (kp *FromAddress) publicKey() ed25519.PublicKey {
	bytes := strkey.MustDecode(strkey.VersionByteAccountID, kp.address)
	var result [ed25519.PublicKeySize]byte

	copy(result[:], bytes)

	slice := result[:]
	return slice
}
