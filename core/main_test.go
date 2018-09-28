package core_test

import (
	"crypto/rand"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/cafe/models"
	. "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"os"
	"testing"
)

var repo = "testdata/.textile"

var node *Textile

func TestInitRepo(t *testing.T) {
	os.RemoveAll(repo)
	accnt := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  *accnt,
		RepoPath: repo,
		LogLevel: logging.DEBUG,
	}); err != nil {
		t.Errorf("init node failed: %s", err)
	}
}

func TestNewTextile(t *testing.T) {
	var err error
	node, err = NewTextile(RunConfig{
		RepoPath: repo,
		CafeAddr: os.Getenv("CAFE_ADDR"),
		LogLevel: logging.DEBUG,
	})
	if err != nil {
		t.Errorf("create node failed: %s", err)
	}
}

func TestCore_Start(t *testing.T) {
	if err := node.Start(); err != nil {
		t.Errorf("start node failed: %s", err)
	}
	<-node.Online()
}

func TestCore_Started(t *testing.T) {
	if !node.Started() {
		t.Errorf("should report started")
	}
}

func TestCore_IsOnline(t *testing.T) {
	if !node.IsOnline() {
		t.Errorf("should report online")
	}
}

func TestCore_CafeRegister(t *testing.T) {
	req := &models.ReferralRequest{
		Key:         os.Getenv("CAFE_REFERRAL_KEY"),
		Count:       1,
		Limit:       1,
		RequestedBy: "test",
	}
	res, err := node.CreateCafeReferral(req)
	if err != nil {
		t.Errorf("create referral for registration failed: %s", err)
		return
	}
	if len(res.RefCodes) == 0 {
		t.Error("create referral for registration got no codes")
		return
	}
	if err := node.CafeRegister(res.RefCodes[0]); err != nil {
		t.Errorf("register failed: %s", err)
		return
	}
}

func TestCore_CafeLogin(t *testing.T) {
	if err := node.CafeLogin(); err != nil {
		t.Errorf("login failed: %s", err)
		return
	}
}

func TestCore_AddThread(t *testing.T) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	thrd, err := node.AddThread("test", sk, true)
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	if thrd == nil {
		t.Error("add thread didn't return thread")
	}
}

func TestCore_AddPhoto(t *testing.T) {
	added, err := node.AddPhoto("../photo/testdata/image.jpg")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	if len(added.Id) == 0 {
		t.Errorf("add photo got bad id")
	}
	// test adding an image w/o the orientation tag
	added2, err := node.AddPhoto("../photo/testdata/image-no-orientation.jpg")
	if err != nil {
		t.Errorf("add photo w/o orientation tag failed: %s", err)
		return
	}
	if len(added2.Id) == 0 {
		t.Errorf("add photo w/o orientation tag got bad id")
	}
}

func TestCore_CafeLogout(t *testing.T) {
	err := node.CafeLogout()
	if err != nil {
		t.Errorf("logout failed: %s", err)
		return
	}
}

func TestCore_Stop(t *testing.T) {
	err := node.Stop()
	if err != nil {
		t.Errorf("stop node failed: %s", err)
	}
}

func TestCore_StartedAgain(t *testing.T) {
	if node.Started() {
		t.Errorf("should report stopped")
	}
}

func TestCore_OnlineAgain(t *testing.T) {
	if node.IsOnline() {
		t.Errorf("should report offline")
	}
}

// test cafe login in stopped state, should re-connect to db
func TestCore_LoginAgain(t *testing.T) {
	if err := node.CafeLogin(); err != nil {
		t.Errorf("login from stopped failed: %s", err)
		return
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(node.GetRepoPath())
}
