package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingCommentBody = errors.New("missing comment body")
var errMissingCommentId = errors.New("missing comment block ID")

func init() {
	register(&commentsCmd{})
}

type commentsCmd struct {
	Add    addCommentsCmd `command:"add"`
	List   lsCommentsCmd  `command:"ls"`
	Get    getCommentsCmd `command:"get"`
	Ignore rmCommentsCmd  `command:"ignore"`
}

func (x *commentsCmd) Name() string {
	return "comments"
}

func (x *commentsCmd) Short() string {
	return "Manage thread comments"
}

func (x *commentsCmd) Long() string {
	return `
Comments are added as blocks in a thread, which target
another block, usually a file(s).
Use this command to add, list, get, and ignore comments.
`
}

func (x *commentsCmd) Shell() *ishell.Cmd {
	return nil
}

type addCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *addCommentsCmd) Name() string {
	return "add"
}

func (x *addCommentsCmd) Short() string {
	return "Add a thread comment"
}

func (x *addCommentsCmd) Long() string {
	return "Adds a comment to a thread block."
}

func (x *addCommentsCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"block": x.Block,
	}
	return callAddComments(args, opts)
}

func (x *addCommentsCmd) Shell() *ishell.Cmd {
	return nil
}

func callAddComments(args []string, opts map[string]string) error {
	if len(args) == 0 {
		return errMissingCommentBody
	}

	var info *core.ThreadCommentInfo
	res, err := executeJsonCmd(POST, "blocks/"+opts["block"]+"/comments", params{
		args: args,
	}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type lsCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *lsCommentsCmd) Name() string {
	return "ls"
}

func (x *lsCommentsCmd) Short() string {
	return "List thread comments"
}

func (x *lsCommentsCmd) Long() string {
	return "List comments on a thread block."
}

func (x *lsCommentsCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"block": x.Block,
	}
	return callLsComments(opts)
}

func (x *lsCommentsCmd) Shell() *ishell.Cmd {
	return nil
}

func callLsComments(opts map[string]string) error {
	var list []core.ThreadCommentInfo
	res, err := executeJsonCmd(GET, "blocks/"+opts["block"]+"/comments", params{}, &list)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type getCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getCommentsCmd) Name() string {
	return "get"
}

func (x *getCommentsCmd) Short() string {
	return "Get a thread comment"
}

func (x *getCommentsCmd) Long() string {
	return "Gets a thread comment by its block ID."
}

func (x *getCommentsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetComments(args)
}

func (x *getCommentsCmd) Shell() *ishell.Cmd {
	return nil
}

func callGetComments(args []string) error {
	if len(args) == 0 {
		return errMissingCommentId
	}
	var info *core.ThreadCommentInfo
	res, err := executeJsonCmd(GET, "blocks/"+args[0]+"/comment", params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type rmCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmCommentsCmd) Name() string {
	return "ignore"
}

func (x *rmCommentsCmd) Short() string {
	return "Ignore a thread comment"
}

func (x *rmCommentsCmd) Long() string {
	return `
Ignores a thread comment by its block ID.
This adds an "ignore" thread block targeted at the comment.
Ignored blocks are by default not returned when listing. 
`
}

func (x *rmCommentsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmBlocks(args)
}

func (x *rmCommentsCmd) Shell() *ishell.Cmd {
	return nil
}
