package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/textileio/textile-go/core"
)

func init() {
	register(&lsCmd{})
}

type lsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  int           `short:"l" long:"limit" description:"List page size." default:"5"`
}

func (x *lsCmd) Name() string {
	return "ls"
}

func (x *lsCmd) Short() string {
	return "Paginate thread content"
}

func (x *lsCmd) Long() string {
	return `
Paginates thread joins, leaves, files, and messages.
Omit the --thread option to paginate all files.
Specify "default" to use the default thread (if selected).
`
}

func (x *lsCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"offset": x.Offset,
		"limit":  strconv.Itoa(x.Limit),
	}
	return callLs(opts)
}

func callLs(opts map[string]string) error {
	var list []core.ThreadFeedItem
	res, err := executeJsonCmd(GET, "feed", params{opts: opts}, &list)
	if err != nil {
		return err
	}

	output(res)

	limit, err := strconv.Atoi(opts["limit"])
	if err != nil {
		return err
	}
	if len(list) < limit {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return callLs(map[string]string{
		"thread": opts["thread"],
		"offset": list[len(list)-1].Block,
		"limit":  opts["limit"],
	})
}
