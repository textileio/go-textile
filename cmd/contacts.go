package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
)

var errMissingAddInfo = fmt.Errorf("missing name or account address")

func ContactAdd(name string, address string, wait int) error {
	if name == "" && address == "" {
		return errMissingAddInfo
	}

	results := handleSearchStream("contacts/search", params{
		opts: map[string]string{
			"name":    name,
			"address": address,
			"limit":   strconv.Itoa(10),
			"wait":    strconv.Itoa(wait),
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
			output("error adding " + result.Id + ": " + res)
		}
	}

	return nil
}

func ContactList() error {
	res, err := executeJsonCmd(http.MethodGet, "contacts", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ContactGet(address string) error {
	_, res, err := getContact(address)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func getContact(address string) (*pb.Contact, string, error) {
	var contact pb.Contact
	res, err := executeJsonPbCmd(http.MethodGet, "contacts/"+address, params{}, &contact)
	if err != nil {
		return nil, "", err
	}
	return &contact, res, nil
}

func ContactDelete(address string) error {
	res, err := executeStringCmd(http.MethodDelete, "contacts/"+address, params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ContactSearch(name string, address string, local bool, remote bool, limit int, wait int) error {
	if name == "" && address == "" {
		return errMissingSearchInfo
	}

	handleSearchStream("contacts/search", params{
		opts: map[string]string{
			"name":    name,
			"address": address,
			"local":   strconv.FormatBool(local),
			"remote":  strconv.FormatBool(remote),
			"limit":   strconv.Itoa(limit),
			"wait":    strconv.Itoa(wait),
		},
	})

	return nil
}
