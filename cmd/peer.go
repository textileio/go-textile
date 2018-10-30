package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func init() {
	register(&peerCmd{})
}

type peerCmd struct{}

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
	res, err := executeStringCmd(GET, x.Name(), params{})
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (x *peerCmd) Shell() *ishell.Cmd {
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
