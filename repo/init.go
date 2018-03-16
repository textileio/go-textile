package repo

import (
	"fmt"
	"os"
	"io"
	"errors"
	"path"
	"context"
	"time"
	"crypto/sha256"
	"crypto/hmac"
	"encoding/base64"
	"bytes"

	"github.com/tyler-smith/go-bip39"

	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/repo/wallet"

	nconfig "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"

	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

const (
	NBitsForKeypair = 4096
)

var ErrRepoExists = errors.New(`ipfs configuration file already exists!
Reinitializing would overwrite your keys.
`)

func DoInit(out io.Writer, repoRoot string, creationDate time.Time, dbInit func(string, []byte, string, time.Time) error) error {
	if _, err := fmt.Fprintf(out, "initializing textile ipfs node at %s\n", repoRoot); err != nil {
		return err
	}

	if err := checkWriteable(repoRoot); err != nil {
		return err
	}

	paths, err := util.NewCustomSchemaManager(util.SchemaContext{
		DataPath: repoRoot,
	})
	if err := paths.BuildSchemaDirectories(); err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoRoot) {
		return ErrRepoExists
	}

	conf, err := config.Init(out, NBitsForKeypair)
	if err != nil {
		return err
	}

	fmt.Fprint(out, "generating Ed25519 keypair...")
	mnemonic, err := createMnemonic(bip39.NewEntropy, bip39.NewMnemonic)
	if err != nil {
		return err
	}
	seed := bip39.NewSeed(mnemonic, "")
	identityKey, err := identityKeyFromSeed(seed, NBitsForKeypair)
	if err != nil {
		return err
	}
	fmt.Printf("Done\n")

	//identity, err := identityFromKey(identityKey)
	//if err != nil {
	//	return err
	//}

	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return err
	}

	if err := dbInit(mnemonic, identityKey, "", creationDate); err != nil {
		return err
	}

	return initializeIpnsKeyspace(out, repoRoot)
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

func initializeIpnsKeyspace(out io.Writer, repoRoot string) error {
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

	// setup our wallet
	err = wallet.NewWalletData(nd)
	if err != nil {
		return fmt.Errorf("init: create empty wallet data failed: %s", err)
	}

	return nil
}

func createMnemonic(newEntropy func(int) ([]byte, error), newMnemonic func([]byte) (string, error)) (string, error) {
	entropy, err := newEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := newMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

func identityFromKey(privkey []byte) (nconfig.Identity, error) {
	ident := nconfig.Identity{}
	sk, err := libp2p.UnmarshalPrivateKey(privkey)
	if err != nil {
		return ident, err
	}
	skbytes, err := sk.Bytes()
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(sk.GetPublic())
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	return ident, nil
}

func identityKeyFromSeed(seed []byte, bits int) ([]byte, error) {
	hm := hmac.New(sha256.New, []byte("scythian horde"))
	hm.Write(seed)
	reader := bytes.NewReader(hm.Sum(nil))
	sk, _, err := libp2p.GenerateKeyPairWithReader(libp2p.Ed25519, bits, reader)
	if err != nil {
		return nil, err
	}
	encodedKey, err := sk.Bytes()
	if err != nil {
		return nil, err
	}
	return encodedKey, nil
}
