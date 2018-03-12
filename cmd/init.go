package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	config "github.com/textileio/textile-go/repo/config"

	cmds "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
	core "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	namesys "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/namesys"
	nconfig "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
	fsrepo "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"

	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	"time"
	"bytes"
)

const (
	nBitsForKeypairDefault = 2048
)

type Photo map[string]string

type WalletData struct {
	Photos []Photo `json:"photos"`
}

type Wallet struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Data WalletData `json:"data"`
}

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
		cmdkit.IntOption("bits", "b", "Number of bits to use in the generated RSA private key.").WithDefault(nBitsForKeypairDefault),
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

		if err := doInit(os.Stdout, repoDir, nBitsForKeypair, profiles, conf); err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
	},
}

var errRepoExists = errors.New(`textile configuration file already exists!
Reinitializing would overwrite your keys.
`)

func initWithDefaults(out io.Writer, repoRoot string) error {
	return doInit(out, repoRoot, nBitsForKeypairDefault, nil, nil)
}

func doInit(out io.Writer, repoRoot string, nBitsForKeypair int, confProfiles []string, conf *nconfig.Config) error {
	if _, err := fmt.Fprintf(out, "initializing Textile node at %s\n", repoRoot); err != nil {
		return err
	}

	if err := checkWriteable(repoRoot); err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoRoot) {
		return errRepoExists
	}

	if conf == nil {
		var err error
		conf, err = config.Init(out, nBitsForKeypair)
		if err != nil {
			return err
		}
	}

	for _, profile := range confProfiles {
		transformer, ok := nconfig.Profiles[profile]
		if !ok {
			return fmt.Errorf("invalid configuration profile: %s", profile)
		}

		if err := transformer(conf); err != nil {
			return err
		}
	}

	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return err
	}

	if err := addDefaultAssets(out, repoRoot); err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoRoot)
}

func checkWriteable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := path.Join(dir, "test")
		fi, err := os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		fi.Close()
		return os.Remove(testfile)
	}

	if os.IsNotExist(err) {
		// dir doesnt exist, check that we can create it
		return os.Mkdir(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}

	return err
}

func addDefaultAssets(out io.Writer, repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	api := coreapi.NewCoreAPI(nd)

	wallet := &Wallet{
		Created: time.Now(),
		Updated: time.Now(),
		Data: WalletData{
			Photos: make([]Photo, 0),
		},
	}

	wb, err := json.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("init: encode init wallet failed: %s", err)
	}

	c, err := api.Dag().Put(ctx, bytes.NewReader(wb))
	if err != nil {
		return fmt.Errorf("init: seeding init wallet failed: %s", err)
	}

	if err := api.Pin().Add(nd.Context(), c); err != nil {
		return fmt.Errorf("init: pinning on init wallet failed: %s", err)
	}

	_, err = api.Name().Publish(nd.Context(), c)
	if err != nil {
		return fmt.Errorf("init: publish wallet address failed: %s", err)
	}

	hash := c.Cid().String()
	log.Debugf("init: seeded init wallet %s", hash)

	if _, err = fmt.Fprint(out, "to view your new wallet, enter:\n"); err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "\n\ttextile dag get %s\n\n", hash)
	return err
}

func initializeIpnsKeyspace(repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	err = nd.SetupOfflineRouting()
	if err != nil {
		return err
	}

	return namesys.InitializeKeyspace(ctx, nd.Namesys, nd.Pinning, nd.PrivateKey)
}
