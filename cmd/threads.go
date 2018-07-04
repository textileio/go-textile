package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet/thread"
	"gopkg.in/abiosoft/ishell.v2"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

func ListThreads(c *ishell.Context) {
	threads := core.Node.Wallet.Threads()
	if len(threads) == 0 {
		c.Println("no threads found")
	} else {
		c.Println(fmt.Sprintf("found %v threads", len(threads)))
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	for _, thrd := range threads {
		c.Println(blue(fmt.Sprintf("name: %s, id: %s", thrd.Name, thrd.Id)))
	}
}

func AddThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	c.Print("key pair mnemonic phrase (optional): ")
	mnemonics := c.ReadLine()
	var mnemonic *string
	if mnemonics != "" {
		mnemonic = &mnemonics
	}

	thrd, mnem, err := core.Node.Wallet.AddThreadWithMnemonic(name, mnemonic)
	if err != nil {
		c.Err(err)
		return
	}

	go Subscribe(c, thrd)

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("added thread '%s'", name)))
	c.Println(cyan(fmt.Sprintf("mnemonic phrase: %s", mnem)))
}

func PublishThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	_, thrd := core.Node.Wallet.GetThreadByName(name)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	head, err := thrd.GetHead()
	if err != nil {
		c.Err(err)
		return
	}
	if head == "" {
		c.Println(blue("nothing to publish"))
		return
	}
	peers := thrd.Peers("", -1)
	if len(peers) == 0 {
		c.Println(blue("no peers to publish to"))
		return
	}

	err = thrd.PostHead()
	if err != nil {
		c.Err(err)
		return
	}

	c.Println(blue(fmt.Sprintf("published %s in thread %s to %d peers", head, thrd.Name, len(peers))))
}

func ListThreadPeers(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	_, thrd := core.Node.Wallet.GetThreadByName(name)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

	peers := thrd.Peers("", -1)
	if len(peers) == 0 {
		c.Println(fmt.Sprintf("no peers found in: %s", name))
	} else {
		c.Println(fmt.Sprintf("found %v peers in: %s", len(peers), name))
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	for _, peer := range peers {
		c.Println(green(peer.Id))
	}
}

func AddThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer pub key"))
		return
	}
	pks := c.Args[0]
	if len(c.Args) == 1 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[1]

	_, thrd := core.Node.Wallet.GetThreadByName(name)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", name)))
		return
	}

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

	if _, err := thrd.AddInvite(pk); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("invite sent!"))
}

func RemoveThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	err := core.Node.Wallet.RemoveThread(name)
	if err != nil {
		c.Err(err)
		return
	}

	red := color.New(color.FgHiRed).SprintFunc()
	c.Println(red(fmt.Sprintf("removed thread '%s'", name)))
}

func Subscribe(shell ishell.Actions, thrd *thread.Thread) {
	cyan := color.New(color.FgCyan).SprintFunc()
	for {
		select {
		case update, ok := <-thrd.Updates():
			if !ok {
				return
			}
			shell.Printf(cyan(fmt.Sprintf("new block %s in thread %s", update.Id, update.ThreadName)))
		}
	}
}
