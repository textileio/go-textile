package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmSwZMWwFZSUpe5muU2xgTUwppH24KfMwdPXiwbEp2c6G5/go-libp2p-swarm"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"sort"
	"strconv"
	"strings"
)

type streamInfo struct {
	Protocol string
}

type connInfo struct {
	Addr    string
	Peer    string
	Latency string
	Muxer   string
	Streams []streamInfo
}

func (ci *connInfo) Less(i, j int) bool {
	return ci.Streams[i].Protocol < ci.Streams[j].Protocol
}

func (ci *connInfo) Len() int {
	return len(ci.Streams)
}

func (ci *connInfo) Swap(i, j int) {
	ci.Streams[i], ci.Streams[j] = ci.Streams[j], ci.Streams[i]
}

type connInfos struct {
	Peers []connInfo
}

func (ci connInfos) Less(i, j int) bool {
	return ci.Peers[i].Addr < ci.Peers[j].Addr
}

func (ci connInfos) Len() int {
	return len(ci.Peers)
}

func (ci connInfos) Swap(i, j int) {
	ci.Peers[i], ci.Peers[j] = ci.Peers[j], ci.Peers[i]
}

func SwarmPeers(c *ishell.Context) {
	conns, err := core.Node.Peers()
	if err != nil {
		c.Err(core.ErrOffline)
		return
	}

	var out connInfos
	for _, c := range conns {
		pid := c.RemotePeer()
		addr := c.RemoteMultiaddr()

		ci := connInfo{
			Addr: addr.String(),
			Peer: pid.Pretty(),
		}

		swcon, ok := c.(*swarm.Conn)
		if ok {
			ci.Muxer = fmt.Sprintf("%T", swcon.StreamConn().Conn())
		}

		sort.Sort(&ci)
		out.Peers = append(out.Peers, ci)
	}
	sort.Sort(&out)

	cyan := color.New(color.FgHiCyan).SprintFunc()
	pipfs := ma.ProtocolWithCode(ma.P_IPFS).Name
	for _, info := range out.Peers {
		ids := fmt.Sprintf("/%s/%s", pipfs, info.Peer)
		if strings.HasSuffix(info.Addr, ids) {
			c.Print(cyan(fmt.Sprintf("%s", info.Addr)))
		} else {
			c.Print(cyan(fmt.Sprintf("%s%s", info.Addr, ids)))
		}
		if info.Latency != "" {
			c.Print(cyan(fmt.Sprintf(" %s", info.Latency)))
		}
		c.Print("\n")

		for _, s := range info.Streams {
			if s.Protocol == "" {
				s.Protocol = "<no protocol name>"
			}

			c.Printf(cyan(fmt.Sprintf("  %s\n", s.Protocol)))
		}
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

func SwarmConnect(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer address"))
		return
	}
	addrs := c.Args

	output, err := core.Node.ConnectPeer(addrs)
	if err != nil {
		c.Err(err)
		return
	}

	// show user their id
	red := color.New(color.FgRed).SprintFunc()
	for _, o := range output {
		c.Println(red(o))
	}
}
