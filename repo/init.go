package repo

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/namesys"
	loader "github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	logging "github.com/ipfs/go-log"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/repo/config"
)

var log = logging.Logger("tex-repo")

var ErrRepoExists = fmt.Errorf("repo not empty, reinitializing would overwrite your account")
var ErrRepoDoesNotExist = fmt.Errorf("repo does not exist, initialization is required")
var ErrMigrationRequired = fmt.Errorf("repo needs migration")
var ErrRepoCorrupted = fmt.Errorf("repo is corrupted")

const Repover = "18"

func Init(repoPath string, mobile bool, server bool) error {
	err := checkWriteable(repoPath)
	if err != nil {
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
	conf, err := config.InitIpfs(peerIdentity, mobile, server)
	if err != nil {
		return err
	}

	_, err = LoadPlugins(repoPath)
	if err != nil {
		return err
	}

	err = fsrepo.Init(repoPath, conf)
	if err != nil {
		return err
	}

	// create a folder for bots
	err = initializeBotFolder(repoPath)
	if err != nil {
		return err
	}

	// write default textile config
	tconf, err := config.Init()
	if err != nil {
		return err
	}
	err = config.Write(repoPath, tconf)
	if err != nil {
		return err
	}

	// write repo version
	repoverFile, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer repoverFile.Close()
	_, err = repoverFile.Write([]byte(Repover))
	if err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoPath)
}

func LoadPlugins(repoPath string) (*loader.PluginLoader, error) {
	// check if repo is accessible before loading plugins
	_, err := checkPermissions(repoPath)
	if err != nil {
		return nil, err
	}

	plugins, err := loader.NewPluginLoader(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error loading preloaded plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return nil, nil
	}

	if err := plugins.Inject(); err != nil {
		return nil, nil
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
		return os.MkdirAll(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}

	return err
}

func initializeBotFolder(repoPath string) error {
	botFolder := filepath.Join(repoPath, "bots")
	_, err := os.Stat(botFolder)
	if os.IsNotExist(err) {
		// dir doesnt exist, check that we can create it
		return os.MkdirAll(botFolder, 0775)
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
