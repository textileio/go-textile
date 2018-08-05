package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/util"
	"gopkg.in/abiosoft/ishell.v2"
)

func ShowId(c *ishell.Context) {
	var id, pk string
	if len(c.Args) != 0 {
		pk = c.Args[0]
		pid, err := util.IdFromEncodedPublicKey(pk)
		if err != nil {
			c.Err(err)
			return
		}
		id = pid.Pretty()
	} else {
		var err error
		id, err = core.Node.Wallet.GetId()
		if err != nil {
			c.Err(err)
			return
		}
		pk, err = core.Node.Wallet.GetPubKeyString()
		if err != nil {
			c.Err(err)
			return
		}
	}
	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("id: %s\npk: %s", id, pk)))
}
