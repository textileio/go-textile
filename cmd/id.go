package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/ipfs"
	"gopkg.in/abiosoft/ishell.v2"
)

func showId(c *ishell.Context) {
	grey := color.New(color.FgHiBlack).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()

	// check for an input to convert
	if len(c.Args) != 0 {
		pk := c.Args[0]
		pid, err := ipfs.IdFromEncodedPublicKey(pk)
		if err != nil {
			c.Err(err)
			return
		}
		id := pid.Pretty()
		c.Println(green(fmt.Sprintf("id: %s\npk: %s", id, pk)))
		return
	}

	// get account
	accnt, err := core.Node.Account()
	if err != nil {
		c.Err(err)
		return
	}
	accntId, err := core.Node.Id()
	if err != nil {
		c.Err(err)
		return
	}

	// get peer id / pk
	pid, err := core.Node.GetPeerId()
	if err != nil {
		c.Err(err)
		return
	}

	c.Println(grey("--- ACCOUNT ---"))
	c.Println(cyan(fmt.Sprintf("Id: %s", accntId.Pretty())))
	c.Println(cyan(fmt.Sprintf("Address: %s", accnt.Address())))
	c.Println(cyan(fmt.Sprintf("Seed: %s", accnt.Seed())))
	c.Println(grey("--- PEER ---"))
	c.Println(green(fmt.Sprintf("Id: %s", pid)))
}
