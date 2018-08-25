package cmd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
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

	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		c.Err(err)
		return
	}

	thrd, err := core.Node.Wallet.AddThread(name, sk, true)
	if err != nil {
		c.Err(err)
		return
	}

	if err := core.Node.Wallet.InviteDevices(thrd); err != nil {
		c.Err(err)
		return
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("added thread %s with name %s", thrd.Id, name)))
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
		c.Println(fmt.Sprintf("%v peers:", len(peers)))
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	for _, p := range peers {
		c.Println(green(p.Id))
	}
}

func ListThreadBlocks(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	threadId := c.Args[0]

	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	blocks := thrd.Blocks("", -1, nil, nil)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no blocks found in: %s", threadId))
	} else {
		c.Println(fmt.Sprintf("%v blocks:", len(blocks)))
	}

	magenta := color.New(color.FgHiMagenta).SprintFunc()
	for _, block := range blocks {
		c.Println(magenta(fmt.Sprintf("%s %s", block.Type.Description(), block.Id)))
	}
}

func GetThreadHead(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	threadId := c.Args[0]

	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	head, err := thrd.GetHead()
	if thrd == nil {
		c.Err(err)
		return
	}

	yellow := color.New(color.FgHiYellow).SprintFunc()
	c.Println(yellow(head))
}

func IgnoreBlock(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]

	block, err := core.Node.Wallet.GetBlock(id)
	if err != nil {
		c.Err(err)
		return
	}
	_, thrd := core.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	if _, err := thrd.Ignore(block.Id); err != nil {
		c.Err(err)
		return
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

	if _, err := thrd.AddInvite(pk); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("invite sent!"))
}

func AcceptThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing invite address"))
		return
	}
	blockId := c.Args[0]

	if _, err := core.Node.Wallet.AcceptThreadInvite(blockId); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("ok, accepted"))
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
	c.Println(green(fmt.Sprintf("added! creds: %s %s", addr.B58String(), string(key))))
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

	if _, err := core.Node.Wallet.AcceptExternalThreadInvite(id, []byte(key)); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("ok, accepted"))
}

func RemoveThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	id := c.Args[0]

	if _, err := core.Node.Wallet.RemoveThread(id); err != nil {
		c.Err(err)
		return
	}

	red := color.New(color.FgHiRed).SprintFunc()
	c.Println(red("removed thread %s", id))
}
