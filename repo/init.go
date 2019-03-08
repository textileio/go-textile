package repo

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/namesys"
	loader "gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/plugin/loader"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo/fsrepo"
	libp2pc "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	logging "gx/ipfs/QmbkT7eMTyXfpeyB3ZMxxcxg7XH8t6uXp49jqzz4HB7BGF/go-log"

	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/repo/config"
)

var log = logging.Logger("tex-repo")

var ErrRepoExists = errors.New("repo not empty, reinitializing would overwrite your account")
var ErrRepoDoesNotExist = errors.New("repo does not exist, initialization is required")
var ErrMigrationRequired = errors.New("repo needs migration")
var ErrRepoCorrupted = errors.New("repo is corrupted")

const Repover = "11"

func Init(repoPath string) error {
	if err := checkWriteable(repoPath); err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoPath) {
		return ErrRepoExists
	}
	log.Infof("initializing repo at %s", repoPath)

	// create an identity for the ipfs peer
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return err
	}
	peerIdentity, err := ipfs.IdentityConfig(sk)
	if err != nil {
		return err
	}

	// initialize ipfs config
	conf, err := config.InitIpfs(peerIdentity)
	if err != nil {
		return err
	}

	if _, err := LoadPlugins(repoPath); err != nil {
		return err
	}

	if err := fsrepo.Init(repoPath, conf); err != nil {
		return err
	}

	// write default textile config
	tconf, err := config.Init()
	if err != nil {
		return err
	}
	if err := config.Write(repoPath, tconf); err != nil {
		return err
	}

	// write repo version
	repoverFile, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer repoverFile.Close()
	if _, err := repoverFile.Write([]byte(Repover)); err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoPath)
}

func LoadPlugins(repoPath string) (*loader.PluginLoader, error) {
	pluginpath := filepath.Join(repoPath, "plugins")

	// check if repo is accessible before loading plugins
	var plugins *loader.PluginLoader
	ok, err := checkPermissions(repoPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		pluginpath = ""
	}
	plugins, err = loader.NewPluginLoader(pluginpath)
	if err != nil {
		log.Error("error loading plugins: ", err)
	}

	if err := plugins.Initialize(); err != nil {
		log.Error("error initializing plugins: ", err)
	}

	if err := plugins.Inject(); err != nil {
		log.Error("error running plugins: ", err)
	}
	return plugins, nil
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

	return namesys.InitializeKeyspace(ctx, nd.Namesys, nd.Pinning, nd.PrivateKey)
}

func checkPermissions(path string) (bool, error) {
	_, err := os.Open(path)
	if os.IsNotExist(err) {
		// repo does not exist yet - don't load plugins, but also don't fail
		return false, nil
	}
	if os.IsPermission(err) {
		// repo is not accessible. error out.
		return false, fmt.Errorf("error opening repository at %s: permission denied", path)
	}

	return true, nil
}
