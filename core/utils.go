package core

import (
	"crypto/sha256"
	"errors"
	cid "gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	ps "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	util "gx/ipfs/QmSU6eubNdhXjFBJBSksTp8kv8YRub8mGAPv8tVJHmL2EU/go-ipfs-util"
	mh "gx/ipfs/QmU9a9NV9RdPNwZQDYd5uKsm6N6LJLSvLbywDDYFbaaC6P/go-multihash"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"time"

	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
)

// Hash with SHA-256 and encode as a multihash
func EncodeCID(b []byte) (*cid.Cid, error) {
	multihash, err := EncodeMultihash(b)
	if err != nil {
		return nil, err
	}
	id := cid.NewCidV1(cid.Raw, *multihash)
	return id, err
}

func EncodeMultihash(b []byte) (*mh.Multihash, error) {
	h := sha256.Sum256(b)
	encoded, err := mh.Encode(h[:], mh.SHA2_256)
	if err != nil {
		return nil, err
	}
	multihash, err := mh.Cast(encoded)
	if err != nil {
		return nil, err
	}
	return &multihash, err
}

// Certain pointers, such as moderators, contain a peerID. This function
// will extract the ID from the underlying PeerInfo object.
func ExtractIDFromPointer(pi ps.PeerInfo) (string, error) {
	if len(pi.Addrs) == 0 {
		return "", errors.New("PeerInfo object has no addresses")
	}
	addr := pi.Addrs[0]
	if addr.Protocols()[0].Code != ma.P_IPFS {
		return "", errors.New("IPFS protocol not found in address")
	}
	val, err := addr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		return "", err
	}
	return val, nil
}

// FormatRFC3339PB returns the given `google_protobuf.Timestamp` as a RFC3339
// formatted string
func FormatRFC3339PB(ts google_protobuf.Timestamp) string {
	return util.FormatRFC3339(time.Unix(ts.Seconds, int64(ts.Nanos)).UTC())
}
