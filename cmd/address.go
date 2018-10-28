package cmd

import (
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

var addressCmd AddressCmd

type AddressCmd struct{}

func (x *AddressCmd) Name() string {
	return "address"
}
func (x *AddressCmd) Short() string {
	return "fixme"
}
func (x *AddressCmd) Long() string {
	return "fixme"
}

func (x *AddressCmd) Execute(args []string) error {
	return executeStringCmd(x.Name(), nil)
}

func (x *AddressCmd) Shell() *ishell.Cmd {
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
