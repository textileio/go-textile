package keypair

import (
	"bytes"

	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/strkey"
	"golang.org/x/crypto/ed25519"
)

type Full struct {
	seed string
}

func (kp *Full) Address() string {
	return strkey.MustEncode(strkey.VersionByteAccountID, kp.publicKey()[:])
}

func (kp *Full) Hint() (r [4]byte) {
	copy(r[:], kp.publicKey()[28:])
	return
}

func (kp *Full) Seed() string {
	return kp.seed
}

func (kp *Full) Id() (peer.ID, error) {
	pub, err := kp.LibP2PPubKey()
	if err != nil {
		return "", nil
	}
	return peer.IDFromPublicKey(pub)
}

func (kp *Full) LibP2PPrivKey() (*libp2pc.Ed25519PrivateKey, error) {
	buf := make([]byte, ed25519.PrivateKeySize)
	copy(buf, kp.rawSeed()[:])
	copy(buf[ed25519.PrivateKeySize-ed25519.PublicKeySize:], kp.publicKey()[:])
	pmes := new(pb.PrivateKey)
	pmes.Data = buf
	sk, err := libp2pc.UnmarshalEd25519PrivateKey(pmes.GetData())
	if err != nil {
		return nil, err
	}
	esk, ok := sk.(*libp2pc.Ed25519PrivateKey)
	if !ok {
		return nil, nil
	}
	return esk, nil
}

func (kp *Full) LibP2PPubKey() (*libp2pc.Ed25519PublicKey, error) {
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

func (kp *Full) Verify(input []byte, sig []byte) error {
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

func (kp *Full) Sign(input []byte) ([]byte, error) {
	_, priv := kp.keys()
	return ed25519.Sign(priv, input)[:], nil
}

func (kp *Full) Encrypt(input []byte) ([]byte, error) {
	pub, err := kp.LibP2PPubKey()
	if err != nil {
		return nil, err
	}
	return crypto.Encrypt(pub, input)
}

func (kp *Full) Decrypt(input []byte) ([]byte, error) {
	priv, err := kp.LibP2PPrivKey()
	if err != nil {
		return nil, err
	}
	return crypto.Decrypt(priv, input)
}

func (kp *Full) publicKey() ed25519.PublicKey {
	pub, _ := kp.keys()
	return pub
}

func (kp *Full) keys() (ed25519.PublicKey, ed25519.PrivateKey) {
	reader := bytes.NewReader(kp.rawSeed())
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		panic(err)
	}
	return pub, priv
}

func (kp *Full) rawSeed() []byte {
	return strkey.MustDecode(strkey.VersionByteSeed, kp.seed)
}
