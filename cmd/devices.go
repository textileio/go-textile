package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

func listDevices(c *ishell.Context) {
	devices := core.Node.Devices()
	if len(devices) == 0 {
		c.Println("no devices found")
	} else {
		c.Println(fmt.Sprintf("found %v devices", len(devices)))
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	for _, dev := range devices {
		c.Println(blue(fmt.Sprintf("name: %s, id: %s", dev.Name, dev.Id)))
	}
}

func addDevice(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing device name"))
		return
	}
	name := c.Args[0]
	if len(c.Args) == 1 {
		c.Err(errors.New("missing device id"))
		return
	}

	did, err := peer.IDB58Decode(c.Args[1])
	if err != nil {
		c.Err(err)
		return
	}

	if err := core.Node.AddDevice(name, did); err != nil {
		c.Err(err)
		return
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("added device '%s'", name)))
}

func removeDevice(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing device id"))
		return
	}
	id := c.Args[0]

	err := core.Node.RemoveDevice(id)
	if err != nil {
		c.Err(err)
		return
	}

	red := color.New(color.FgHiRed).SprintFunc()
	c.Println(red(fmt.Sprintf("removed device '%s'", id)))
}
