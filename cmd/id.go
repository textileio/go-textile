package cmd

import (
	"github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"

	"github.com/textileio/textile-go/core"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

func ShowId(c *ishell.Context) {
	// peer id
	psk, err := core.Node.UnmarshalPrivatePeerKey()
	if err != nil {
		c.Err(err)
		return
	}
	pid, err := peer.IDFromPrivateKey(psk)
	if err != nil {
		c.Err(err)
		return
	}

	// show user their id
	red := color.New(color.FgRed).SprintFunc()
	c.Println(red(pid.Pretty()))
}
