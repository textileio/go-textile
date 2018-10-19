package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

func listAccountPeers(c *ishell.Context) {
	peers := core.Node.AccountPeers()
	if len(peers) == 0 {
		c.Println("no peers found")
	} else {
		c.Println(fmt.Sprintf("found %v peers", len(peers)))
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	for _, dev := range peers {
		c.Println(blue(fmt.Sprintf("name: %s, id: %s", dev.Name, dev.Id)))
	}
}

func addAccountPeer(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing account peer name"))
		return
	}
	name := c.Args[0]
	if len(c.Args) == 1 {
		c.Err(errors.New("missing account peer id"))
		return
	}

	pid, err := peer.IDB58Decode(c.Args[1])
	if err != nil {
		c.Err(err)
		return
	}

	if err := core.Node.AddAccountPeer(pid, name); err != nil {
		c.Err(err)
		return
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("added account peer '%s'", name)))
}

func removeAccountPeer(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing account peer id"))
		return
	}
	id := c.Args[0]

	err := core.Node.RemoveAccountPeer(id)
	if err != nil {
		c.Err(err)
		return
	}

	red := color.New(color.FgHiRed).SprintFunc()
	c.Println(red(fmt.Sprintf("removed account peer '%s'", id)))
}
