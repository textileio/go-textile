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
	"io/ioutil"
	"os"
	"path"
	"time"
)

var log = logging.MustGetLogger("repo")

var ErrRepoExists = errors.New("repo not empty, reinitializing would overwrite your keys")

const versionFilename = "textile_version"

func DoInit(repoRoot string, version string, mnemonic *string, initDB func(string) error, initConfig func(time.Time) error) (string, error) {
	if err := checkWriteable(repoRoot); err != nil {
		return "", err
	}

	versionPath := fmt.Sprintf("%s/%s", repoRoot, versionFilename)
	if fsrepo.IsInitialized(repoRoot) {
		// check version
		var onDiskVersion []byte
		onDiskVersion, _ = ioutil.ReadFile(versionPath)
		if version == string(onDiskVersion) {
			return "", ErrRepoExists
		} else {
			log.Info("old repo found, destroying...")
			if err := destroyRepo(repoRoot); err != nil {
				return "", err
			}
		}
	}
	log.Infof("initializing textile ipfs node at %s", repoRoot)

	if err := ioutil.WriteFile(versionPath, []byte(version), 0644); err != nil {
		return "", err
	}

	paths, err := schema.NewCustomSchemaManager(schema.Context{
		DataPath: repoRoot,
	})
	if err := paths.BuildSchemaDirectories(); err != nil {
		return "", err
	}

	sk, mnem, err := wutil.PrivKeyFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	identity, err := wutil.IdentityConfig(sk)
	if err != nil {
		return "", err
	}

	conf, err := config.Init(identity, version)
	if err != nil {
		return "", err
	}

	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return "", err
	}

	if err := initDB(""); err != nil {
		return "", err
	}

	if err := initConfig(time.Now()); err != nil {
		return "", err
	}

	return mnem, initializeIpnsKeyspace(repoRoot)
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
