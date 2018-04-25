package cmd

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func ListAlbums(c *ishell.Context) {
	rooms := core.Node.IpfsNode.Floodsub.GetTopics()
	albums := core.Node.Datastore.Albums().GetAlbums("")

	if len(albums) == 0 {
		c.Println("no threads found")
	} else {
		c.Println(fmt.Sprintf("found %v threads", len(albums)))
	}

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
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	c.Print("key pair mnemonic phrase (optional): ")
	mnemonic := c.ReadLine()

	if err := core.Node.CreateAlbum(mnemonic, name); err != nil {
		c.Err(err)
		return
	}

	a := core.Node.Datastore.Albums().GetAlbumByName(name)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	go core.Node.JoinRoom(a.Id, make(chan string))

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("created thread #%s", name)))
}

func EnableAlbum(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	a := core.Node.Datastore.Albums().GetAlbumByName(name)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	if core.Node.LeftRoomChs[a.Id] != nil {
		c.Printf("already enabled: %s\n", a.Id)
		return
	}

	go core.Node.JoinRoom(a.Id, make(chan string))

	c.Printf("ok, now enabled: %s\n", a.Id)
}

func DisableAlbum(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	a := core.Node.Datastore.Albums().GetAlbumByName(name)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	if core.Node.LeftRoomChs[a.Id] == nil {
		c.Printf("already disabled: %s\n", a.Id)
		return
	}

	core.Node.LeaveRoom(a.Id)
	<-core.Node.LeftRoomChs[a.Id]

	c.Printf("ok, now disabled: %s\n", a.Id)
}

func AlbumMnemonic(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	a := core.Node.Datastore.Albums().GetAlbumByName(name)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	c.Println(green(a.Mnemonic))
}
