package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

var errMissingBlockId = fmt.Errorf("missing block ID")

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

-  MERGE:    3-way merge added.
-  IGNORE:   A block was ignored.
-  FLAG:     A block was flagged.
-  JOIN:     Peer joined.
-  ANNOUNCE: Peer set username / avatar / inbox addresses
-  LEAVE:    Peer left.
-  TEXT:     Text message added.
-  FILES:    File(s) added.
-  COMMENT:  Comment added to another block.
-  LIKE:     Like added to another block.

Use this command to list and get blocks in a thread.`
}

type lsBlocksCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  int           `short:"l" long:"limit" description:"List page size." default:"5"`
	Dots   bool          `short:"d" long:"dots" description:"Return GraphViz dots instead of JSON."`
}

func (x *lsBlocksCmd) Usage() string {
	return `

Paginates blocks in a thread.
Use the --dots option to return GraphViz dots instead of JSON.`
}

func (x *lsBlocksCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"offset": x.Offset,
		"limit":  strconv.Itoa(x.Limit),
		"dots":   strconv.FormatBool(x.Dots),
	}

	if x.Dots {
		return callLsDots(opts)
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
		"dots":   opts["dots"],
	})
}

func callLsDots(opts map[string]string) error {
	if opts["thread"] == "" {
		opts["thread"] = "default"
	}

	var viz pb.BlockViz
	_, err := executeJsonPbCmd(GET, "blocks", params{opts: opts}, &viz)
	if err != nil {
		return err
	}
	if viz.Count > 0 {
		output(viz.Dots)
	}

	if viz.Next == "" {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return callLsDots(map[string]string{
		"thread": opts["thread"],
		"offset": viz.Next,
		"limit":  opts["limit"],
		"dots":   opts["dots"],
	})
}

type getBlocksCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getBlocksCmd) Usage() string {
	return `

Gets a thread block by ID.`
}

func (x *getBlocksCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingBlockId
	}

	_, res, err := callGetBlocks(args[0])
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func callGetBlocks(id string) (*pb.Block, string, error) {
	var block pb.Block
	res, err := executeJsonPbCmd(GET, "blocks/"+id, params{}, &block)
	if err != nil {
		return nil, "", err
	}
	return &block, res, nil
}

func callRmBlocks(args []string) error {
	if len(args) == 0 {
		return errMissingBlockId
	}

	res, err := executeJsonCmd(DEL, "blocks/"+util.TrimQuotes(args[0]), params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
