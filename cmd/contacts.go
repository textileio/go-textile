package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"
)

var errMissingPeerId = errors.New("missing peer id")
var errMissingPeerAddress = errors.New("missing peer address")
var errMissingSearchInfo = errors.New("missing search info")

func init() {
	register(&contactsCmd{})
}

type contactsCmd struct {
	Ls   lsContactsCmd   `command:"ls" description:"List known contacts"`
	Get  getContactsCmd  `command:"get" description:"Get contact information"`
	Add  addContactsCmd  `command:"add" description:"Add a new contact"`
	Find findContactsCmd `command:"find" description:"Find a contact on the network"`
}

func (x *contactsCmd) Name() string {
	return "contacts"
}

func (x *contactsCmd) Short() string {
	return "Get, add, and list local contacts"
}

func (x *contactsCmd) Long() string {
	return "Get, add, and list local contacts."
}

type lsContactsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all known contacts."`
}

func (x *lsContactsCmd) Usage() string {
	return `

List known contacts.

Include the --thread flag to list contacts for a given thread.`
}

func (x *lsContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	var list []core.ContactInfo
	res, err := executeJsonCmd(GET, "contacts", params{
		opts: map[string]string{
			"thread": x.Thread,
		},
	}, &list)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type getContactsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getContactsCmd) Usage() string {
	return `

Get contact information.`
}

func (x *getContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingPeerId
	}
	var info core.ContactInfo
	res, err := executeJsonCmd(GET, "contacts/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type addContactsCmd struct {
	Client   ClientOptions `group:"Client Options"`
	Username string        `short:"u" long:"username" description:"New contact's username. Omit to use peer id."`
}

func (x *addContactsCmd) Usage() string {
	return `

Add a new contact.

Use the --username flag to specify a username.`
}

func (x *addContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) < 2 {
		return errMissingPeerAddress
	} else if len(args) < 1 {
		return errMissingPeerId
	}
	var info core.ContactInfo
	res, err := executeJsonCmd(POST, "contacts", params{
		args: args,
		opts: map[string]string{
			"username": x.Username,
		},
	}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type findContactsCmd struct {
	Client   ClientOptions `group:"Client Options"`
	Username string        `short:"u" long:"username" description:"A username to use in the search."`
	Peer     string        `short:"p" long:"peer" description:"A Peer ID use in the search."`
	Address  string        `short:"a" long:"address" description:"An account address to use in the search."`
}

func (x *findContactsCmd) Usage() string {
	return `

Find contact information.`
}

func (x *findContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Username == "" && x.Peer == "" && x.Address == "" {
		return errMissingSearchInfo
	}
	var infos core.FindContactInfoResult
	res, err := executeJsonCmd(POST, "contacts/search", params{
		opts: map[string]string{
			"username": x.Username,
			"peer":     x.Peer,
			"address":  x.Address,
		},
	}, &infos)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
