package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"

	//"github.com/libp2p/go-libp2p-crypto"
	//"github.com/libp2p/go-libp2p-host"
	//"github.com/libp2p/go-libp2p-net"
	//"github.com/libp2p/go-libp2p-peer"
	//"github.com/libp2p/go-libp2p-peerstore"
	//"github.com/libp2p/go-libp2p-swarm"
	//"github.com/libp2p/go-libp2p/p2p/host/basic"
	//"github.com/multiformats/go-multiaddr"
	//"gx/ipfs/QmSwZMWwFZSUpe5muU2xgTUwppH24KfMwdPXiwbEp2c6G5/go-libp2p-swarm"
	//"gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p/p2p/host/basic"
	"gx/ipfs/QmSwZMWwFZSUpe5muU2xgTUwppH24KfMwdPXiwbEp2c6G5/go-libp2p-swarm"
	"gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p/p2p/host/basic"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
	"gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	"gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
)

/*
* addAddrToPeerstore parses a peer multiaddress and adds
* it to the given host's peerstore, so it knows how to
* contact it. It returns the peer ID of the remote peer.
* @credit examples/http-proxy/proxy.go
 */
func addAddrToPeerstore(h host.Host, addr string) peer.ID {
	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Fatalln(err)
	}
	pid, err := ipfsaddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add
	// it to the peerstore so LibP2P knows how to contact it
	h.Peerstore().AddAddr(peerid, targetAddr, peerstore.PermanentAddrTTL)
	return peerid
}

func handleStream(s net.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}

}

// TODO: This should be broken down, converted to use ipfs info, and moved to mobile/node.go
func main(sourcePort int, dest string, debug bool) {

	r := rand.Reader

	// TODO: Use already generated in ipfs
	// Creates a new RSA key pair for this host
	prvKey, pubKey, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {
		panic(err)
	}

	// Getting host ID from public key.
	// host ID is the hash of public key
	nodeID, _ := peer.IDFromPublicKey(pubKey)

	// 0.0.0.0 will listen on any interface device
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))

	// Adding self to the peerstore.
	ps := peerstore.NewPeerstore()
	ps.AddPrivKey(nodeID, prvKey)
	ps.AddPubKey(nodeID, pubKey)

	// Creating a new Swarm network.
	network, err := swarm.NewNetwork(context.Background(), []multiaddr.Multiaddr{sourceMultiAddr}, nodeID, ps, nil)

	if err != nil {
		panic(err)
	}

	// NewHost constructs a new *BasicHost and activates it by attaching its
	// stream and connection handlers to the given inet.Network (network).
	// Other options like NATManager can also be added here.
	// See docs: https://godoc.org/github.com/libp2p/go-libp2p/p2p/host/basic#HostOpts
	host := basichost.New(network)

	if *dest == "" {
		// Set a function as stream handler.
		// This function  is called when a peer initiate a connection and starts a stream with this peer.
		// Only applicable on the receiving side.
		host.SetStreamHandler("/chat/1.0.0", handleStream)

		fmt.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%d/ipfs/%s' on another console.\n You can replace 127.0.0.1 with public IP as well.\n", *sourcePort, host.ID().Pretty())
		fmt.Printf("\nWaiting for incoming connection\n\n")
		// Hang forever
		<-make(chan struct{})

	} else {

		// Add destination peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		peerID := addAddrToPeerstore(host, *dest)

		fmt.Println("This node's multiaddress: ")
		// IP will be 0.0.0.0 (listen on any interface) and port will be 0 (choose one for me).
		// Although this node will not listen for any connection. It will just initiate a connect with
		// one of its peer and use that stream to communicate.
		fmt.Printf("%s/ipfs/%s\n", sourceMultiAddr, host.ID().Pretty())

		// Start a stream with peer with peer Id: 'peerId'.
		// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
		s, err := host.NewStream(context.Background(), peerID, "/chat/1.0.0")

		if err != nil {
			panic(err)
		}

		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		// Hang forever.
		select {}

	}
}