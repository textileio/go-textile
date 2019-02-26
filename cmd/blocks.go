package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/textileio/textile-go/pb"
)

var errMissingBlockId = errors.New("missing block ID")

func init() {
	register(&blocksCmd{})
}

type blocksCmd struct {
	List lsBlocksCmd  `command:"ls" description:"Paginate thread blocks"`
	Get  getBlocksCmd `command:"get" description:"Get a thread block"`
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

There are several block types:

-  JOIN:     Peer joined.
-  ANNOUNCE: Peer set username / inbox address
-  LEAVE:    Peer left.
-  FILES:    File(s) added.
-  MESSAGE:  Text message added.
-  COMMENT:  Comment added to another block.
-  LIKE:     Like added to another block.
-  MERGE:    3-way merge added.
-  IGNORE:   Another block was ignored.
-  FLAG:     A flag was added to another block.
  
Use this command to get and list blocks in a thread.
`
}

type lsBlocksCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  int           `short:"l" long:"limit" description:"List page size." default:"5"`
}

func (x *lsBlocksCmd) Usage() string {
	return `

Paginates blocks in a thread.
`
}

func (x *lsBlocksCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"offset": x.Offset,
		"limit":  strconv.Itoa(x.Limit),
	}
	return callLsBlocks(opts)
}

func callLsBlocks(opts map[string]string) error {
	if opts["thread"] == "" {
		opts["thread"] = "default"
	}

	var list pb.BlockList
	res, err := executeJsonPbCmd(GET, "blocks", params{opts: opts}, &list)
	if err != nil {
		return err
	}
	if len(list.Items) > 0 {
		output(res)
	}

	limit, err := strconv.Atoi(opts["limit"])
	if err != nil {
		return err
	}
	if len(list.Items) < limit {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return callLsBlocks(map[string]string{
		"thread": opts["thread"],
		"offset": list.Items[len(list.Items)-1].Id,
		"limit":  opts["limit"],
	})
}

type getBlocksCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getBlocksCmd) Usage() string {
	return `

Gets a thread block by ID.
`
}

func (x *getBlocksCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingBlockId
	}

	res, err := executeJsonCmd(GET, "blocks/"+args[0], params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func callRmBlocks(args []string) error {
	if len(args) == 0 {
		return errMissingBlockId
	}

	res, err := executeJsonCmd(DEL, "blocks/"+args[0], params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
