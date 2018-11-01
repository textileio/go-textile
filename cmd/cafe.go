package cmd

import (
	"errors"
	"github.com/textileio/textile-go/repo"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingCafeId = errors.New("missing cafe id")

func init() {
	register(&cafesCmd{})
}

type cafesCmd struct {
	Add       addCafesCmd       `command:"add"`
	List      lsCafesCmd        `command:"ls"`
	Get       getCafesCmd       `command:"get"`
	Remove    rmCafesCmd        `command:"rm"`
	CheckMail checkMailCafesCmd `command:"check-mail"`
}

func (x *cafesCmd) Name() string {
	return "cafes"
}

func (x *cafesCmd) Short() string {
	return "Manage cafes"
}

func (x *cafesCmd) Long() string {
	return `
Cafes are other peers on the network who offer pinning, backup, and inbox services. 
Use this command to add, list, get, and remove cafes.
`
}

func (x *cafesCmd) Shell() *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
	}
	cmd.AddCmd((&addCafesCmd{}).Shell())
	cmd.AddCmd((&lsCafesCmd{}).Shell())
	cmd.AddCmd((&getCafesCmd{}).Shell())
	cmd.AddCmd((&rmCafesCmd{}).Shell())
	return cmd
}

type addCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *addCafesCmd) Name() string {
	return "add"
}

func (x *addCafesCmd) Short() string {
	return "Register with a cafe"
}

func (x *addCafesCmd) Long() string {
	return "Registers with a cafe and saves an expiring service session token."
}

func (x *addCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callAddCafes(args, nil)
}

func (x *addCafesCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callAddCafes(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callAddCafes(args []string, ctx *ishell.Context) error {
	var info *repo.CafeSession
	res, err := executeJsonCmd(POST, "cafes", params{args: args}, &info)
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}

type lsCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsCafesCmd) Name() string {
	return "ls"
}

func (x *lsCafesCmd) Short() string {
	return "List cafes"
}

func (x *lsCafesCmd) Long() string {
	return "List info about all active cafe sessions."
}

func (x *lsCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callLsCafes(args, nil)
}

func (x *lsCafesCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callLsCafes(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callLsCafes(_ []string, ctx *ishell.Context) error {
	var list *[]repo.CafeSession
	res, err := executeJsonCmd(GET, "cafes", params{}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}

type getCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getCafesCmd) Name() string {
	return "get"
}

func (x *getCafesCmd) Short() string {
	return "Get a cafe"
}

func (x *getCafesCmd) Long() string {
	return "Gets and displays info about a cafe session."
}

func (x *getCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetCafes(args, nil)
}

func (x *getCafesCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callGetCafes(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callGetCafes(args []string, ctx *ishell.Context) error {
	if len(args) == 0 {
		return errMissingCafeId
	}
	var info *repo.CafeSession
	res, err := executeJsonCmd(GET, "cafes/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}

type rmCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmCafesCmd) Name() string {
	return "rm"
}

func (x *rmCafesCmd) Short() string {
	return "Remove a cafe"
}

func (x *rmCafesCmd) Long() string {
	return "Deregisters a cafe (content will expire based on the cafe's service rules)."
}

func (x *rmCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmCafes(args, nil)
}

func (x *rmCafesCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callRmCafes(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callRmCafes(args []string, ctx *ishell.Context) error {
	if len(args) == 0 {
		return errMissingCafeId
	}
	res, err := executeStringCmd(DEL, "cafes/"+args[0], params{})
	if err != nil {
		return nil
	}
	output(res, ctx)
	return nil
}

type checkMailCafesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *checkMailCafesCmd) Name() string {
	return "check-mail"
}

func (x *checkMailCafesCmd) Short() string {
	return "Checks mail at all cafes"
}

func (x *checkMailCafesCmd) Long() string {
	return "Check for mail at all cafes. New messages are downloaded and processed opportunistically."
}

func (x *checkMailCafesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callCheckMailCafes(args, nil)
}

func (x *checkMailCafesCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callCheckMailCafes(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callCheckMailCafes(args []string, ctx *ishell.Context) error {
	res, err := executeStringCmd(POST, "cafes/check_mail", params{})
	if err != nil {
		return nil
	}
	output(res, ctx)
	return nil
}
