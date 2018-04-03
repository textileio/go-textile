package commands

import (
	oldcmds "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/commands"
	core "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/commands"

	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
	"gx/ipfs/QmfAkMSt9Fwzk48QDJecPcwCUjnf2uG7MLnmCGTp4C6ouL/go-ipfs-cmds"
	logging "gx/ipfs/QmRb5jh8z2E8hMGN2tkvs1yHynUanqnZ3UeKwgN1i9P1F8/go-log"
)

var log = logging.Logger("core/commands")

var Root = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:  "Your gateway to the decentralized web, built around IPFS.",
		Synopsis: "textile [--debug=<debug> | -D] [--help=<help>] [-h=<h>] <command> ...",
		Subcommands: `
BASIC COMMANDS
  init          Initialize textile local configuration
  wallet        Manage your wallet

ADVANCED COMMANDS
  start        Start a long-running daemon process

TOOL COMMANDS
  version       Show textile version information
  commands      List all available commands

Use 'textile <command> --help' to learn more about each command.

textile is a wrapper around ipfs. ipfs uses a repository in the local file system. By default, the repo is
located at ~/.ipfs.

EXIT STATUS

The CLI will exit with one of the following values:

0     Successful execution.
1     Failed executions.
`,
	},
	Options: []cmdkit.Option{
		cmdkit.BoolOption("debug", "D", "Operate in debug mode."),
		cmdkit.BoolOption("help", "Show the full command help text."),
		cmdkit.BoolOption("h", "Show a short version of the command help text."),

		// global options, added to every command
		cmds.OptionEncodingType,
		cmds.OptionStreamChannels,
		cmds.OptionTimeout,
	},
}

// commandsDaemonCmd is the "ipfs commands" command for daemon
var CommandsDaemonCmd = CommandsCmd(Root)

var rootSubcommands = map[string]*cmds.Command{
	"wallet":   WalletCmd,
	"commands": CommandsDaemonCmd,
	//"version":   lgc.NewCommand(VersionCmd),
	//"shutdown":  lgc.NewCommand(daemonShutdownCmd),
}

// RootRO is the readonly version of Root
var RootRO = &cmds.Command{}

var CommandsDaemonROCmd = CommandsCmd(RootRO)

var RefsROCmd = &oldcmds.Command{}

var rootROSubcommands = map[string]*cmds.Command{
	"commands": CommandsDaemonROCmd,
	"wallet": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"cat": walletCatCmd,
		},
	},
	//"version": lgc.NewCommand(VersionCmd),
}

func init() {
	Root.ProcessHelp()
	*RootRO = *Root

	// sanitize readonly refs command
	*RefsROCmd = *core.RefsCmd
	RefsROCmd.Subcommands = map[string]*oldcmds.Command{}

	Root.Subcommands = rootSubcommands

	RootRO.Subcommands = rootROSubcommands
}

type MessageOutput struct {
	Message string
}
