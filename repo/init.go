package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/op/go-logging"

	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/schema"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/namesys"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/repo/fsrepo"
)

const (
	NBitsForKeypair = 4096
)

var log = logging.MustGetLogger("repo")

var ErrRepoExists = errors.New(`ipfs configuration file already exists!
Reinitializing would overwrite your keys.
`)

func DoInit(repoRoot string, isMobile bool, dbInit func(string) error, dbConfigure func(time.Time) error) error {
	if err := checkWriteable(repoRoot); err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoRoot) {
		return ErrRepoExists
	}
	log.Infof("initializing textile ipfs node at %s", repoRoot)

	paths, err := schema.NewCustomSchemaManager(schema.SchemaContext{
		DataPath: repoRoot,
	})
	if err := paths.BuildSchemaDirectories(); err != nil {
		return err
	}

	conf, err := config.Init(NBitsForKeypair, isMobile)
	if err != nil {
		return err
	}

	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return err
	}

	if err := dbInit(""); err != nil {
		return err
	}

	if err := dbConfigure(time.Now()); err != nil {
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
