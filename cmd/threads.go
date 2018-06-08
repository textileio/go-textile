package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet"
	"gopkg.in/abiosoft/ishell.v2"
	"sort"
)

func ListThreads(c *ishell.Context) {
	threads := core.Node.Wallet.Threads()
	if len(threads) == 0 {
		c.Println("no threads found")
	} else {
		c.Println(fmt.Sprintf("found %v threads", len(threads)))
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	for _, thread := range threads {
		mem := "disabled"
		if thread.Listening() {
			mem = "enabled"
		}
		c.Println(blue(fmt.Sprintf("name: %s, id: %s, status: %s", thread.Name, thread.Id, mem)))
	}
}

func CreateThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	c.Print("key pair mnemonic phrase (optional): ")
	mnemonic := c.ReadLine()

	thread, err := core.Node.Wallet.AddThread(name, mnemonic)
	if err != nil {
		c.Err(err)
		return
	}

	Subscribe(c, thread)

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("created thread #%s", name)))
}

func EnableThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	thread := core.Node.Wallet.GetThreadByName(name)
	if thread == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	if thread.Listening() {
		c.Printf("already enabled: %s\n", thread.Id)
		return
	}

	Subscribe(c, thread)

	c.Printf("ok, now enabled: %s\n", thread.Id)
}

func DisableAlbum(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	thread := core.Node.Wallet.GetThreadByName(name)
	if thread == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	if !thread.Listening() {
		c.Printf("already disabled: %s\n", thread.Id)
		return
	}

	thread.Unsubscribe()
	<-thread.LeftCh

	c.Printf("ok, now disabled: %s\n", thread.Id)
}

func PublishThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	thread := core.Node.Wallet.GetThreadByName(name)
	if thread == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}
	thread.Publish()

	blue := color.New(color.FgHiBlue).SprintFunc()
	c.Println(blue(fmt.Sprintf("published %s to %s", thread.Head, thread.Id)))
}

func ListThreadPeers(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	thread := core.Node.Wallet.GetThreadByName(name)
	if thread == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	peers := core.Node.Wallet.Ipfs.Floodsub.ListPeers(thread.Id)
	var list []string
	for _, peer := range peers {
		list = append(list, peer.Pretty())
	}
	sort.Strings(list)

	if len(list) == 0 {
		c.Println(fmt.Sprintf("no peers found in: %s", name))
	} else {
		c.Println(fmt.Sprintf("found %v peers in: %s", len(list), name))
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	for _, peer := range list {
		c.Println(green(peer))
	}
}

func Subscribe(shell ishell.Actions, thread *wallet.Thread) {
	cyan := color.New(color.FgCyan).SprintFunc()
	datac := make(chan wallet.ThreadUpdate)
	go thread.Subscribe(datac)
	go func() {
		for {
			select {
			case update, ok := <-datac:
				if !ok {
					return
				}
				msg := fmt.Sprintf("\nnew photo %s in %s thread\n", update.Id, update.Thread)
				shell.ShowPrompt(false)
				shell.Printf(cyan(msg))
				shell.ShowPrompt(true)
			}
		}
	}()
}
