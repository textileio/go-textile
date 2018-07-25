package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func GetProfile(c *ishell.Context) {
	prof, err := core.Node.Wallet.GetProfile()
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	if prof.Id != "" {
		c.Println(green(fmt.Sprintf("id:       %s", prof.Id)))
	}
	if prof.Username != "" {
		c.Println(green(fmt.Sprintf("username: %s", prof.Username)))
	}
	if prof.AvatarId != "" {
		c.Println(green(fmt.Sprintf("avatar:   %s", prof.AvatarId)))
	}
}

func PublishProfile(c *ishell.Context) {
	entry, err := core.Node.Wallet.PublishProfile()
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("ok, published %s -> %s", entry.Name, entry.Value)))
}

func ResolveProfile(c *ishell.Context) {
	var name string
	if len(c.Args) == 0 {
		id, err := core.Node.Wallet.GetId()
		if err != nil {
			c.Err(err)
			return
		}
		name = id
	} else {
		name = c.Args[0]
	}

	entry, err := core.Node.Wallet.ResolveProfile(name)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(entry.String()))
}
