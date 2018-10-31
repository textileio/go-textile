package cmd

import (
	"gopkg.in/abiosoft/ishell.v2"
)

func init() {
	register(&pingCmd{})
}

type pingCmd struct{}

func (x *pingCmd) Name() string {
	return "ping"
}

func (x *pingCmd) Short() string {
	return "Ping another peer"
}

func (x *pingCmd) Long() string {
	return "Pings another peer on the network, returning online|offline."
}

func (x *pingCmd) Execute(args []string) error {
	return callPing(args, nil)
}

func (x *pingCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callPing(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callPing(args []string, ctx *ishell.Context) error {
	res, err := executeStringCmd(GET, "ping", params{args: args})
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}
