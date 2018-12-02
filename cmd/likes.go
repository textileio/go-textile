package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingLikeId = errors.New("missing like block ID")

func init() {
	register(&likesCmd{})
}

type likesCmd struct {
	Add    addLikesCmd `command:"add"`
	List   lsLikesCmd  `command:"ls"`
	Get    getLikesCmd `command:"get"`
	Ignore rmLikesCmd  `command:"ignore"`
}

func (x *likesCmd) Name() string {
	return "likes"
}

func (x *likesCmd) Short() string {
	return "Manage thread likes"
}

func (x *likesCmd) Long() string {
	return `
Likes are added as blocks in a thread which target
another block, usually a file(s).
Use this command to add, list, get, and ignore likes.
`
}

func (x *likesCmd) Shell() *ishell.Cmd {
	return nil
}

type addLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *addLikesCmd) Name() string {
	return "add"
}

func (x *addLikesCmd) Short() string {
	return "Add a thread like"
}

func (x *addLikesCmd) Long() string {
	return "Adds a like to a thread block."
}

func (x *addLikesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"block": x.Block,
	}
	return callAddLikes(opts)
}

func (x *addLikesCmd) Shell() *ishell.Cmd {
	return nil
}

func callAddLikes(opts map[string]string) error {
	var info *core.ThreadLikeInfo
	res, err := executeJsonCmd(POST, "blocks/"+opts["block"]+"/likes", params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type lsLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *lsLikesCmd) Name() string {
	return "ls"
}

func (x *lsLikesCmd) Short() string {
	return "List thread likes"
}

func (x *lsLikesCmd) Long() string {
	return "List likes on a thread block."
}

func (x *lsLikesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"block": x.Block,
	}
	return callLsLikes(opts)
}

func (x *lsLikesCmd) Shell() *ishell.Cmd {
	return nil
}

func callLsLikes(opts map[string]string) error {
	var list []core.ThreadLikeInfo
	res, err := executeJsonCmd(GET, "blocks/"+opts["block"]+"/likes", params{}, &list)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type getLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getLikesCmd) Name() string {
	return "get"
}

func (x *getLikesCmd) Short() string {
	return "Get a thread like"
}

func (x *getLikesCmd) Long() string {
	return "Gets a thread like by its block ID."
}

func (x *getLikesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetLikes(args)
}

func (x *getLikesCmd) Shell() *ishell.Cmd {
	return nil
}

func callGetLikes(args []string) error {
	if len(args) == 0 {
		return errMissingLikeId
	}
	var info *core.ThreadLikeInfo
	res, err := executeJsonCmd(GET, "blocks/"+args[0]+"/like", params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type rmLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmLikesCmd) Name() string {
	return "ignore"
}

func (x *rmLikesCmd) Short() string {
	return "Ignore a thread like"
}

func (x *rmLikesCmd) Long() string {
	return `
Ignores a thread like by its block ID.
This adds an "ignore" thread block targeted at the like.
Ignored blocks are by default not returned when listing. 
`
}

func (x *rmLikesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmBlocks(args)
}

func (x *rmLikesCmd) Shell() *ishell.Cmd {
	return nil
}
