package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/textileio/textile-go/repo"

	cmds "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/commands"
	nconfig "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/config"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo"

	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
)

var initCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Initializes textile config file.",
		ShortDescription: `
Initializes textile configuration files and generates a new keypair.
`,
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
		// needs to be called at least once
		res.SetOutput(nil)

		if req.InvocContext().Online {
			res.SetError(errors.New("init must be run offline only!"), cmdkit.ErrNormal)
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

		if err := repo.DoInit(os.Stdout, req.InvocContext().ConfigRoot, nil); err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
	},
}
