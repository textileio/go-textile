package cmd

import (
	"errors"
	"strconv"

	"github.com/textileio/textile-go/core"
)

var errMissingPeerId = errors.New("missing peer id")
var errMissingSearchInfo = errors.New("missing search info")

func init() {
	register(&contactsCmd{})
}

type contactsCmd struct {
	Ls     lsContactsCmd   `command:"ls" description:"List known contacts"`
	Get    getContactsCmd  `command:"get" description:"Get a known contact"`
	Remove rmContactsCmd   `command:"rm" description:"Remove a known contact"`
	Find   findContactsCmd `command:"find" description:"Find contacts"`
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

Lists known contacts.

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

Gets a known contact.`
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

type rmContactsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmContactsCmd) Usage() string {
	return `

Removes a known contact.`
}

func (x *rmContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingPeerId
	}
	res, err := executeStringCmd(DEL, "contacts/"+args[0], params{})
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
	Local    bool          `long:"local" description:"Only search local contacts."`
	Limit    int           `long:"limit" description:"Stops searching after limit results are found." default:"5"`
	Wait     int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed." default:"5"`
	Add      bool          `long:"add" description:"Add results to local contacts. Not allowed when searching by username."`
}

func (x *findContactsCmd) Usage() string {
	return `

Finds contacts known locally and on the network.
Use the --add option to save remote results as local contacts.
`
}

func (x *findContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Username == "" && x.Peer == "" && x.Address == "" {
		return errMissingSearchInfo
	}
	var infos core.ContactInfoQueryResult
	res, err := executeJsonCmd(POST, "contacts/search", params{
		opts: map[string]string{
			"username": x.Username,
			"peer":     x.Peer,
			"address":  x.Address,
			"local":    strconv.FormatBool(x.Local),
			"limit":    strconv.Itoa(x.Limit),
			"wait":     strconv.Itoa(x.Wait),
			"add":      strconv.FormatBool(x.Add),
		},
	}, &infos)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
