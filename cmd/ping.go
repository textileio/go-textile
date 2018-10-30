package cmd

import (
	"errors"
	"fmt"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
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
	res, err := executeStringCmd(GET, x.Name(), params{args: args})
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (x *pingCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing peer id"))
				return
			}
			pid, err := peer.IDB58Decode(c.Args[0])
			if err != nil {
				c.Println(fmt.Errorf("bad peer id: %s", err))
				return
			}
			status, err := core.Node.Ping(pid)
			if err != nil {
				c.Println(fmt.Errorf("ping failed: %s", err))
				return
			}
			c.Println(status)
		},
	}
}
