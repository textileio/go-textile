package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func ShowId(c *ishell.Context) {
	id, err := core.Node.Wallet.GetId()
	if err != nil {
		c.Err(err)
		return
	}
	ipfsId, err := core.Node.Wallet.GetIPFSPeerId()
	if err != nil {
		c.Err(err)
		return
	}
	red := color.New(color.FgRed).SprintFunc()
	c.Println(red(fmt.Sprintf("textile: %s, ipfs: %s", id, ipfsId)))
}
