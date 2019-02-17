package cmd

import "strconv"

func init() {
	register(&accountCmd{})
}

type accountCmd struct {
	Address addressCmd `command:"address" description:"Show wallet account address"`
	Backups backupsCmd `command:"backups" description:"Find account thread backups"`
}

func (x *accountCmd) Name() string {
	return "account"
}

func (x *accountCmd) Short() string {
	return "Manage a wallet account"
}

func (x *accountCmd) Long() string {
	return `
Use this command to view account address and backups and view and sync with account peers.
`
}

type addressCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *addressCmd) Usage() string {
	return "Shows the local node's wallet account address."
}

func (x *addressCmd) Execute(args []string) error {
	setApi(x.Client)
	res, err := executeStringCmd(GET, "account/address", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type backupsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 10s)." default:"5"`
}

func (x *backupsCmd) Usage() string {
	return `

Finds account thread backups on the network.
`
}

func (x *backupsCmd) Execute(args []string) error {
	setApi(x.Client)

	handleSearchStream("account/backups", params{
		opts: map[string]string{
			"wait": strconv.Itoa(x.Wait),
		},
	})
	return nil
}
