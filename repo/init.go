package repo

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path"

	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/core"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/namesys"
	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/repo/fsrepo"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo/config"
)

var log = logging.Logger("tex-repo")

var ErrRepoExists = errors.New("repo not empty, reinitializing would overwrite your account")
var ErrRepoDoesNotExist = errors.New("repo does not exist, initialization is required")
var ErrMigrationRequired = errors.New("repo needs migration")
var ErrRepoCorrupted = errors.New("repo is corrupted")

const repover = "10"

func Init(repoPath string, version string) error {
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
	if err := fsrepo.Init(repoPath, conf); err != nil {
		return err
	}

	// write default textile config
	tconf, err := config.Init(version)
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
	if _, err := repoverFile.Write([]byte(repover)); err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoPath)
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
