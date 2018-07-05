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
	pk, err := core.Node.Wallet.GetPubKeyString()
	if err != nil {
		c.Err(err)
		return
	}
	red := color.New(color.FgRed).SprintFunc()
	c.Println(red(fmt.Sprintf("id: %s, pk: %s", id, pk)))
}
