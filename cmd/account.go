package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
)

var errMissingBackupId = errors.New("missing backup ID")

func init() {
	register(&accountCmd{})
}

type accountCmd struct {
	Address accountAddressCmd `command:"address" description:"Show wallet account address"`
	Peers   accountPeersCmd   `command:"peers" description:"List known account peers"`
	Backups accountBackupsCmd `command:"backups" description:"Manage account thread backups"`
	Sync    accountSyncCmd    `command:"sync" description:"Sync account with all network backups"`
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

type accountPeersCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *accountPeersCmd) Usage() string {
	return `

Lists all known wallet account peers.`
}

func (x *accountPeersCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(GET, "account/peers", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type accountBackupsCmd struct {
	List  lsAccountBackupsCmd    `command:"ls" description:"Search for wallet account thread backups"`
	Apply applyAccountBackupsCmd `command:"apply" description:"Apply a single thread backup"`
}

func (x *accountBackupsCmd) Usage() string {
	return `

Use this command to List and apply wallet account backups.`
}

type lsAccountBackupsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 10s)." default:"2"`
}

func (x *lsAccountBackupsCmd) Usage() string {
	return `

Searches the network for wallet account thread backups.
`
}

func (x *lsAccountBackupsCmd) Execute(args []string) error {
	setApi(x.Client)

	handleSearchStream("account/backups", params{
		opts: map[string]string{
			"wait": strconv.Itoa(x.Wait),
		},
	})
	return nil
}

type applyAccountBackupsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 10s)." default:"2"`
}

func (x *applyAccountBackupsCmd) Usage() string {
	return `

Applies a single wallet account thread backup.
`
}

func (x *applyAccountBackupsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingBackupId
	}
	id := args[0]

	results := handleSearchStream("account/backups", params{
		opts: map[string]string{
			"wait": strconv.Itoa(x.Wait),
		},
	})

	var result *pb.QueryResult
	for _, r := range results {
		if r.Id == id {
			result = &r
		}
	}

	if result == nil {
		output("Could not find backup with ID: " + id)
		return nil
	}

	if err := applyBackup(result); err != nil {
		return err
	}

	return nil
}

type accountSyncCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 10s)." default:"2"`
}

func (x *accountSyncCmd) Usage() string {
	return `

Syncs the local wallet account with all thread backups found on the network.
`
}

func (x *accountSyncCmd) Execute(args []string) error {
	setApi(x.Client)

	results := handleSearchStream("account/backups", params{
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
		output("No backups were found")
		return nil
	}

	var postfix string
	if len(remote) > 1 {
		postfix = "s"
	}
	if !confirm(fmt.Sprintf("Apply %d backup%s?", len(remote), postfix)) {
		return nil
	}

	for _, result := range remote {
		if err := applyBackup(&result); err != nil {
			return err
		}
	}

	return nil
}

func applyBackup(result *pb.QueryResult) error {
	backup := new(pb.Thread)
	if err := ptypes.UnmarshalAny(result.Value, backup); err != nil {
		return err
	}
	data, err := pbMarshaler.MarshalToString(result.Value)
	if err != nil {
		return err
	}

	res, err := executeStringCmd(PUT, "threads/"+backup.Id, params{
		payload: strings.NewReader(data),
		ctype:   "application/json",
	})
	if err != nil {
		return err
	}
	if res == "ok" {
		output("applied " + result.Id)
	} else {
		output("error applying " + result.Id + ": " + res)
	}
	return nil
}
