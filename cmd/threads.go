package cmd

import (
	"errors"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingThreadId = errors.New("missing thread id")

func init() {
	register(&threadsCmd{})
}

type threadsCmd struct {
	Add    addThreadsCmd `command:"add"`
	List   lsThreadsCmd  `command:"ls"`
	Get    getThreadsCmd `command:"get"`
	Remove rmThreadsCmd  `command:"rm"`
}

func (x *threadsCmd) Name() string {
	return "threads"
}

func (x *threadsCmd) Short() string {
	return "Manage threads"
}

func (x *threadsCmd) Long() string {
	return `
Threads are distributed sets of encrypted files between peers,
governed by build-in or custom Schemas.
Use this command to add, list, get, and remove threads.

Open threads are the most common thread type. Open threads allow 
any member to invite new members.

Private threads are primarily used internally for backup/recovery 
purpose and 1-to-1 communication channels.
`
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
	cmd.AddCmd((&rmThreadsCmd{}).Shell())
	return cmd
}

type addThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Key    string        `short:"k" long:"key" description:"A locally unique key used by an app to identify this thread on recovery."`
	Type   string        `short:"t" long:"type" description:"Thread type [open, private]." default:"open"`
	Schema string        `short:"s" long:"schema" description:"Thread schema [photos]." default:"photos"`
}

func (x *addThreadsCmd) Name() string {
	return "add"
}

func (x *addThreadsCmd) Short() string {
	return "Add a new thread"
}

func (x *addThreadsCmd) Long() string {
	return "Adds and joins a new thread."
}

func (x *addThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"key":    x.Key,
		"type":   x.Type,
		"schema": x.Schema,
	}
	return callAddThreads(args, opts, nil)
}

func (x *addThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callAddThreads(args []string, opts map[string]string, ctx *ishell.Context) error {
	var info *core.ThreadInfo
	res, err := executeJsonCmd(POST, "threads", params{args: args, opts: opts}, &info)
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}

type lsThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsThreadsCmd) Name() string {
	return "ls"
}

func (x *lsThreadsCmd) Short() string {
	return "List threads"
}

func (x *lsThreadsCmd) Long() string {
	return "List info about all threads."
}

func (x *lsThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callLsThreads(args, nil)
}

func (x *lsThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callLsThreads(_ []string, ctx *ishell.Context) error {
	var list *[]core.ThreadInfo
	res, err := executeJsonCmd(GET, "threads", params{}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}

type getThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getThreadsCmd) Name() string {
	return "get"
}

func (x *getThreadsCmd) Short() string {
	return "Get a thread"
}

func (x *getThreadsCmd) Long() string {
	return "Gets and displays info about a thread."
}

func (x *getThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetThreads(args, nil)
}

func (x *getThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callGetThreads(args []string, ctx *ishell.Context) error {
	if len(args) == 0 {
		return errMissingThreadId
	}
	var info *core.ThreadInfo
	res, err := executeJsonCmd(GET, "threads/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}

type rmThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmThreadsCmd) Name() string {
	return "rm"
}

func (x *rmThreadsCmd) Short() string {
	return "Remove a thread"
}

func (x *rmThreadsCmd) Long() string {
	return "Leaves and removes a thread."
}

func (x *rmThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmThreads(args, nil)
}

func (x *rmThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callRmThreads(args []string, ctx *ishell.Context) error {
	if len(args) == 0 {
		return errMissingThreadId
	}
	res, err := executeStringCmd(DEL, "threads/"+args[0], params{})
	if err != nil {
		return nil
	}
	output(res, ctx)
	return nil
}

//////////////////////////////////////////////////////////////////////////////

//func listThreadPeers(c *ishell.Context) {
//	if len(c.Args) == 0 {
//		c.Err(errors.New("missing thread id"))
//		return
//	}
//	id := c.Args[0]
//
//	thrd := core.Node.Thread(id)
//	if thrd == nil {
//		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", id)))
//		return
//	}
//
//	peers := thrd.Peers()
//	if len(peers) == 0 {
//		c.Println(fmt.Sprintf("no peers found in: %s", id))
//	} else {
//		c.Println(fmt.Sprintf("%v peers:", len(peers)))
//	}
//
//	green := color.New(color.FgHiGreen).SprintFunc()
//	for _, p := range peers {
//		c.Println(green(p.Id))
//	}
//}
//
//func listThreadBlocks(c *ishell.Context) {
//	if len(c.Args) == 0 {
//		c.Err(errors.New("missing thread id"))
//		return
//	}
//	threadId := c.Args[0]
//
//	thrd := core.Node.Thread(threadId)
//	if thrd == nil {
//		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
//		return
//	}
//
//	blocks := core.Node.Blocks("", -1, "threadId='"+thrd.Id+"'")
//	if len(blocks) == 0 {
//		c.Println(fmt.Sprintf("no blocks found in: %s", threadId))
//	} else {
//		c.Println(fmt.Sprintf("%v blocks:", len(blocks)))
//	}
//
//	magenta := color.New(color.FgHiMagenta).SprintFunc()
//	for _, block := range blocks {
//		c.Println(magenta(fmt.Sprintf("%s %s", block.Type.Description(), block.Id)))
//	}
//}
//
//func ignoreBlock(c *ishell.Context) {
//	if len(c.Args) == 0 {
//		c.Err(errors.New("missing block id"))
//		return
//	}
//	id := c.Args[0]
//
//	block, err := core.Node.Block(id)
//	if err != nil {
//		c.Err(err)
//		return
//	}
//	thrd := core.Node.Thread(block.ThreadId)
//	if thrd == nil {
//		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
//		return
//	}
//
//	if _, err := thrd.Ignore(block.Id); err != nil {
//		c.Err(err)
//		return
//	}
//}
//
//func addThreadInvite(c *ishell.Context) {
//	if len(c.Args) == 0 {
//		c.Err(errors.New("missing peer id"))
//		return
//	}
//	peerId := c.Args[0]
//	if len(c.Args) == 1 {
//		c.Err(errors.New("missing thread id"))
//		return
//	}
//	threadId := c.Args[1]
//
//	thrd := core.Node.Thread(threadId)
//	if thrd == nil {
//		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
//		return
//	}
//
//	pid, err := peer.IDB58Decode(peerId)
//	if err != nil {
//		c.Err(err)
//		return
//	}
//
//	if _, err := thrd.AddInvite(pid); err != nil {
//		c.Err(err)
//		return
//	}
//
//	green := color.New(color.FgHiGreen).SprintFunc()
//	c.Println(green("invite sent!"))
//}
//
//func acceptThreadInvite(c *ishell.Context) {
//	if len(c.Args) == 0 {
//		c.Err(errors.New("missing invite address"))
//		return
//	}
//	blockId := c.Args[0]
//
//	if _, err := core.Node.AcceptThreadInvite(blockId); err != nil {
//		c.Err(err)
//		return
//	}
//
//	green := color.New(color.FgHiGreen).SprintFunc()
//	c.Println(green("ok, accepted"))
//}
//
//func addExternalThreadInvite(c *ishell.Context) {
//	if len(c.Args) == 0 {
//		c.Err(errors.New("missing thread id"))
//		return
//	}
//	id := c.Args[0]
//
//	thrd := core.Node.Thread(id)
//	if thrd == nil {
//		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", id)))
//		return
//	}
//
//	hash, key, err := thrd.AddExternalInvite()
//	if err != nil {
//		c.Err(err)
//		return
//	}
//
//	green := color.New(color.FgHiGreen).SprintFunc()
//	c.Println(green(fmt.Sprintf("added! creds: %s %s", hash.B58String(), string(key))))
//}
//
//func acceptExternalThreadInvite(c *ishell.Context) {
//	if len(c.Args) == 0 {
//		c.Err(errors.New("missing invite id"))
//		return
//	}
//	id := c.Args[0]
//	if len(c.Args) == 1 {
//		c.Err(errors.New("missing invite key"))
//		return
//	}
//	key := c.Args[1]
//
//	if _, err := core.Node.AcceptExternalThreadInvite(id, []byte(key)); err != nil {
//		c.Err(err)
//		return
//	}
//
//	green := color.New(color.FgHiGreen).SprintFunc()
//	c.Println(green("ok, accepted"))
//}
