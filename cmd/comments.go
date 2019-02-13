package cmd

import (
	"errors"
)

var errMissingCommentBody = errors.New("missing comment body")
var errMissingCommentId = errors.New("missing comment block ID")

func init() {
	register(&commentsCmd{})
}

type commentsCmd struct {
	Add    addCommentsCmd `command:"add" description:"Add a thread comment"`
	List   lsCommentsCmd  `command:"ls" description:"List thread comments"`
	Get    getCommentsCmd `command:"get" description:"Get a thread comment"`
	Ignore rmCommentsCmd  `command:"ignore" description:"Ignore a thread comment"`
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

type addCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *addCommentsCmd) Usage() string {
	return `

Adds a comment to a thread block.`
}

func (x *addCommentsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingCommentBody
	}

	res, err := executeJsonCmd(POST, "blocks/"+x.Block+"/comments", params{args: args}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type lsCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *lsCommentsCmd) Usage() string {
	return `

Lists comments on a thread block.`
}

func (x *lsCommentsCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(GET, "blocks/"+x.Block+"/comments", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type getCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getCommentsCmd) Usage() string {
	return `

Gets a thread comment by block ID.`
}

func (x *getCommentsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingCommentId
	}

	res, err := executeJsonCmd(GET, "blocks/"+args[0]+"/comment", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type rmCommentsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmCommentsCmd) Usage() string {
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
