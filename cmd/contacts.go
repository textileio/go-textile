package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
)

var errMissingAddInfo = errors.New("missing peer id or account address")
var errMissingPeerId = errors.New("missing peer id")

func init() {
	register(&contactsCmd{})
}

type contactsCmd struct {
	Add    addContactsCmd    `command:"add" description:"Add a contact"`
	Ls     lsContactsCmd     `command:"ls" description:"List known contacts"`
	Get    getContactsCmd    `command:"get" description:"Get a known contact"`
	Remove rmContactsCmd     `command:"rm" description:"Remove a known contact"`
	Search searchContactsCmd `command:"search" description:"Find contacts"`
}

func (x *contactsCmd) Name() string {
	return "contacts"
}

func (x *contactsCmd) Short() string {
	return "Manage contacts"
}

func (x *contactsCmd) Long() string {
	return `
Use this command to add, list, get, and remove local contacts and find other contacts on the network.
`
}

type addContactsCmd struct {
	Client   ClientOptions `group:"Client Options"`
	Username string        `short:"u" long:"username" description:"Add by username."`
	Peer     string        `short:"p" long:"peer" description:"Add by peer ID."`
	Address  string        `short:"a" long:"address" description:"Add by account address."`
	Wait     int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 10s)." default:"2"`
}

func (x *addContactsCmd) Usage() string {
	return `

Adds a contact by username, peer ID, or account address to known contacts.
`
}

func (x *addContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Username == "" && x.Peer == "" && x.Address == "" {
		return errMissingAddInfo
	}

	var limit int
	if x.Peer != "" {
		limit = 1
	} else if x.Username != "" || x.Address != "" {
		limit = 10
	}

	results := handleSearchStream("contacts/search", params{
		opts: map[string]string{
			"username": x.Username,
			"peer":     x.Peer,
			"address":  x.Address,
			"limit":    strconv.Itoa(limit),
			"wait":     strconv.Itoa(x.Wait),
		},
	})

	if len(results) == 0 {
		output("No contacts were found")
		return nil
	}

	var remote []pb.QueryResult
	for _, res := range results {
		if !res.Local {
			remote = append(remote, res)
		}
	}
	if len(remote) == 0 {
		output("No new contacts were found")
		return nil
	}

	var postfix string
	if len(remote) > 1 {
		postfix = "s"
	}
	if !confirm(fmt.Sprintf("Add %d contact%s?", len(remote), postfix)) {
		return nil
	}

	for _, result := range remote {
		contact := new(pb.Contact)
		if err := ptypes.UnmarshalAny(result.Value, contact); err != nil {
			return err
		}
		data, err := pbMarshaler.MarshalToString(result.Value)
		if err != nil {
			return err
		}

		res, err := executeStringCmd(PUT, "contacts/"+contact.Id, params{
			payload: strings.NewReader(data),
			ctype:   "application/json",
		})
		if err != nil {
			return err
		}
		if res == "ok" {
			output("added " + result.Id)
		} else {
			output("error adding " + result.Id + ": " + res)
		}
	}

	return nil
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

	res, err := executeJsonCmd(GET, "contacts", params{
		opts: map[string]string{
			"thread": x.Thread,
		},
	}, nil)
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

	res, err := executeJsonCmd(GET, "contacts/"+args[0], params{}, nil)
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

type searchContactsCmd struct {
	Client   ClientOptions `group:"Client Options"`
	Username string        `short:"u" long:"username" description:"Search by username."`
	Peer     string        `short:"p" long:"peer" description:"Search by peer ID."`
	Address  string        `short:"a" long:"address" description:"Search by account address."`
	Local    bool          `long:"local" description:"Only search local contacts."`
	Limit    int           `long:"limit" description:"Stops searching after limit results are found." default:"5"`
	Wait     int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 10s)." default:"2"`
}

func (x *searchContactsCmd) Usage() string {
	return `

Searches locally and on the network for contacts.
`
}

func (x *searchContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Username == "" && x.Peer == "" && x.Address == "" {
		return errMissingSearchInfo
	}

	handleSearchStream("contacts/search", params{
		opts: map[string]string{
			"username": x.Username,
			"peer":     x.Peer,
			"address":  x.Address,
			"local":    strconv.FormatBool(x.Local),
			"limit":    strconv.Itoa(x.Limit),
			"wait":     strconv.Itoa(x.Wait),
		},
	})
	return nil
}
