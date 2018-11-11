package core_test

import (
	"crypto/rand"
	"github.com/segmentio/ksuid"
	. "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	logger "gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"os"
	"testing"
)

var repoPath = "testdata/.textile"
var node *Textile

func TestInitRepo(t *testing.T) {
	os.RemoveAll(repoPath)
	accnt := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  accnt,
		RepoPath: repoPath,
		LogLevel: logger.ERROR,
	}); err != nil {
		t.Errorf("init node failed: %s", err)
	}
}

func TestNewTextile(t *testing.T) {
	var err error
	node, err = NewTextile(RunConfig{
		RepoPath: repoPath,
	})
	if err != nil {
		t.Errorf("create node failed: %s", err)
	}
}

func TestCore_Start(t *testing.T) {
	if err := node.Start(); err != nil {
		t.Errorf("start node failed: %s", err)
	}
	<-node.OnlineCh()
}

func TestCore_Started(t *testing.T) {
	if !node.Started() {
		t.Errorf("should report started")
	}
}

func TestCore_Online(t *testing.T) {
	if !node.Online() {
		t.Errorf("should report online")
	}
}

func TestCore_CafeRegister(t *testing.T) {
	// TODO
}

func TestCore_AddThread(t *testing.T) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	schema, err := mh.FromB58String("QmUp6zZ6mNCCqcWfaoofcXPFB1CBhBtXJVLCE2gMPTuoVS")
	if err != nil {
		t.Error(err)
	}
	config := NewThreadConfig{
		Key:    ksuid.New().String(),
		Name:   "test",
		Schema: schema,
		Type:   repo.OpenThread,
		Join:   true,
	}
	thrd, err := node.AddThread(sk, config)
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	if thrd == nil {
		t.Error("add thread didn't return thread")
	}
}

//func TestCore_AddImage(t *testing.T) {
//	added, err := node.AddImageByPath("../images/testdata/image.jpg")
//	if err != nil {
//		t.Errorf("add image failed: %s", err)
//		return
//	}
//	if len(added.Id) == 0 {
//		t.Errorf("add image got bad id")
//	}
//	// test adding an image w/o the orientation tag
//	added2, err := node.AddImageByPath("../images/testdata/image-no-orientation.jpg")
//	if err != nil {
//		t.Errorf("add image w/o orientation tag failed: %s", err)
//		return
//	}
//	if len(added2.Id) == 0 {
//		t.Errorf("add photo w/o orientation tag got bad id")
//	}
//}

func TestCore_Stop(t *testing.T) {
	if err := node.Stop(); err != nil {
		t.Errorf("stop node failed: %s", err)
	}
}

func TestCore_StartedAgain(t *testing.T) {
	if node.Started() {
		t.Errorf("should report stopped")
	}
}

func TestCore_OnlineAgain(t *testing.T) {
	if node.Online() {
		t.Errorf("should report offline")
	}
}

func TestCore_Teardown(t *testing.T) {
	node = nil
}
