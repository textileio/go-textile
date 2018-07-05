package wallet_test

import (
	. "github.com/textileio/textile-go/wallet"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/thread"
	"os"
	"testing"
)

var trepo = "testdata/.ipfs1"

var twallet *Wallet
var wonline <-chan struct{}

var thrd *thread.Thread
var wadded *model.AddResult
var tadded *model.AddResult

func Test_SetupThread(t *testing.T) {
	os.RemoveAll(trepo)
	wconfig := Config{
		RepoPath: trepo,
	}
	var err error
	twallet, _, err = NewWallet(wconfig)
	if err != nil {
		t.Errorf("create wallet failed: %s", err)
	}
	wonline, err = twallet.Start()
	if err != nil {
		t.Errorf("start wallet failed: %s", err)
	}
}

func TestNewThread_WalletOffline(t *testing.T) {
	var err error
	thrd, _, err = twallet.AddThreadWithMnemonic("thread1", nil)
	if err != nil {
		t.Errorf("create thread while offline failed: %s", err)
	}
}

func TestNewThread_WalletOnline(t *testing.T) {
	<-wonline
	var err error
	_, _, err = twallet.AddThreadWithMnemonic("thread2", nil)
	if err != nil {
		t.Errorf("create thread while online failed: %s", err)
	}
}

func TestThread_AddPhotoSetup(t *testing.T) {
	var err error
	wadded, err = twallet.AddPhoto("testdata/image.jpg")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	if len(wadded.Id) == 0 {
		t.Errorf("add photo got bad id")
	}
}

func TestThread_PostHeadPreContent(t *testing.T) {
	err := thrd.PostHead()
	if err != nil {
		t.Errorf("post head with no content failed: %s", err)
	}
}

func TestThread_AddPhoto(t *testing.T) {
	var err error
	tadded, err = thrd.AddPhoto(wadded.Id, "howdy", wadded.Key)
	if err != nil {
		t.Errorf("add photo to thread failed: %s", err)
	}
	if tadded.Id == "" {
		t.Error("add photo to thread got bad id")
	}
}

func TestThread_GetBlockData(t *testing.T) {
	// TODO
}

func TestThread_GetBlockDataBase64(t *testing.T) {
	// TODO
}

func TestThread_GetFileKey(t *testing.T) {
	// TODO
}

func TestThread_GetFileData(t *testing.T) {
	// TODO
}

func TestThread_GetFileDataBase64(t *testing.T) {
	// TODO
}

func TestThread_GetPhotoMetaData(t *testing.T) {
	// TODO
}

func TestThread_Subscribe(t *testing.T) {
	// TODO
}

func TestThread_Unsubscribe(t *testing.T) {
	// TODO
}

func TestThread_Listening(t *testing.T) {
	// TODO
}

func TestThread_Blocks(t *testing.T) {
	// TODO
}

func TestThread_Encrypt(t *testing.T) {
	// TODO
}

func TestThread_Decrypt(t *testing.T) {
	// TODO
}

func TestThread_PostHead(t *testing.T) {
	err := thrd.PostHead()
	if err != nil {
		t.Errorf("post head with content failed: %s", err)
	}
}

func TestThread_Peers(t *testing.T) {
	// TODO
}

func Test_TeardownThread(t *testing.T) {
	os.RemoveAll(twallet.GetRepoPath())
}
