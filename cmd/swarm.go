package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"

	"github.com/textileio/textile-go/core"

	iaddr "gx/ipfs/QmQViVWBHbU6HmYjXcdNq7tVASCNgdg64ZGcauuDkLCivW/go-ipfs-addr"
	"gx/ipfs/QmSwZMWwFZSUpe5muU2xgTUwppH24KfMwdPXiwbEp2c6G5/go-libp2p-swarm"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmfAkMSt9Fwzk48QDJecPcwCUjnf2uG7MLnmCGTp4C6ouL/go-ipfs-cmds"
	"strconv"
)

func SwarmConnect(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer address"))
		return
	}
	addrs := c.Args

	if core.Node.IpfsNode.PeerHost == nil {
		c.Err(errors.New("not online"))
		return
	}

	snet, ok := core.Node.IpfsNode.PeerHost.Network().(*swarm.Network)
	if !ok {
		c.Err(errors.New("peerhost network was not swarm"))
		return
	}

	swrm := snet.Swarm()

	pis, err := peersWithAddresses(addrs)
	if err != nil {
		c.Err(err)
		return
	}

	output := make([]string, len(pis))
	for i, pi := range pis {
		swrm.Backoff().Clear(pi.ID)

		output[i] = "connect " + pi.ID.Pretty()

		err := core.Node.IpfsNode.PeerHost.Connect(core.Node.IpfsNode.Context(), pi)
		if err != nil {
			c.Err(fmt.Errorf("%s failure: %s", output[i], err))
			return
		}
		output[i] += " success"
	}

	// show user their id
	red := color.New(color.FgRed).SprintFunc()
	for _, o := range output {
		c.Println(red(o))
	}
}

func SwarmPing(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer address"))
		return
	}
	addrs := c.Args[0]
	num := 1
	if len(c.Args) > 1 {
		parsed, err := strconv.ParseInt(c.Args[1], 10, 64)
		if err != nil {
			c.Err(err)
			return
		}
		num = int(parsed)
	}

	out := make(chan string)
	go func() {
		err := core.Node.PingPeer(addrs, num, out)
		if err != nil {
			c.Err(err)
		}
	}()

	green := color.New(color.FgHiGreen).SprintFunc()
	cnt := 0
	for {
		select {
		case msg, ok := <-out:
			if !ok {
				return
			}
			c.Println(green(msg))
			cnt++
			if cnt == num {
				return
			}
		}
	}
}

// parseAddresses is a function that takes in a slice of string peer addresses
// (multiaddr + peerid) and returns slices of multiaddrs and peerids.
func parseAddresses(addrs []string) (iaddrs []iaddr.IPFSAddr, err error) {
	iaddrs = make([]iaddr.IPFSAddr, len(addrs))
	for i, saddr := range addrs {
		iaddrs[i], err = iaddr.ParseString(saddr)
		if err != nil {
			return nil, cmds.ClientError("invalid peer address: " + err.Error())
		}
	}
	return
}

// peersWithAddresses is a function that takes in a slice of string peer addresses
// (multiaddr + peerid) and returns a slice of properly constructed peers
func peersWithAddresses(addrs []string) (pis []pstore.PeerInfo, err error) {
	iaddrs, err := parseAddresses(addrs)
	if err != nil {
		return nil, err
	}

	for _, a := range iaddrs {
		pis = append(pis, pstore.PeerInfo{
			ID:    a.ID(),
			Addrs: []ma.Multiaddr{a.Transport()},
		})
	}
	return pis, nil
}
