package cmd

import (
	"errors"

	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingUsername = errors.New("missing username")
var errMissingAvatar = errors.New("missing avatar file image hash")

func init() {
	register(&profileCmd{})
}

type profileCmd struct {
	Get getProfileCmd `command:"get"`
	Set setProfileCmd `command:"set"`
}

func (x *profileCmd) Name() string {
	return "profile"
}

func (x *profileCmd) Short() string {
	return "Manage public profile"
}

func (x *profileCmd) Long() string {
	return `
Every node has a public IPNS-based profile. 
Use this command to get and set profile username and avatar.
A Textile Account will have different profiles for each of its nodes,
i.e., mobile, desktop, etc.
`
}

func (x *profileCmd) Shell() *ishell.Cmd {
	return nil
}

type getProfileCmd struct {
	Client ClientOptions `group:"Client Options"`
	Peer   string        `short:"p" long:"peer" description:"Fetch a remote peer's public profile'."`
}

func (x *getProfileCmd) Name() string {
	return "get"
}

func (x *getProfileCmd) Short() string {
	return "Get profile"
}

func (x *getProfileCmd) Long() string {
	return "Gets the local node's public IPNS-based profile."
}

func (x *getProfileCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetProfile()
}

func (x *getProfileCmd) Shell() *ishell.Cmd {
	return nil
}

func callGetProfile() error {
	var profile *core.Profile
	res, err := executeJsonCmd(GET, "profile", params{}, &profile)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type setProfileCmd struct {
	Username setUsernameCmd `command:"username"`
	Avatar   setAvatarCmd   `command:"avatar"`
}

func (x *setProfileCmd) Name() string {
	return "set"
}

func (x *setProfileCmd) Short() string {
	return "Set profile fields"
}

func (x *setProfileCmd) Long() string {
	return "Sets public profile username and avatar."
}

func (x *setProfileCmd) Shell() *ishell.Cmd {
	return nil
}

type setUsernameCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *setUsernameCmd) Name() string {
	return "username"
}

func (x *setUsernameCmd) Short() string {
	return "Set username"
}

func (x *setUsernameCmd) Long() string {
	return "Sets public profile username."
}

func (x *setUsernameCmd) Execute(args []string) error {
	setApi(x.Client)
	return callSetUsername(args)
}

func (x *setUsernameCmd) Shell() *ishell.Cmd {
	return nil
}

func callSetUsername(args []string) error {
	if len(args) == 0 {
		return errMissingUsername
	}
	res, err := executeStringCmd(POST, "profile/username", params{args: args})
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type setAvatarCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *setAvatarCmd) Name() string {
	return "avatar"
}

func (x *setAvatarCmd) Short() string {
	return "Set avatar"
}

func (x *setAvatarCmd) Long() string {
	return "Sets public profile avatar by specifying an existing image file hash."
}

func (x *setAvatarCmd) Execute(args []string) error {
	setApi(x.Client)
	return callSetAvatar(args)
}

func (x *setAvatarCmd) Shell() *ishell.Cmd {
	return nil
}

func callSetAvatar(args []string) error {
	if len(args) == 0 {
		return errMissingAvatar
	}
	res, err := executeStringCmd(POST, "profile/avatar", params{args: args})
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
