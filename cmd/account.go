package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

var errMissingSnapshotId = errors.New("missing snapshot ID")

func init() {
	register(&accountCmd{})
}

type accountCmd struct {
	Address accountAddressCmd `command:"address" description:"Show wallet account address"`
	Seed    accountSeedCmd    `command:"seed" description:"Show wallet account seed"`
	Contact accountContactCmd `command:"contact" description:"Show own contact"`
	Sync    accountSyncCmd    `command:"sync" description:"Sync account with all network snapshots"`
}

func (x *accountCmd) Name() string {
	return "account"
}

func (x *accountCmd) Short() string {
	return "Manage a wallet account"
}

func (x *accountCmd) Long() string {
	return `
Use this command to manage a wallet account.
`
}

type accountAddressCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *accountAddressCmd) Usage() string {
	return `

Shows the local wallet account address.`
}

func (x *accountAddressCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(GET, "account/address", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type accountSeedCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *accountSeedCmd) Usage() string {
	return `

Shows the local wallet account seed.`
}

func (x *accountSeedCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(GET, "account/seed", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type accountContactCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *accountContactCmd) Usage() string {
	return `

Shows own contact.`
}

func (x *accountContactCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(GET, "account/contact", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type accountSyncCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 30s)." default:"2"`
}

func (x *accountSyncCmd) Usage() string {
	return `

Syncs the local wallet account with all thread snapshots found on the network.
`
}

func (x *accountSyncCmd) Execute(args []string) error {
	setApi(x.Client)

	results := handleSearchStream("snapshots/search", params{
		opts: map[string]string{
			"wait": strconv.Itoa(x.Wait),
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
		if err := applySnapshot(&result); err != nil {
			return err
		}
	}

	if _, err := callCreateSnapshotsThreads(); err != nil {
		return err
	}

	return nil
}
