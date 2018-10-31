package cmd

import (
	"gopkg.in/abiosoft/ishell.v2"
)

func init() {
	register(&addressCmd{})
}

type addressCmd struct{}

func (x *addressCmd) Name() string {
	return "address"
}
func (x *addressCmd) Short() string {
	return "Show wallet address"
}
func (x *addressCmd) Long() string {
	return "Shows the local node's wallet address."
}

func (x *addressCmd) Execute(args []string) error {
	return callAddress(args, nil)
}

func (x *addressCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name: x.Name(),
		Help: x.Short(),
		Func: func(c *ishell.Context) {
			if err := callAddress(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callAddress(_ []string, ctx *ishell.Context) error {
	res, err := executeStringCmd(GET, "address", params{})
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}
