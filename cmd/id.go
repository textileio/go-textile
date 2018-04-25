package cmd

import (
	"github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"

	"github.com/textileio/textile-go/core"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

func GetIds(c *ishell.Context) {

	//// textile wallet id
	//wsk, err := core.Node.UnmarshalPrivateKey()
	//if err != nil {
	//	c.Err(err)
	//	return
	//}
	//wid, err := peer.IDFromPrivateKey(wsk)
	//if err != nil {
	//	c.Err(err)
	//	return
	//}
	//mn, err := core.Node.Datastore.Config().GetMnemonic()
	//if err != nil {
	//	c.Err(err)
	//	return
	//}

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
	magenta := color.New(color.FgMagenta).SprintFunc()
	//blue := color.New(color.FgBlue).SprintFunc()
	//green := color.New(color.FgGreen).SprintFunc()
	c.Println(magenta("peer id: " + pid.Pretty()))
	//c.Println(blue("wallet id: " + wid.Pretty()))
	//c.Println(green("wallet secret: " + mn))
}
