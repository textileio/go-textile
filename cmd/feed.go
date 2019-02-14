package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/textileio/textile-go/pb"
)

func init() {
	register(&lsCmd{})
}

type lsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  int           `short:"l" long:"limit" description:"List page size." default:"5"`
	Type   string        `long:"type" description:"Feed type. One of: flat, annotated, hybrid (default: annotated)"`
}

func (x *lsCmd) Name() string {
	return "ls"
}

func (x *lsCmd) Short() string {
	return "Paginate thread content"
}

func (x *lsCmd) Long() string {
	return `
Paginates top-level (joins, leaves, files, and messages) and annotation (comments and likes) block types.
The --type option dictates how the feed is displayed.

-  FLAT: All feed types are listed. Annotation types include their top-level target, e.g., thing thing a comment is about.
-  ANNOTATED: Annotation types are nested under top-level targets.
-  HYBRID: Annotation types are nested under top-level targets. However, if the top-level target changes, a subsequent
   annotation is shown on its own and includes its top-level target, which nests additional annotations.  

Omit the --thread option to paginate all files.
Specify "default" to use the default thread (if selected).
`
}

func (x *lsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Type == "" {
		x.Type = "annotated"
	}

	opts := map[string]string{
		"thread": x.Thread,
		"offset": x.Offset,
		"limit":  strconv.Itoa(x.Limit),
		"type":   x.Type,
	}
	return callLs(opts)
}

func callLs(opts map[string]string) error {
	var list pb.FeedItemList
	res, err := executeJsonPbCmd(GET, "feed", params{opts: opts}, &list)
	if err != nil {
		return err
	}

	output(res)

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

	return callLs(map[string]string{
		"thread":    opts["thread"],
		"offset":    list.Items[len(list.Items)-1].Block,
		"limit":     opts["limit"],
		"annotated": opts["annotated"],
	})
}
