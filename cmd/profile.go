package cmd

import (
	"net/http"

	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

func ProfileGet() error {
	res, _, err := getProfile()
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func getProfile() (string, *pb.Peer, error) {
	var profile pb.Peer
	res, err := executeJsonPbCmd(http.MethodGet, "profile", params{}, &profile)
	if err != nil {
		return "", nil, err
	}
	return res, &profile, err
}

func ProfileSet(name string, avatar string) error {
	if name != "" {
		res, err := executeStringCmd(http.MethodPost, "profile/name", params{args: []string{name}})
		if err != nil {
			return err
		}
		output(res)
	}

	if avatar != "" {
		_, contact, err := getAccount()
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

		if err := FileAdd(avatar, id.Pretty(), "avatar", false, false); err != nil {
			return err
		}

		res, err := executeStringCmd(http.MethodPost, "profile/avatar", params{})
		if err != nil {
			return err
		}
		output(res)
	}

	return nil
}
