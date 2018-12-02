package cmd

import (
	"errors"

	"github.com/textileio/textile-go/repo"
)

var errMissingCafeId = errors.New("missing cafe id")

func init() {
	register(&cafesCmd{})
}

type cafesCmd struct {
	Add      addCafesCmd          `command:"add" description:"Register with a cafe"`
	List     lsCafesCmd           `command:"ls" description:"List cafes"`
	Get      getCafesCmd          `command:"get" description:"Get a cafe"`
	Remove   rmCafesCmd           `command:"rm" description:"Remove a cafe"`
	Messages checkCafeMessagesCmd `command:"messages" description:"Checks cafe messages"`
}

func (x *cafesCmd) Name() string {
	return "cafes"
}

func (x *cafesCmd) Short() string {
	return "Manage cafes"
}

func (x *cafesCmd) Long() string {
	return `
Cafes are other peers on the network who offer pinning, backup, and inbox services. 
Use this command to add, list, get, remove cafes and check messages.
`
}

type addCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *addCafesCmd) Usage() string {
	return ` 

Registers with a cafe and saves an expiring service session token.`
}

func (x *addCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	var info *repo.CafeSession
	res, err := executeJsonCmd(POST, "cafes", params{args: args}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
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
	var list []repo.CafeSession
	res, err := executeJsonCmd(GET, "cafes", params{}, &list)
	if err != nil {
		return err
	}
	output(res, nil)
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
	var info *repo.CafeSession
	res, err := executeJsonCmd(GET, "cafes/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
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
	output(res, nil)
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
	output(res, nil)
	return nil
}
