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
	Mode   string        `short:"m" long:"mode" description:"Feed mode. One of: chrono, annotated, stacks." default:"chrono"`
}

func (x *lsCmd) Name() string {
	return "ls"
}

func (x *lsCmd) Short() string {
	return "Paginate thread content"
}

func (x *lsCmd) Long() string {
	return `
Paginates post (join|leave|files|message) and annotation (comment|like) block types.
The --mode option dictates how the feed is displayed:

-  "chrono": All feed block types are shown. Annotations always nest their target post, i.e., the post a comment is about.
-  "annotated": Annotations are nested under post targets, but are not shown in the top-level feed.
-  "stacks": Related blocks are chronologically grouped into "stacks". A new stack is started if an unrelated block
   breaks continuity. This mode is used by Textile Photos. Stacks may include:

*  The initial post with some nested annotations. Newer annotations may have already been listed. 
*  One or more annotations about a post. The newest annotation assumes the "top" position in the stack. Additional
     annotations are nested under the target. Newer annotations may have already been listed in the case as well.

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
		"mode":   x.Mode,
	}
	return callLs(opts)
}

func callLs(opts map[string]string) error {
	var list pb.FeedItemList
	res, err := executeJsonPbCmd(GET, "feed", params{opts: opts}, &list)
	if err != nil {
		return err
	}
	if list.Count > 0 {
		output(res)
	}

	if list.Next == "" {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return callLs(map[string]string{
		"thread": opts["thread"],
		"offset": list.Next,
		"limit":  opts["limit"],
		"mode":   opts["mode"],
	})
}
