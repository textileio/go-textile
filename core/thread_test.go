package core_test

import (
	"crypto/rand"
	. "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/core/thread"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"os"
	"testing"
)

var trepo = "testdata/.textile1"

var twallet *Textile

var thrd *thread.Thread
var wadded *AddDataResult
var tadded mh.Multihash

func Test_SetupThread(t *testing.T) {
	os.RemoveAll(trepo)
	wconfig := Config{
		RepoPath: trepo,
	}
	var err error
	twallet, err = NewTextile(wconfig)
	if err != nil {
		t.Errorf("create wallet failed: %s", err)
		return
	}
	if err := twallet.Start(); err != nil {
		t.Errorf("start wallet failed: %s", err)
	}
}

func TestNewThread_WalletOffline(t *testing.T) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
		return
	}
	thrd, err = twallet.AddThread("thread1", sk, true)
	if err != nil {
		t.Errorf("create thread while offline failed: %s", err)
	}
}

func TestNewThread_WalletOnline(t *testing.T) {
	<-twallet.Online()
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = twallet.AddThread("thread2", sk, true)
	if err != nil {
		t.Errorf("create thread while online failed: %s", err)
	}
}

func TestThread_AddPhotoSetup(t *testing.T) {
	var err error
	wadded, err = twallet.AddPhoto("../photo/testdata/image.jpg")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	if len(wadded.Id) == 0 {
		t.Errorf("add photo got bad id")
	}
}

func TestThread_AddPhoto(t *testing.T) {
	var err error
	tadded, err = thrd.AddPhoto(wadded.Id, "howdy", []byte(wadded.Key))
	if err != nil {
		t.Errorf("add photo to thread failed: %s", err)
	}
	if tadded == nil {
		t.Error("add photo to thread got bad result")
	}
}

func Test_TeardownThread(t *testing.T) {
	os.RemoveAll(twallet.GetRepoPath())
}
