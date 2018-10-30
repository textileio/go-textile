package cmd

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

func init() {
	register(&threadsCmd{})
}

type threadsCmd struct {
	Add  addThreadsCmd `command:"add"`
	List lsThreadsCmd  `command:"ls"`
	Get  getThreadsCmd `command:"get"`
	//Delete delThreadsCmd `command:"del"`
}

func (x *threadsCmd) Name() string {
	return "threads"
}

func (x *threadsCmd) Short() string {
	return "Manage threads"
}

func (x *threadsCmd) Long() string {
	return "Add, ls, get, and del threads."
}

func (x *threadsCmd) Shell() *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
	}
	cmd.AddCmd((&addThreadsCmd{}).Shell())
	cmd.AddCmd((&lsThreadsCmd{}).Shell())
	cmd.AddCmd((&getThreadsCmd{}).Shell())
	//cmd.AddCmd((&delThreadsCmd{}).Shell())
	return cmd
}

type addThreadsCmd struct{}

func (x *addThreadsCmd) Name() string {
	return "add"
}

func (x *addThreadsCmd) Short() string {
	return "Add a new thread"
}

func (x *addThreadsCmd) Long() string {
	return "Adds a new thread for tracking a set of files between peers."
}

func (x *addThreadsCmd) Execute(args []string) error {
	res, err := executeStringCmd(POST, "threads/"+x.Name(), params{args: args})
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (x *addThreadsCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing thread name"))
				return
			}
			sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
			if err != nil {
				c.Err(err)
				return
			}
			thrd, err := core.Node.AddThread(c.Args[0], sk, true)
			if err != nil {
				c.Err(err)
				return
			}
			c.Println(Grey("id:  ") + Cyan(thrd.Id))
		},
	}
}

type lsThreadsCmd struct{}

func (x *lsThreadsCmd) Name() string {
	return "ls"
}

func (x *lsThreadsCmd) Short() string {
	return "List threads"
}

func (x *lsThreadsCmd) Long() string {
	return "Adds a new thread for tracking a set of files between peers."
}

func (x *lsThreadsCmd) Execute(args []string) error {
	var list *struct {
		Items []core.ThreadInfo `json:"items"`
	}
	if err := executeJsonCmd(GET, "threads", params{}, &list); err != nil {
		return err
	}
	jsonb, err := json.MarshalIndent(list.Items, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonb))
	return nil
}

func (x *lsThreadsCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			var infos []*core.ThreadInfo
			for _, thrd := range core.Node.Threads() {
				info, err := thrd.Info()
				if err != nil {
					c.Err(err)
					return
				}
				infos = append(infos, info)
			}
			if len(infos) == 0 {
				c.Println(Grey("[]"))
				return
			}
			c.Println("[")
			for i, info := range infos {
				jsonb, err := json.MarshalIndent(info, "", "    ")
				if err != nil {
					c.Err(err)
					return
				}
				c.Println(Grey(string(jsonb)))
				if i != len(infos)-1 {
					c.Print(Grey(","))
				}
			}
			c.Println("]")
		},
	}
}

type getThreadsCmd struct{}

func (x *getThreadsCmd) Name() string {
	return "get"
}

func (x *getThreadsCmd) Short() string {
	return "Get a thread"
}

func (x *getThreadsCmd) Long() string {
	return "Gets and displays info for a thread."
}

func (x *getThreadsCmd) Execute(args []string) error {
	var info *core.ThreadInfo
	if err := executeJsonCmd(GET, "threads/"+x.Name(), params{args: args}, &info); err != nil {
		return err
	}
	jsonb, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonb))
	return nil
}

func (x *getThreadsCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing thread id"))
				return
			}
			_, thrd := core.Node.Thread(c.Args[0])
			if thrd == nil {
				c.Err(errors.New(fmt.Sprintf("could not find thread: %s", c.Args[0])))
				return
			}
			info, err := thrd.Info()
			if thrd == nil {
				c.Err(err)
				return
			}
			jsonb, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				c.Err(err)
				return
			}
			c.Println(Grey(string(jsonb)))
		},
	}
}

//////////////////////////////////////////////////////////////////////////////

func listThreadPeers(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	id := c.Args[0]

	_, thrd := core.Node.Thread(id)
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

func listThreadBlocks(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	threadId := c.Args[0]

	_, thrd := core.Node.Thread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	blocks := core.Node.Blocks("", -1, "threadId='"+thrd.Id+"'")
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

func ignoreBlock(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]

	block, err := core.Node.Block(id)
	if err != nil {
		c.Err(err)
		return
	}
	_, thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	if _, err := thrd.Ignore(block.Id); err != nil {
		c.Err(err)
		return
	}
}

func addThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing peer id"))
		return
	}
	peerId := c.Args[0]
	if len(c.Args) == 1 {
		c.Err(errors.New("missing thread id"))
		return
	}
	threadId := c.Args[1]

	_, thrd := core.Node.Thread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		c.Err(err)
		return
	}

	if _, err := thrd.AddInvite(pid); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("invite sent!"))
}

func acceptThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing invite address"))
		return
	}
	blockId := c.Args[0]

	if _, err := core.Node.AcceptThreadInvite(blockId); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("ok, accepted"))
}

func addExternalThreadInvite(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	id := c.Args[0]

	_, thrd := core.Node.Thread(id)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", id)))
		return
	}

	hash, key, err := thrd.AddExternalInvite()
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("added! creds: %s %s", hash.B58String(), string(key))))
}

func acceptExternalThreadInvite(c *ishell.Context) {
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

	if _, err := core.Node.AcceptExternalThreadInvite(id, []byte(key)); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("ok, accepted"))
}

func removeThread(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	id := c.Args[0]

	if _, err := core.Node.RemoveThread(id); err != nil {
		c.Err(err)
		return
	}

	red := color.New(color.FgHiRed).SprintFunc()
	c.Println(red("removed thread %s", id))
}
