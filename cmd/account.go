package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

func getAccount() (string, *pb.Contact, error) {
	var contact pb.Contact
	res, err := executeJsonPbCmd(http.MethodGet, "account", params{}, &contact)
	if err != nil {
		return "", nil, err
	}
	return res, &contact, err
}

func AccountGet() error {
	res, _, err := getAccount()
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func AccountSeed() error {
	res, err := executeStringCmd(http.MethodGet, "account/seed", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func AccountAddress() error {
	res, err := executeStringCmd(http.MethodGet, "account/address", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func AccountSync(wait int) error {
	results := handleSearchStream("snapshots/search", params{
		opts: map[string]string{
			"wait": strconv.Itoa(wait),
		},
	})

	var remote []pb.QueryResult
	for _, res := range results {
		if !res.Local {
			remote = append(remote, res)
		}
	}
	if len(remote) == 0 {
		output("No snapshots were found")
		return nil
	}

	var postfix string
	if len(remote) > 1 {
		postfix = "s"
	}
	if !confirm(fmt.Sprintf("Apply %d snapshot%s?", len(remote), postfix)) {
		return nil
	}

	for _, result := range remote {
		if err := applyThreadSnapshot(&result); err != nil {
			return err
		}
	}

	return nil
}
