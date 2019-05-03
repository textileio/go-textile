package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

var errMissingAddInfo = fmt.Errorf("missing name or account address")
var errMissingAddress = fmt.Errorf("missing account address")

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
Use this command to add, list, get, and remove local contacts and find other contacts on the network.`
}

type addContactsCmd struct {
	Client  ClientOptions `group:"Client Options"`
	Name    string        `short:"n" long:"name" description:"Add by display name."`
	Address string        `short:"a" long:"address" description:"Add by account address."`
	Wait    int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 30s)." default:"2"`
}

func (x *addContactsCmd) Usage() string {
	return `

Adds a contact by display name or account address to known contacts.`
}

func (x *addContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Name == "" && x.Address == "" {
		return errMissingAddInfo
	}

	results := handleSearchStream("contacts/search", params{
		opts: map[string]string{
			"name":    x.Name,
			"address": x.Address,
			"limit":   strconv.Itoa(10),
			"wait":    strconv.Itoa(x.Wait),
		},
	})

	if len(results) == 0 {
		output("No contacts were found")
		return nil
	}

	remote := make(map[string]pb.QueryResult)
	for _, res := range results {
		if !res.Local {
			remote[res.Id] = res // overwrite with newer / more complete result
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

		res, err := executeStringCmd(PUT, "contacts/"+contact.Address, params{
			payload: strings.NewReader(data),
			ctype:   "application/json",
		})
		if err != nil {
			return err
		}
		if res == "" {
			output("added " + result.Id)
		} else {
			output("error adding " + result.Id + ": " + res)
		}
	}

	return nil
}

type lsContactsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsContactsCmd) Usage() string {
	return `

Lists known contacts.`
}

func (x *lsContactsCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(GET, "contacts", params{}, nil)
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
		return errMissingAddress
	}

	_, res, err := callGetContacts(args[0])
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func callGetContacts(address string) (*pb.Contact, string, error) {
	var contact pb.Contact
	res, err := executeJsonPbCmd(GET, "contacts/"+address, params{}, &contact)
	if err != nil {
		return nil, "", err
	}
	return &contact, res, nil
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
		return errMissingAddress
	}

	res, err := executeStringCmd(DEL, "contacts/"+util.TrimQuotes(args[0]), params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type searchContactsCmd struct {
	Client     ClientOptions `group:"Client Options"`
	Name       string        `short:"n" long:"name" description:"Search by display name."`
	Address    string        `short:"a" long:"address" description:"Search by account address."`
	LocalOnly  bool          `long:"only-local" description:"Only search local contacts."`
	RemoteOnly bool          `long:"only-remote" description:"Only search remote contacts."`
	Limit      int           `long:"limit" description:"Stops searching after limit results are found." default:"5"`
	Wait       int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 30s)." default:"2"`
}

func (x *searchContactsCmd) Usage() string {
	return `

Searches locally and on the network for contacts.`
}

func (x *searchContactsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Name == "" && x.Address == "" {
		return errMissingSearchInfo
	}

	handleSearchStream("contacts/search", params{
		opts: map[string]string{
			"name":    x.Name,
			"address": x.Address,
			"local":   strconv.FormatBool(x.LocalOnly),
			"remote":  strconv.FormatBool(x.RemoteOnly),
			"limit":   strconv.Itoa(x.Limit),
			"wait":    strconv.Itoa(x.Wait),
		},
	})
	return nil
}
