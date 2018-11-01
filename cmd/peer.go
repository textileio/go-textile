package cmd

import (
	"gopkg.in/abiosoft/ishell.v2"
)

func init() {
	register(&peerCmd{})
}

type peerCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *peerCmd) Name() string {
	return "peer"
}

func (x *peerCmd) Short() string {
	return "Show peer ID"
}

func (x *peerCmd) Long() string {
	return "Shows the local node's peer ID."
}

func (x *peerCmd) Execute(args []string) error {
	setApi(x.Client)
	return callPeer(args, nil)
}

func (x *peerCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callPeer(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callPeer(_ []string, ctx *ishell.Context) error {
	res, err := executeStringCmd(GET, "peer", params{})
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}
