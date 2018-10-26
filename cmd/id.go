package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func showId(c *ishell.Context) {
	grey := color.New(color.FgHiBlack).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()

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
	pid, err := core.Node.PeerId()
	if err != nil {
		c.Err(err)
		return
	}

	c.Println(grey("--- PEER ---"))
	c.Println(green(fmt.Sprintf("Id: %s", pid.Pretty())))
	c.Println(grey("--- ACCOUNT ---"))
	c.Println(cyan(fmt.Sprintf("Id: %s", accntId.Pretty())))
	c.Println(cyan(fmt.Sprintf("Address: %s", accnt.Address())))
	c.Println(cyan(fmt.Sprintf("Seed: %s", accnt.Seed())))
}
