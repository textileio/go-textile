package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func ListAlbums(c *ishell.Context) {
	// cross check pubsub rooms and albums
	rooms := core.Node.IpfsNode.Floodsub.GetTopics()
	albums := core.Node.Datastore.Albums().GetAlbums("")

	yellow := color.New(color.FgYellow).SprintFunc()
	for _, a := range albums {
		mem := "disabled"
		for _, r := range rooms {
			if r == a.Id {
				mem = "enabled"
			}
		}
		c.Println(yellow(fmt.Sprintf("name: %s, id: %s, status: %s", a.Name, a.Id, mem)))
	}
}

func CreateAlbum(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing album name"))
		return
	}
	name := c.Args[0]
	if err := core.Node.CreateAlbum("", name); err != nil {
		c.Err(err)
	}
}

func EnableAlbum(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing album name"))
		return
	}
	name := c.Args[0]

	a := core.Node.Datastore.Albums().GetAlbumByName(name)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find album: %s", name)))
		return
	}

	if core.Node.LeftRoomChs[a.Id] != nil {
		c.Printf("album already enabled: %s\n", a.Id)
		return
	}

	go core.Node.JoinRoom(a.Id, make(chan string))

	c.Printf("ok, album is now enabled: %s\n", a.Id)
}

func DisableAlbum(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing album name"))
		return
	}
	name := c.Args[0]

	a := core.Node.Datastore.Albums().GetAlbumByName(name)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find album: %s", name)))
		return
	}

	if core.Node.LeftRoomChs[a.Id] == nil {
		c.Printf("album already disabled: %s\n", a.Id)
		return
	}

	core.Node.LeaveRoom(a.Id)
	<-core.Node.LeftRoomChs[a.Id]

	c.Printf("ok, album is now disabled: %s\n", a.Id)
}

func AlbumMnemonic(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing album name"))
		return
	}
	name := c.Args[0]

	a := core.Node.Datastore.Albums().GetAlbumByName(name)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find album: %s", name)))
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	c.Println(green(a.Mnemonic))
}
