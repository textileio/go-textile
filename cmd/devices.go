package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

func ListDevices(c *ishell.Context) {
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

func AddDevice(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing device name"))
		return
	}
	name := c.Args[0]
	if len(c.Args) == 1 {
		c.Err(errors.New("missing device pub key"))
		return
	}
	pks := c.Args[1]

	pkb, err := libp2pc.ConfigDecodeKey(pks)
	if err != nil {
		c.Err(err)
		return
	}
	pk, err := libp2pc.UnmarshalPublicKey(pkb)
	if err != nil {
		c.Err(err)
		return
	}

	err = core.Node.AddDevice(name, pk)
	if err != nil {
		c.Err(err)
		return
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("added device '%s'", name)))
}

func RemoveDevice(c *ishell.Context) {
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
