package core

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
)

func CreateAndStartPeer(conf InitConfig, wait bool) (*Textile, error) {
	conf.Account = keypair.Random()

	repo, err := conf.Repo()
	if err != nil {
		return nil, err
	}

	_ = os.RemoveAll(repo)

	err = InitRepo(conf)
	if err != nil {
		return nil, err
	}
	node, err := NewTextile(RunConfig{
		RepoPath: repo,
		Debug:    conf.Debug,
	})
	if err != nil {
		return nil, err
	}
	err = node.Start()
	if err != nil {
		return nil, err
	}

	if wait {
		<-node.OnlineCh()
	}

	return node, nil
}

func addTestThread(node *Textile, conf *pb.AddThreadConfig) (*Thread, error) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	thrd, err := node.AddThread(*conf, sk, node.Account().Address(), true, true)
	if err != nil {
		return nil, err
	}
	if thrd == nil {
		return nil, fmt.Errorf("thread is nil")
	}

	return thrd, nil
}

func addData(node *Textile, names []string, thread *Thread, caption string) (*pb.Files, error) {
	var files []*pb.FileIndex

	for _, name := range names {
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}

		mil := &mill.Blob{}
		media, err := node.GetMillMedia(f, mil)
		if err != nil {
			return nil, err
		}

		_, _ = f.Seek(0, 0)
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		_, fname := filepath.Split(f.Name())
		file, err := node.AddFileIndex(mil, AddFileConfig{
			Input: data,
			Name:  fname,
			Media: media,
		})
		if err != nil {
			return nil, err
		}
		files = append(files, file)
		f.Close()
	}

	nd, keys, err := node.AddNodeFromFiles(files)
	if err != nil {
		return nil, err
	}

	hash, err := thread.AddFiles(nd, "", caption, keys.Files)
	if err != nil {
		return nil, err
	}

	return node.File(hash.B58String())
}
