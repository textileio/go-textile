package cmd

import (
	"fmt"

	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

var errMissingName = fmt.Errorf("missing name")

func init() {
	register(&profileCmd{})
}

type profileCmd struct {
	Get getProfileCmd `command:"get" description:"Get profile"`
	Set setProfileCmd `command:"set" description:"Set profile name and avatar"`
}

func (x *profileCmd) Name() string {
	return "profile"
}

func (x *profileCmd) Short() string {
	return "Manage public profile"
}

func (x *profileCmd) Long() string {
	return `
Use this command to view and update the peer profile. A Textile account will
show a profile for each of its peers, e.g., mobile, desktop, etc.
`
}

type getProfileCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getProfileCmd) Usage() string {
	return `

Gets the local peer profile.`
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

func callGetProfile() (string, *pb.Peer, error) {
	var profile pb.Peer
	res, err := executeJsonPbCmd(GET, "profile", params{}, &profile)
	if err != nil {
		return "", nil, err
	}
	return res, &profile, err
}

type setProfileCmd struct {
	Name   setNameCmd   `command:"name" description:"Set display name"`
	Avatar setAvatarCmd `command:"avatar" description:"Set avatar"`
}

func (x *setProfileCmd) Usage() string {
	return `

Sets the peer display name and avatar.`
}

type setNameCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *setNameCmd) Usage() string {
	return `

Sets the peer display name.`
}

func (x *setNameCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingName
	}
	res, err := executeStringCmd(POST, "profile/name", params{args: args})
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

Sets the peer avatar from an image path (JPEG, PNG, or GIF).`
}

func (x *setAvatarCmd) Execute(args []string) error {
	setApi(x.Client)

	_, contact, err := callGetAccount()
	if err != nil {
		return err
	}
	kp, err := keypair.Parse(contact.Address)
	if err != nil {
		return err
	}
	id, err := kp.Id()
	if err != nil {
		return err
	}

	opts := map[string]string{
		"thread": id.Pretty(),
	}
	if err := callAddFiles(args, opts); err != nil {
		return err
	}

	res, err := executeStringCmd(POST, "profile/avatar", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
