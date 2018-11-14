package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"

	"gopkg.in/abiosoft/ishell.v2"
)

func init() {
	register(&inviteCmd{})
}

var errMissingInviteId = errors.New("missing invite id")

type inviteCmd struct {
	Create createInviteCmd `command:"create"`
	Accept acceptInviteCmd `command:"accept"`
	//Ignore ignoreInviteCmd `command:"ignore"`
}

func (x *inviteCmd) Name() string {
	return "invite"
}

func (x *inviteCmd) Short() string {
	return "Manage thread invites"
}

func (x *inviteCmd) Long() string {
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

func (x *inviteCmd) Shell() *ishell.Cmd {
	return nil
}

type createInviteCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Peer   string        `short:"p" long:"peer" description:"Peer ID. Omit to create an external invite."`
}

func (x *createInviteCmd) Name() string {
	return "create"
}

func (x *createInviteCmd) Short() string {
	return "Create peer-to-peer or external invites to a thread"
}

func (x *createInviteCmd) Long() string {
	return `
Creates a direct peer-to-peer or external invite to a thread.
Omit the --peer option to create an external invite.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *createInviteCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"peer":   x.Peer,
	}
	return callCreateInvite(opts)
}

func (x *createInviteCmd) Shell() *ishell.Cmd {
	return nil
}

func callCreateInvite(opts map[string]string) error {
	threadId := opts["thread"]
	if threadId == "" {
		threadId = "default"
	}
	var result map[string]string
	res, err := executeJsonCmd(POST, "invite", params{
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

type acceptInviteCmd struct {
	Client ClientOptions `group:"Client Options"`
	Key    string        `short:"k" long:"key" description:"Key for an external invite."`
}

func (x *acceptInviteCmd) Name() string {
	return "accept"
}

func (x *acceptInviteCmd) Short() string {
	return "Accept an invite to a thread"
}

func (x *acceptInviteCmd) Long() string {
	return `
Accepts a direct peer-to-peer or external invite to a thread.
Use the --key option with an external invite.
`
}

func (x *acceptInviteCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"key": x.Key,
	}
	return callAcceptInvite(args, opts)
}

func (x *acceptInviteCmd) Shell() *ishell.Cmd {
	return nil
}

func callAcceptInvite(args []string, opts map[string]string) error {
	if len(args) == 0 {
		return errMissingInviteId
	}
	var info core.BlockInfo
	res, err := executeJsonCmd(POST, "invite/"+args[0]+"/accept", params{
		args: args,
		opts: opts,
	}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
