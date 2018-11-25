package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func init() {
	register(&invitesCmd{})
}

var errMissingInviteId = errors.New("missing invite id")

type invitesCmd struct {
	Create createInvitesCmd `command:"create"`
	List   lsInvitesCmd     `command:"ls"`
	Accept acceptInvitesCmd `command:"accept"`
	Ignore ignoreInvitesCmd `command:"ignore"`
}

func (x *invitesCmd) Name() string {
	return "invites"
}

func (x *invitesCmd) Short() string {
	return "Manage thread invites"
}

func (x *invitesCmd) Long() string {
	return `
Invites allow other peers to join threads. There are two types of
invites: direct peer-to-peer and external.

Peer-to-peer invites are encrypted with the invitee's public key.

External invites are encrypted with a single-use key and are useful for 
onboarding new users. Careful though. Once an external invite and its key are
shared, the thread should be considered public, since any number of peers
can use it to join.
`
}

func (x *invitesCmd) Shell() *ishell.Cmd {
	return nil
}

type createInvitesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Peer   string        `short:"p" long:"peer" description:"Peer ID. Omit to create an external invite."`
}

func (x *createInvitesCmd) Name() string {
	return "create"
}

func (x *createInvitesCmd) Short() string {
	return "Create peer-to-peer or external invites to a thread"
}

func (x *createInvitesCmd) Long() string {
	return `
Creates a direct peer-to-peer or external invite to a thread.
Omit the --peer option to create an external invite.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *createInvitesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"peer":   x.Peer,
	}
	return callCreateInvites(opts)
}

func (x *createInvitesCmd) Shell() *ishell.Cmd {
	return nil
}

func callCreateInvites(opts map[string]string) error {
	threadId := opts["thread"]
	if threadId == "" {
		threadId = "default"
	}
	var result map[string]string
	res, err := executeJsonCmd(POST, "invites", params{
		opts: map[string]string{
			"thread": threadId,
			"peer":   opts["peer"],
		},
	}, &result)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type lsInvitesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsInvitesCmd) Name() string {
	return "ls"
}

func (x *lsInvitesCmd) Short() string {
	return "List thread invites"
}

func (x *lsInvitesCmd) Long() string {
	return "Lists all pending thread invites."
}

func (x *lsInvitesCmd) Execute(_ []string) error {
	setApi(x.Client)
	return callLsInvites()
}

func (x *lsInvitesCmd) Shell() *ishell.Cmd {
	return nil
}

func callLsInvites() error {
	var list []core.ThreadInviteInfo
	res, err := executeJsonCmd(GET, "invites", params{}, &list)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type acceptInvitesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Key    string        `short:"k" long:"key" description:"Key for an external invite."`
}

func (x *acceptInvitesCmd) Name() string {
	return "accept"
}

func (x *acceptInvitesCmd) Short() string {
	return "Accept an invite to a thread"
}

func (x *acceptInvitesCmd) Long() string {
	return `
Accepts a direct peer-to-peer or external invite to a thread.
Use the --key option with an external invite.
`
}

func (x *acceptInvitesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"key": x.Key,
	}
	return callAcceptInvites(args, opts)
}

func (x *acceptInvitesCmd) Shell() *ishell.Cmd {
	return nil
}

func callAcceptInvites(args []string, opts map[string]string) error {
	if len(args) == 0 {
		return errMissingInviteId
	}
	var info core.BlockInfo
	res, err := executeJsonCmd(POST, "invites/"+args[0]+"/accept", params{
		args: args,
		opts: opts,
	}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type ignoreInvitesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *ignoreInvitesCmd) Name() string {
	return "ignore"
}

func (x *ignoreInvitesCmd) Short() string {
	return "Ignore direct invite to a thread"
}

func (x *ignoreInvitesCmd) Long() string {
	return `
Ignores a direct peer-to-peer invite to a thread.
`
}

func (x *ignoreInvitesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callIgnoreInvites(args)
}

func (x *ignoreInvitesCmd) Shell() *ishell.Cmd {
	return nil
}

func callIgnoreInvites(args []string) error {
	if len(args) == 0 {
		return errMissingInviteId
	}
	res, err := executeStringCmd(POST, "invites/"+args[0]+"/ignore", params{
		args: args,
	})
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
