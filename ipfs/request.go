package ipfs

import (
	"context"

	"github.com/op/go-logging"
	cmds "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
	cli "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands/cli"
	commands "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/commands"
	"time"
)

var log = logging.MustGetLogger("ipfs")

var Root = &cmds.Command{
	Options:  commands.Root.Options,
	Helptext: commands.Root.Helptext,
}

var commandsClientCmd = commands.CommandsCmd(Root)

var localCommands = map[string]*cmds.Command{
	"commands": commandsClientCmd,
}

func init() {
	Root.Subcommands = localCommands
	for k, v := range commands.Root.Subcommands {
		if _, found := Root.Subcommands[k]; !found {
			Root.Subcommands[k] = v
		}
	}
}

func NewRequest(ctx cmds.Context, args []string) (cmds.Request, *cmds.Command, error) {
	Root.Subcommands = localCommands
	for k, v := range commands.Root.Subcommands {
		if _, found := Root.Subcommands[k]; !found {
			Root.Subcommands[k] = v
		}
	}
	req, cmd, _, err := cli.Parse(args, nil, Root)
	cctx := context.Background()
	rerr := req.SetRootContext(cctx)
	if rerr != nil {
		return nil, nil, rerr
	}
	req.SetInvocContext(ctx)
	return req, cmd, err
}

func NewRequestWithTimeout(ctx cmds.Context, args []string, timeout time.Duration) (cmds.Request, *cmds.Command, error) {
	Root.Subcommands = localCommands
	for k, v := range commands.Root.Subcommands {
		if _, found := Root.Subcommands[k]; !found {
			Root.Subcommands[k] = v
		}
	}
	req, cmd, _, err := cli.Parse(args, nil, Root)
	cctx, _ := context.WithTimeout(context.Background(), timeout)
	rerr := req.SetRootContext(cctx)
	if rerr != nil {
		return nil, nil, rerr
	}
	req.SetInvocContext(ctx)
	return req, cmd, err
}
