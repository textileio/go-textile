package cmd

import (
	"fmt"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

var errMissingSnapshotId = fmt.Errorf("missing snapshot ID")

func init() {
	register(&accountCmd{})
}

type accountCmd struct {
	Get     accountGetCmd     `command:"get" description:"Show account contact"`
	Seed    accountSeedCmd    `command:"seed" description:"Show wallet account seed"`
	Address accountAddressCmd `command:"address" description:"Show wallet account address"`
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
Use this command to manage a wallet account.`
}

type accountGetCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *accountGetCmd) Usage() string {
	return `

Shows the local peer's account info as a contact.`
}

func (x *accountGetCmd) Execute(args []string) error {
	setApi(x.Client)

	res, _, err := callGetAccount()
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func callGetAccount() (string, *pb.Contact, error) {
	var contact pb.Contact
	res, err := executeJsonPbCmd(GET, "account", params{}, &contact)
	if err != nil {
		return "", nil, err
	}
	return res, &contact, err
}

type accountSeedCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *accountSeedCmd) Usage() string {
	return `

Shows the local peer's account seed.`
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

type accountAddressCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *accountAddressCmd) Usage() string {
	return `

Shows the local peer's account address.`
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

type accountSyncCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 30s)." default:"2"`
}

func (x *accountSyncCmd) Usage() string {
	return `

Syncs the local account peer with other peers found on the network.`
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
