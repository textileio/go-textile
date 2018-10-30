package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
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
	res, err := executeStringCmd(GET, x.Name(), params{})
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (x *addressCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name: x.Name(),
		Help: x.Short(),
		Func: func(c *ishell.Context) {
			cyan := color.New(color.FgHiCyan).SprintFunc()
			addr, err := core.Node.Address()
			if err != nil {
				c.Err(err)
				return
			}
			c.Println(cyan(addr))
		},
	}
}
