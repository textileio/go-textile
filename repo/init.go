package repo

import (
	"fmt"
	"os"
	"time"
	"bytes"
	"io"
	"errors"
	"path"
	"context"
	"encoding/json"
	"github.com/textileio/textile-go/repo/config"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"
	nconfig "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/namesys"
)

const (
	NBitsForKeypairDefault = 2048
)

var errRepoExists = errors.New(`textile configuration file already exists!
Reinitializing would overwrite your keys.
`)

func InitWithDefaults(out io.Writer, repoRoot string) error {
	return DoInit(out, repoRoot, NBitsForKeypairDefault, nil, nil)
}

func DoInit(out io.Writer, repoRoot string, nBitsForKeypair int, confProfiles []string, conf *nconfig.Config) error {
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
	hash := c.Cid().String()

	if err := api.Pin().Add(nd.Context(), c); err != nil {
		return fmt.Errorf("init: pinning on init wallet failed: %s", err)
	}

	_, err = api.Name().Publish(nd.Context(), c)
	if err != nil {
		return fmt.Errorf("init: publish wallet address failed: %s", err)
	}

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
