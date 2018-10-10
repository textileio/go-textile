package net

import (
	"github.com/golang/protobuf/proto"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/pb"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

var log = logging.MustGetLogger("net")

// newEnvelope returns a wrapper around a signed message
func newEnvelope(sk libp2pc.PrivKey, message *pb.Message) (*pb.Envelope, error) {
	serialized, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	authorSig, err := sk.Sign(serialized)
	if err != nil {
		return nil, err
	}
	authorPk, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	return &pb.Envelope{Message: message, Pk: authorPk, Sig: authorSig}, nil
}
