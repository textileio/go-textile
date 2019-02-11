package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

var errMissingStdin = errors.New("missing stdin")
var errInvalidContact = errors.New("invalid contact format")
var errMissingPeerId = errors.New("missing peer id")

func init() {
	register(&contactsCmd{})
}

type contactsCmd struct {
	Add    addContactsCmd  `command:"add" description:"Add a contact"`
	Ls     lsContactsCmd   `command:"ls" description:"List known contacts"`
	Get    getContactsCmd  `command:"get" description:"Get a known contact"`
	Remove rmContactsCmd   `command:"rm" description:"Remove a known contact"`
	Find   findContactsCmd `command:"find" description:"Find contacts"`
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
	Client ClientOptions `group:"Client Options"`
}

func (x *addContactsCmd) Usage() string {
	return `

Add to known contacts.

NOTE: This command only accepts input from stdin.
A common workflow is to pipe 'textile contacts find' into 'textile contacts add',
just be sure you know what the results of the find are before adding.
`
}

func (x *addContactsCmd) Execute(args []string) error {
	setApi(x.Client)

	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		return errMissingStdin
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var body []byte
	var contact *repo.Contact
	if err := json.Unmarshal(input, &contact); err != nil {
		return errInvalidContact
	}
	if contact.Address != "" {
		body = input
	} else {
		var result pb.QueryResult
		if err := pbUnmarshaler.Unmarshal(bytes.NewReader(input), &result); err != nil {
			return errInvalidContact
		}
		if result.Value == nil {
			return errInvalidContact
		}
		data, err := pbMarshaler.MarshalToString(result.Value)
		if err != nil {
			return err
		}
		body = []byte(data)
	}

	res, err := executeStringCmd(POST, "contacts", params{
		payload: bytes.NewReader(body),
		ctype:   "application/json",
	})
	if err != nil {
		return err
	}
	output(res)
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
	Wait     int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 10s)." default:"5"`
}

func (x *findContactsCmd) Usage() string {
	return `

Finds contacts known locally and on the network.
`
}

func (x *findContactsCmd) Execute(args []string) error {
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
