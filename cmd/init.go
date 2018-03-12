package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/textileio/textile-go/repo"

	cmds "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
	nconfig "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"

	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
)

var initCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Initializes textile config file.",
		ShortDescription: `
Initializes textile configuration files and generates a new keypair.
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.FileArg("default-config", false, false, "Initialize with the given configuration.").EnableStdin(),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("dir", "d", "Repo directory.").WithDefault("~/.ipfs"),
		cmdkit.IntOption("bits", "b", "Number of bits to use in the generated RSA private key.").WithDefault(repo.NBitsForKeypairDefault),
		cmdkit.StringOption("profile", "p", "Apply profile settings to config. Multiple profiles can be separated by ','"),

		// TODO need to decide whether to expose the override as a file or a
		// directory. That is: should we allow the user to also specify the
		// name of the file?
		// TODO cmdkit.StringOption("event-logs", "l", "Location for machine-readable event logs."),
	},
	PreRun: func(req cmds.Request) error {
		daemonLocked, err := fsrepo.LockedByOtherProcess(req.InvocContext().ConfigRoot)
		if err != nil {
			return err
		}

		log.Info("checking if daemon is running...")
		if daemonLocked {
			log.Debug("textile daemon is running")
			e := "textile daemon is running. please stop it to run this command"
			return cmds.ClientError(e)
		}

		return nil
	},
	Run: func(req cmds.Request, res cmds.Response) {

		repoDir, _, err := req.Option("r").String()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		// needs to be called at least once
		res.SetOutput(nil)

		if req.InvocContext().Online {
			res.SetError(errors.New("init must be run offline only!"), cmdkit.ErrNormal)
			return
		}

		nBitsForKeypair, _, err := req.Option("b").Int()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		var conf *nconfig.Config

		f := req.Files()
		if f != nil {
			confFile, err := f.NextFile()
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}

			conf = &nconfig.Config{}
			if err := json.NewDecoder(confFile).Decode(conf); err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
		}

		profile, _, err := req.Option("profile").String()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		var profiles []string
		if profile != "" {
			profiles = strings.Split(profile, ",")
		}

		if err := repo.DoInit(os.Stdout, repoDir, nBitsForKeypair, profiles, conf); err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
	},
}
