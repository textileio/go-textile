package cmd

import (
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func ListRooms(c *ishell.Context) {

	// TODO: reconcile this with galleries table
	rooms := core.Node.IpfsNode.Floodsub.GetTopics()

	// show user
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, r := range rooms {
		c.Println(yellow(r))
	}
}
