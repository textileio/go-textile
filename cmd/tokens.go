package cmd

import (
	"github.com/textileio/textile-go/pb"
)

// func init() {
// 	register(&tokensCmd{})
// }

type tokensCmd struct {
	Add    addTokensCmd `command:"add" description:"Create a new access token"`
	List   lsTokensCmd  `command:"ls" description:"List available access tokens"`
	Get    getTokensCmd `command:"get" description:"Get a specific access token"`
	Remove rmTokensCmd  `command:"rm" description:"Remove a specific access token"`
}

func (x *tokensCmd) Name() string {
	return "tokens"
}

func (x *tokensCmd) Short() string {
	return "Manage Cafe access tokens"
}

func (x *tokensCmd) Long() string {
	return `
Tokens allow other peers to register with a Cafe peer.
Use this command to add, list, get, and remove tokens required for access to this peer's Cafe.
`
}

type addCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *addTokensCmd) Usage() string {
	return `

Generates an access token and saves a salted and encrypted version for future lookup.
The response contains a base58 encoded version of the random bytes token.
`
}

func (x *addTokensCmd) Execute(args []string) error {
	setApi(x.Client)
	var info *pb.CafeSession
	res, err := executeJsonCmd(POST, "cafes", params{args: args}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type lsCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsCafesCmd) Usage() string {
	return `

List info about all active cafe sessions.`
}

func (x *lsCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	var list []pb.CafeSession
	res, err := executeJsonCmd(GET, "cafes", params{}, &list)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type getCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getCafesCmd) Usage() string {
	return `

Gets and displays info about a cafe session.
`
}

func (x *getCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingCafeId
	}
	var info *pb.CafeSession
	res, err := executeJsonCmd(GET, "cafes/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type rmCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmCafesCmd) Usage() string {
	return "Deregisters a cafe (content will expire based on the cafe's service rules)."
}

func (x *rmCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingCafeId
	}
	res, err := executeStringCmd(DEL, "cafes/"+args[0], params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type checkCafeMessagesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *checkCafeMessagesCmd) Usage() string {
	return `

Check for messages at all cafes. New messages are downloaded and processed opportunistically.
`
}

func (x *checkCafeMessagesCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(POST, "cafes/messages", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
