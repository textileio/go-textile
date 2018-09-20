package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/schema"
	wutil "github.com/textileio/textile-go/util"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/namesys"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/repo/fsrepo"
	"os"
	"path"
)

var log = logging.MustGetLogger("repo")

var ErrRepoExists = errors.New("repo not empty, reinitializing would overwrite your keys")

const repover = "5"

func DoInit(repoRoot string, version string, initDatastore func() error) error {
	if err := checkWriteable(repoRoot); err != nil {
		return err
	}

	// handle migrations
	if fsrepo.IsInitialized(repoRoot) {
		// run all migrations if needed
		if err := migrateUp(repoRoot, "", false); err != nil {
			return err
		}
		return ErrRepoExists
	}
	log.Infof("initializing textile ipfs node at %s", repoRoot)

	// custom directories
	paths, err := schema.NewCustomSchemaManager(schema.Context{
		DataPath: repoRoot,
	})
	if err := paths.BuildSchemaDirectories(); err != nil {
		return err
	}

	// TODO: remove
	sk, _, err := wutil.PrivKeyFromMnemonic(nil)
	if err != nil {
		return err
	}

	// create an identity for the ipfs peer
	peerIdentity, err := wutil.IdentityConfig(sk)
	if err != nil {
		return err
	}
	conf, err := config.Init(peerIdentity, version)
	if err != nil {
		return err
	}
	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return err
	}

	// initialize sqlite datastore
	if err := initDatastore(); err != nil {
		return err
	}

	// write repo version
	repoverFile, err := os.Create(path.Join(repoRoot, "repover"))
	if err != nil {
		return err
	}
	defer repoverFile.Close()
	if _, err := repoverFile.Write([]byte(repover)); err != nil {
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

func destroyRepo(root string) error {
	// exclude logs
	paths := []string{"blocks", "datastore", "keystore", "tmp", "config", "datastore_spec", "repo.lock", "version"}
	for _, p := range paths {
		err := os.RemoveAll(fmt.Sprintf("%s/%s", root, p))
		if err != nil {
			return err
		}
	}
	return nil
}
