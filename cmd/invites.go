package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

func InviteCreate(threadID string, address string, wait int) error {
	if threadID == "" {
		threadID = "default"
	}

	if address != "" {
		contact, _, _ := GetContact(address)
		if contact != nil {
			return CreateInvite(threadID, address)
		}

		output("Could not find contact locally, searching network...")

		results := handleSearchStream("contacts/search", params{
			opts: map[string]string{
				"address": address,
				"limit":   strconv.Itoa(10),
				"wait":    strconv.Itoa(wait),
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
		result, ok := remote[address]
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

		res, err := executeStringCmd(http.MethodPut, "contacts/"+contact.Address, params{
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

	return CreateInvite(threadID, address)
}

func CreateInvite(threadID string, address string) error {
	res, err := executeJsonCmd(http.MethodPost, "invites", params{
		opts: map[string]string{
			"thread":  threadID,
			"address": address,
		},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func InviteList() error {
	res, err := executeJsonCmd(http.MethodGet, "invites", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func InviteAccept(inviteID string, key string) error {
	res, err := executeJsonCmd(http.MethodPost, "invites/"+util.TrimQuotes(inviteID)+"/accept", params{
		args: []string{inviteID},
		opts: map[string]string{
			"key": key,
		},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}


func InviteIgnore(inviteID string) error {
	res, err := executeStringCmd(http.MethodPost, "invites/"+util.TrimQuotes(inviteID)+"/ignore", params{
		args: []string{inviteID},
	})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
