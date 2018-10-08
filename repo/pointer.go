package repo

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	routing "gx/ipfs/QmTktQYCKzQjhxF6dk5xJPRuhHn3JBiKGvMLoiDy1mYmxC/go-libp2p-kad-dht"
	pb "gx/ipfs/QmTktQYCKzQjhxF6dk5xJPRuhHn3JBiKGvMLoiDy1mYmxC/go-libp2p-kad-dht/pb"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	ps "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"strconv"
	"sync"
	"time"
)

const MAGIC string = "000000000000000000000000"

type Purpose int

const (
	MESSAGE   Purpose = 1
	MODERATOR Purpose = 2
	TAG       Purpose = 3
	CHANNEL   Purpose = 4
)

/* A pointer is a custom provider inserted into the DHT which points to a location of a file.
   For offline messaging purposes we use a hash of the recipient's ID as the key and set the
   provider to the location of the ciphertext. We set the Peer ID of the provider object to
   a magic number so we distinguish it from regular providers and use a longer ttl.
   Note this will only be compatible with the OpenBazaar/go-ipfs fork. */
type Pointer struct {
	Cid      *cid.Cid
	Value    ps.PeerInfo
	Purpose  Purpose
	Date     time.Time
	CancelId *peer.ID
}

// entropy is a sequence of bytes that should be deterministic based on the content of the pointer
// it is hashed and used to fill the remaining 20 bytes of the magic id
func NewPointer(mhKey multihash.Multihash, prefixLen int, addr ma.Multiaddr, entropy []byte) (Pointer, error) {
	keyhash := CreatePointerKey(mhKey, prefixLen)
	k, err := cid.Decode(keyhash.B58String())
	if err != nil {
		return Pointer{}, err
	}

	magicID, err := getMagicId(entropy)
	if err != nil {
		return Pointer{}, err
	}
	pi := ps.PeerInfo{
		ID:    magicID,
		Addrs: []ma.Multiaddr{addr},
	}
	return Pointer{Cid: k, Value: pi}, nil
}

func PublishPointer(node *core.IpfsNode, ctx context.Context, pointer Pointer) error {
	return addPointer(node, ctx, pointer.Cid, pointer.Value)
}

func FindPointersAsync(dht *routing.IpfsDHT, ctx context.Context, mhKey multihash.Multihash, prefixLen int) <-chan ps.PeerInfo {
	keyhash := CreatePointerKey(mhKey, prefixLen)
	key, _ := cid.Decode(keyhash.B58String())
	peerout := dht.FindProvidersAsync(ctx, key, 100000)
	return peerout
}

func FindPointers(dht *routing.IpfsDHT, ctx context.Context, mhKey multihash.Multihash, prefixLen int) ([]ps.PeerInfo, error) {
	var providers []ps.PeerInfo
	for p := range FindPointersAsync(dht, ctx, mhKey, prefixLen) {
		providers = append(providers, p)
	}
	return providers, nil
}

func PutPointerToPeer(node *core.IpfsNode, ctx context.Context, peer peer.ID, pointer Pointer) error {
	dht := node.Routing.(*routing.IpfsDHT)
	return putPointer(ctx, dht, peer, pointer.Value, pointer.Cid.KeyString())
}

//func GetPointersFromPeer(node *core.IpfsNode, ctx context.Context, p peer.ID, key *cid.Cid) ([]*ps.PeerInfo, error) {
//	dht := node.Routing.(*routing.IpfsDHT)
//	pmes := pb.NewMessage(pb.Message_GET_PROVIDERS, key.KeyString(), 0)
//	resp, err := dht.SendRequest(ctx, p, pmes)
//	if err != nil {
//		return []*ps.PeerInfo{}, err
//	}
//	return pb.PBPeersToPeerInfos(resp.GetProviderPeers()), nil
//}

func addPointer(node *core.IpfsNode, ctx context.Context, k *cid.Cid, pi ps.PeerInfo) error {
	dht := node.Routing.(*routing.IpfsDHT)
	peers, err := dht.GetClosestPeers(ctx, k.KeyString())
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for p := range peers {
		wg.Add(1)
		go func(p peer.ID) {
			defer wg.Done()
			putPointer(ctx, dht, p, pi, k.KeyString())
		}(p)
	}
	wg.Wait()
	return nil
}

func putPointer(ctx context.Context, dht *routing.IpfsDHT, p peer.ID, pi ps.PeerInfo, skey string) error {
	pmes := pb.NewMessage(pb.Message_ADD_PROVIDER, skey, 0)
	pmes.ProviderPeers = pb.RawPeerInfosToPBPeers([]ps.PeerInfo{pi})

	//err := dht.SendMessage(ctx, p, pmes)
	//if err != nil {
	//	return err
	//}
	return nil
}

func CreatePointerKey(mh multihash.Multihash, prefixLen int) multihash.Multihash {
	// Grab the first 8 bytes from the multihash digest
	m, _ := multihash.Decode(mh)
	prefix64 := binary.BigEndian.Uint64(m.Digest[:8])

	// Convert to binary string
	bin := strconv.FormatUint(prefix64, 2)

	// Pad with leading zeros
	leadingZeros := 64 - len(bin)
	for i := 0; i < leadingZeros; i++ {
		bin = "0" + bin
	}

	// Grab the bits corresponding to the prefix length and convert to int
	intPrefix, _ := strconv.ParseUint(bin[:prefixLen], 2, 64)

	// Convert to 8 byte array
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, intPrefix)

	// Hash the array
	hash := sha256.New()
	hash.Write(bs)
	md := hash.Sum(nil)

	// Encode as multihash
	keyHash, _ := multihash.Encode(md, multihash.SHA2_256)
	return keyHash
}

func getMagicId(entropy []byte) (peer.ID, error) {
	magicBytes, err := hex.DecodeString(MAGIC)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	hash.Write(entropy)
	hashedEntropy := hash.Sum(nil)
	magicBytes = append(magicBytes, hashedEntropy[:20]...)
	h, err := multihash.Encode(magicBytes, multihash.SHA2_256)
	if err != nil {
		return "", err
	}
	id, err := peer.IDFromBytes(h)
	if err != nil {
		return "", err
	}
	return id, nil
}
