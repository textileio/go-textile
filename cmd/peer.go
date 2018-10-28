package cmd

import (
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

var peerCmd PeerCmd

type PeerCmd struct{}

func (x *PeerCmd) Name() string {
	return "peer"
}

func (x *PeerCmd) Short() string {
	return "fixme"
}

func (x *PeerCmd) Long() string {
	return "fixme"
}

func (x *PeerCmd) Execute(args []string) error {
	return executeStringCmd(x.Name(), nil)
}

func (x *PeerCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			green := color.New(color.FgHiGreen).SprintFunc()
			pid, err := core.Node.PeerId()
			if err != nil {
				c.Err(err)
				return
			}
			c.Println(green(pid.Pretty()))
		},
	}
}
