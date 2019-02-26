package cmd

import (
	"errors"

	"github.com/textileio/textile-go/pb"
)

var errMissingUsername = errors.New("missing username")
var errMissingAvatar = errors.New("missing avatar file image hash")

func init() {
	register(&profileCmd{})
}

type profileCmd struct {
	Get getProfileCmd `command:"get" description:"Get profile"`
	Set setProfileCmd `command:"set" description:"Set profile fields"`
}

func (x *profileCmd) Name() string {
	return "profile"
}

func (x *profileCmd) Short() string {
	return "Manage public profile"
}

func (x *profileCmd) Long() string {
	return `
Every node has a public profile. 
Use this command to get and set profile username and avatar.
A Textile Account will have different profiles for each of its nodes,
i.e., mobile, desktop, etc.
`
}

type getProfileCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getProfileCmd) Usage() string {
	return `

Gets the local node's public profile.`
}

func (x *getProfileCmd) Execute(args []string) error {
	setApi(x.Client)
	res, _, err := callGetProfile()
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func callGetProfile() (string, *pb.Contact, error) {
	var profile pb.Contact
	res, err := executeJsonPbCmd(GET, "profile", params{}, &profile)
	if err != nil {
		return "", nil, err
	}
	return res, &profile, err
}

type setProfileCmd struct {
	Username setUsernameCmd `command:"username" description:"Set username"`
	Avatar   setAvatarCmd   `command:"avatar" description:"Set avatar"`
}

func (x *setProfileCmd) Usage() string {
	return `

Sets public profile username and avatar.`
}

type setUsernameCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *setUsernameCmd) Usage() string {
	return `

Sets public profile username.`
}

func (x *setUsernameCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingUsername
	}
	res, err := executeStringCmd(POST, "profile/username", params{args: args})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type setAvatarCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *setAvatarCmd) Usage() string {
	return `

Sets public profile avatar by specifying an existing image file hash.`
}

func (x *setAvatarCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingAvatar
	}
	res, err := executeStringCmd(POST, "profile/avatar", params{args: args})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
