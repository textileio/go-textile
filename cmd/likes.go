package cmd

import (
	"errors"
)

var errMissingLikeId = errors.New("missing like block ID")

func init() {
	register(&likesCmd{})
}

type likesCmd struct {
	Add    addLikesCmd `command:"add" description:"Add a thread like"`
	List   lsLikesCmd  `command:"ls" description:"List thread likes"`
	Get    getLikesCmd `command:"get" description:"Get a thread like"`
	Ignore rmLikesCmd  `command:"ignore" description:"Ignore a thread like"`
}

func (x *likesCmd) Name() string {
	return "likes"
}

func (x *likesCmd) Short() string {
	return "Manage thread likes"
}

func (x *likesCmd) Long() string {
	return `
Likes are added as blocks in a thread, which target
another block, usually a file(s).
Use this command to add, list, get, and ignore likes.
`
}

type addLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *addLikesCmd) Usage() string {
	return `

Adds a like to a thread block.`
}

func (x *addLikesCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(POST, "blocks/"+x.Block+"/likes", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type lsLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *lsLikesCmd) Usage() string {
	return `

Lists likes on a thread block.`
}

func (x *lsLikesCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(GET, "blocks/"+x.Block+"/likes", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type getLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getLikesCmd) Usage() string {
	return `

Gets a thread like by block ID.`
}

func (x *getLikesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingLikeId
	}

	res, err := executeJsonCmd(GET, "blocks/"+args[0]+"/like", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type rmLikesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmLikesCmd) Usage() string {
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
