package cmd

import (
	"errors"

	"github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"

	"github.com/textileio/textile-go/core"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

func GetIds(c *ishell.Context) {
	if !core.Node.IsDatastoreConfigured() {
		c.Err(errors.New("datastore not initialized, please run textile init"))
		return
	}

	// textile wallet id
	wsk, err := core.Node.UnmarshalPrivateKey()
	if err != nil {
		c.Err(err)
		return
	}
	wid, err := peer.IDFromPrivateKey(wsk)
	if err != nil {
		c.Err(err)
		return
	}

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

	// show user their id's
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	c.Println(blue("wallet id: " + wid.Pretty()))
	c.Println(magenta("peer id: " + pid.Pretty()))
}
