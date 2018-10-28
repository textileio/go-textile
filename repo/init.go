package repo

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/schema"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/namesys"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/fsrepo"
	"os"
	"path"
)

var log = logging.Logger("tex-repo")

var ErrRepoExists = errors.New("repo not empty, reinitializing would overwrite your account")
var ErrRepoDoesNotExist = errors.New("repo does not exist, initialization is required")
var ErrMigrationRequired = errors.New("repo needs migration")
var ErrRepoCorrupted = errors.New("repo is corrupted")

const repover = "6"

func DoInit(repoRoot string, initDatastore func() error) error {
	if err := checkWriteable(repoRoot); err != nil {
		return err
	}

	// double check if initialized
	if fsrepo.IsInitialized(repoRoot) {
		return ErrRepoExists
	}
	log.Infof("initializing repo at %s", repoRoot)

	// custom directories
	paths, err := schema.NewCustomSchemaManager(schema.Context{
		DataPath: repoRoot,
	})
	if err := paths.BuildSchemaDirectories(); err != nil {
		return err
	}

	// create an identity for the ipfs peer
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return err
	}
	peerIdentity, err := ipfs.IdentityConfig(sk)
	if err != nil {
		return err
	}
	conf, err := config.Init(peerIdentity)
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
