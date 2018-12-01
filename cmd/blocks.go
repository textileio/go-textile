package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingBlockId = errors.New("missing block ID")

func init() {
	register(&blocksCmd{})
}

type blocksCmd struct {
	List lsBlocksCmd  `command:"ls"`
	Get  getBlocksCmd `command:"get"`
}

func (x *blocksCmd) Name() string {
	return "blocks"
}

func (x *blocksCmd) Short() string {
	return "View thread blocks"
}

func (x *blocksCmd) Long() string {
	return `
Blocks are the raw components in a thread. Think of them as an
append-only log of thread updates where each update is hash-linked
to its parent(s). New / recovering peers can sync history by simply
traversing the hash tree.

There are several thread types:

  JOIN:     Peer joined.
  ANNOUNCE: Peer set username / inbox address
  LEAVE:    Peer left.
  FILES:    File(s) added.
  MESSAGE:  Text message added.
  COMMENT:  Comment added to another block.
  LIKE:     Like added to another block.
  MERGE:    3-way merge added.
  IGNORE:   Another block was ignored.
  FLAG:     A flag was added to another block.
  
Use this command to get and list blocks in a thread.
`
}

func (x *blocksCmd) Shell() *ishell.Cmd {
	return nil
}

type lsBlocksCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
}

func (x *lsBlocksCmd) Name() string {
	return "ls"
}

func (x *lsBlocksCmd) Short() string {
	return "List thread blocks"
}

func (x *lsBlocksCmd) Long() string {
	return "List blocks on a thread block."
}

func (x *lsBlocksCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
	}
	return callLsBlocks(opts)
}

func (x *lsBlocksCmd) Shell() *ishell.Cmd {
	return nil
}

func callLsBlocks(opts map[string]string) error {
	if opts["thread"] == "" {
		opts["thread"] = "default"
	}

	var list *[]core.BlockInfo
	res, err := executeJsonCmd(GET, "blocks", params{opts: opts}, &list)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type getBlocksCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getBlocksCmd) Name() string {
	return "get"
}

func (x *getBlocksCmd) Short() string {
	return "Get a thread block"
}

func (x *getBlocksCmd) Long() string {
	return "Gets a thread block by its ID."
}

func (x *getBlocksCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetBlocks(args)
}

func (x *getBlocksCmd) Shell() *ishell.Cmd {
	return nil
}

func callGetBlocks(args []string) error {
	if len(args) == 0 {
		return errMissingBlockId
	}
	var info *core.BlockInfo
	res, err := executeJsonCmd(GET, "blocks/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

func callRmBlocks(args []string) error {
	if len(args) == 0 {
		return errMissingLikeId
	}
	var info *core.BlockInfo
	res, err := executeJsonCmd(DEL, "blocks/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
