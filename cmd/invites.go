package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

func init() {
	register(&invitesCmd{})
}

var errMissingInviteId = errors.New("missing invite id")

type invitesCmd struct {
	Create createInvitesCmd `command:"create" description:"Create account-to-account or external invites to a thread"`
	List   lsInvitesCmd     `command:"ls" description:"List thread invites"`
	Accept acceptInvitesCmd `command:"accept" description:"Accept an invite to a thread"`
	Ignore ignoreInvitesCmd `command:"ignore" description:"Ignore direct invite to a thread"`
}

func (x *invitesCmd) Name() string {
	return "invites"
}

func (x *invitesCmd) Short() string {
	return "Manage thread invites"
}

func (x *invitesCmd) Long() string {
	return `
Invites allow other users to join threads. There are two types of
invites, direct account-to-account and external:

- Account-to-account invites are encrypted with the invitee's account address (public key).
- External invites are encrypted with a single-use key and are useful for onboarding new users.`
}

type createInvitesCmd struct {
	Client  ClientOptions `group:"Client Options"`
	Thread  string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Address string        `short:"a" long:"address" description:"Account address. Omit to create an external invite."`
	Wait    int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 30s)." default:"2"`
}

func (x *createInvitesCmd) Usage() string {
	return `

Creates a direct account-to-account or external invite to a thread.
Omit the --address option to create an external invite.
Omit the --thread option to use the default thread (if selected).`
}

func (x *createInvitesCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Thread == "" {
		x.Thread = "default"
	}

	if x.Address != "" {
		contact, _, _ := callGetContacts(x.Address)
		if contact != nil {
			return callCreateInvites(x.Thread, x.Address)
		}

		output("Could not find contact locally, searching network...")

		results := handleSearchStream("contacts/search", params{
			opts: map[string]string{
				"address": x.Address,
				"limit":   strconv.Itoa(10),
				"wait":    strconv.Itoa(x.Wait),
			},
		})

		if len(results) == 0 {
			output("Could not find contact")
			return nil
		}

		remote := make(map[string]pb.QueryResult)
		for _, res := range results {
			if !res.Local {
				remote[res.Id] = res // overwrite with newer / more complete result
			}
		}
		result, ok := remote[x.Address]
		if !ok {
			output("Could not find contact")
			return nil
		}

		if !confirm(fmt.Sprintf("Add and invite %s?", result.Id)) {
			return nil
		}

		contact = new(pb.Contact)
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
			return fmt.Errorf("error adding %s: %s", result.Id, res)
		}
	}

	return callCreateInvites(x.Thread, x.Address)
}

func callCreateInvites(thread string, address string) error {
	res, err := executeJsonCmd(POST, "invites", params{
		opts: map[string]string{
			"thread":  thread,
			"address": address,
		},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type lsInvitesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsInvitesCmd) Usage() string {
	return `

Lists all pending thread invites.`
}

func (x *lsInvitesCmd) Execute(_ []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(GET, "invites", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type acceptInvitesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Key    string        `short:"k" long:"key" description:"Key for an external invite."`
}

func (x *acceptInvitesCmd) Usage() string {
	return `

Accepts a direct account-to-account or external invite to a thread.
Use the --key option with an external invite.`
}

func (x *acceptInvitesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingInviteId
	}

	res, err := executeJsonCmd(POST, "invites/"+util.TrimQuotes(args[0])+"/accept", params{
		args: args,
		opts: map[string]string{
			"key": x.Key,
		},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type ignoreInvitesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *ignoreInvitesCmd) Usage() string {
	return `

Ignores a direct account-to-account invite to a thread.`
}

func (x *ignoreInvitesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingInviteId
	}

	res, err := executeStringCmd(POST, "invites/"+util.TrimQuotes(args[0])+"/ignore", params{
		args: args,
	})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
