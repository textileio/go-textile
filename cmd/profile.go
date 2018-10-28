package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/ipfs"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

func publishProfile(c *ishell.Context) {
	entry, err := core.Node.PublishProfile(nil)
	if err != nil {
		c.Err(err)
		return
	}
	if entry == nil {
		c.Println(color.New(color.FgHiRed).SprintFunc()("profile does not exist"))
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("ok, published %s -> %s", entry.Name, entry.Value)))
}

func resolveProfile(c *ishell.Context) {
	var pid peer.ID
	if len(c.Args) == 0 {
		self, err := core.Node.PeerId()
		if err != nil {
			c.Err(err)
			return
		}
		pid = self
	} else {
		var err error
		pid, err = peer.IDB58Decode(c.Args[0])
		if err != nil {
			c.Err(err)
			return
		}
	}

	entry, err := core.Node.ResolveProfile(pid)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(entry.String()))
}

func getProfile(c *ishell.Context) {
	var pid peer.ID
	if len(c.Args) == 0 {
		self, err := core.Node.PeerId()
		if err != nil {
			c.Err(err)
			return
		}
		pid = self
	} else {
		var err error
		pid, err = peer.IDB58Decode(c.Args[0])
		if err != nil {
			c.Err(err)
			return
		}
	}

	prof, err := core.Node.GetProfile(pid)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	if prof.Address != "" {
		c.Println(green(fmt.Sprintf("address:    %s", prof.Address)))
	}
	if prof.Username != "" {
		c.Println(green(fmt.Sprintf("username:   %s", prof.Username)))
	}
	if prof.AvatarUri != "" {
		c.Println(green(fmt.Sprintf("avatar_uri: %s", prof.AvatarUri)))
	}
}

func getSubs(c *ishell.Context) {
	subs, err := ipfs.IpnsSubs(core.Node.Ipfs())
	if err != nil {
		c.Err(err)
		return
	}
	green := color.New(color.FgHiGreen).SprintFunc()
	for _, sub := range subs {
		c.Println(green(sub))
	}
}

func setUsername(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing username"))
		return
	}
	id := c.Args[0]

	if err := core.Node.SetUsername(id); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("ok, updated"))
}

func setAvatar(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo id"))
		return
	}
	id := c.Args[0]

	if err := core.Node.SetAvatar(id); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("ok, updated"))
}
