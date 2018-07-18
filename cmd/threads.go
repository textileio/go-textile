package cmd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet/thread"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
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

	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		c.Err(err)
		return
	}

	if _, err := core.Node.Wallet.AddThread(name, sk); err != nil {
		c.Err(err)
		return
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("added thread '%s'", name)))
}

func ListThreadPeers(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	id := c.Args[0]

	_, thrd := core.Node.Wallet.GetThread(id)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", id)))
		return
	}

	peers := thrd.Peers()
	if len(peers) == 0 {
		c.Println(fmt.Sprintf("no peers found in: %s", id))
	} else {
		c.Println(fmt.Sprintf("found %v peers in: %s", len(peers), id))
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	for _, p := range peers {
		c.Println(green(p.Id))
	}
}

func AddThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer pub key"))
		return
	}
	pks := c.Args[0]
	if len(c.Args) == 1 {
		c.Err(errors.New("missing thread id"))
		return
	}
	id := c.Args[1]

	_, thrd := core.Node.Wallet.GetThread(id)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", id)))
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

	addr, err := thrd.AddInvite(pk)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("invite sent! added block %s.", addr.B58String())))
}

func AcceptThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing invite address"))
		return
	}
	blockId := c.Args[0]

	addr, err := core.Node.Wallet.AcceptThreadInvite(blockId)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("ok, accepted. added block %s.", addr.B58String())))
}

func AddExternalThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	id := c.Args[0]

	_, thrd := core.Node.Wallet.GetThread(id)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", id)))
		return
	}

	addr, key, err := thrd.AddExternalInvite()
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("id: %s, key: %s", addr.B58String(), string(key))))
}

func AcceptExternalThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing invite id"))
		return
	}
	id := c.Args[0]
	if len(c.Args) == 1 {
		c.Err(errors.New("missing invite key"))
		return
	}
	key := c.Args[1]

	addr, err := core.Node.Wallet.AcceptExternalThreadInvite(id, []byte(key))
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("ok, accepted. added block %s.", addr.B58String())))
}

func RemoveThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	name := c.Args[0]

	addr, err := core.Node.Wallet.RemoveThread(name)
	if err != nil {
		c.Err(err)
		return
	}

	red := color.New(color.FgHiRed).SprintFunc()
	c.Println(red(fmt.Sprintf("removed thread '%s'. added block %s.", name, addr.B58String())))
}

func Subscribe(thrd *thread.Thread, peerId string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	for {
		select {
		case update, ok := <-thrd.Updates():
			if !ok {
				return
			}
			authorPkb, err := libp2pc.ConfigDecodeKey(update.Index.AuthorPk)
			if err != nil {
				fmt.Printf(red(err.Error()))
				return
			}
			authorPk, err := libp2pc.UnmarshalPublicKey(authorPkb)
			if err != nil {
				fmt.Printf(red(err.Error()))
				return
			}
			authorId, err := peer.IDFromPublicKey(authorPk)
			if err != nil {
				fmt.Printf(red(err.Error()))
				return
			}
			if authorId.Pretty() == peerId {
				return
			}
			fmt.Printf(cyan(fmt.Sprintf("\nnew block %s in thread %s from %s", update.Index.Id, update.ThreadName, authorId.Pretty())))
		}
	}
}
