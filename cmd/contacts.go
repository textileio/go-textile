package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"
)

var errMissingPeerId = errors.New("missing peer id")

func init() {
	register(&contactsCmd{})
}

type contactsCmd struct {
	Ls  lsContactsCmd  `command:"ls" description:"List known contacts"`
	Get getContactsCmd `command:"get" description:"Get contact information"`
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
	return `List known contacts.
	
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
	return `Get contact information.`
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
