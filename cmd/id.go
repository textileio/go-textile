package cmd

import (
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func ShowId(c *ishell.Context) {
	id, err := core.Node.Wallet.GetIPFSPeerID()
	if err != nil {
		c.Err(err)
		return
	}
	red := color.New(color.FgRed).SprintFunc()
	c.Println(red(id))
}
