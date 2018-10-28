package cmd

import (
	"errors"
	"fmt"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

var pingCmd PingCmd

type PingCmd struct{}

func (x *PingCmd) Name() string {
	return "ping"
}

func (x *PingCmd) Short() string {
	return "fixme"
}

func (x *PingCmd) Long() string {
	return "fixme"
}

func (x *PingCmd) Execute(args []string) error {
	return executeStringCmd(x.Name(), args)
}

func (x *PingCmd) Shell() *ishell.Cmd {
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
