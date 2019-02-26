package core_test

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"

	"github.com/segmentio/ksuid"
	. "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/schema/textile"
)

var repoPath = "testdata/.textile"
var otherPath = "testdata/.textile2"
var node *Textile
var other *Textile
var token string
var contact = &pb.Contact{
	Id:       "abcde",
	Address:  "address1",
	Username: "joe",
	Avatar:   "Qm123",
	Inboxes: []*pb.Cafe{{
		Peer:     "peer",
		Address:  "address",
		Api:      "v0",
		Protocol: "/textile/cafe/1.0.0",
		Node:     "v1.0.0",
		Url:      "https://mycafe.com",
	}},
}

var schemaHash string

func TestInitRepo(t *testing.T) {
	os.RemoveAll(repoPath)
	accnt := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  accnt,
		RepoPath: repoPath,
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

func TestSetLogLevel(t *testing.T) {
	logLevel := &pb.LogLevel{Systems: map[string]pb.LogLevel_Level{
		"tex-core":      pb.LogLevel_DEBUG,
		"tex-datastore": pb.LogLevel_INFO,
	}}
	if err := node.SetLogLevel(logLevel); err != nil {
		t.Errorf("set log levels failed: %s", err)
	}
}

func TestTextile_Start(t *testing.T) {
	if err := node.Start(); err != nil {
		t.Errorf("start node failed: %s", err)
	}
	<-node.OnlineCh()
}

func TestTextile_CafeSetup(t *testing.T) {
	// start another
	os.RemoveAll(otherPath)
	accnt := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:     accnt,
		RepoPath:    otherPath,
		CafeApiAddr: "127.0.0.1:5000",
		CafeOpen:    true,
	}); err != nil {
		t.Errorf("init other failed: %s", err)
		return
	}
	var err error
	other, err = NewTextile(RunConfig{
		RepoPath: otherPath,
	})
	if err != nil {
		t.Errorf("create other failed: %s", err)
		return
	}
	other.Start()

	// wait for cafe to be online
	<-other.OnlineCh()
}

func TestTextile_Started(t *testing.T) {
	if !node.Started() {
		t.Errorf("should report node started")
	}
	if !other.Started() {
		t.Errorf("should report other started")
	}
}

func TestTextile_Online(t *testing.T) {
	if !node.Online() {
		t.Errorf("should report node online")
	}
	if !other.Online() {
		t.Errorf("should report other online")
	}
}

func TestTextile_CafeTokens(t *testing.T) {
	var err error
	token, err = other.CreateCafeToken("", true)
	if err != nil {
		t.Error(fmt.Errorf("error creating cafe token: %s", err))
		return
	}
	if len(token) == 0 {
		t.Error("invalid token created")
		return
	}

	tokens, _ := other.CafeTokens()
	if len(tokens) < 1 {
		t.Error("token database not updated (should be length 1)")
	}

	if ok, err := other.ValidateCafeToken("blah"); err == nil || ok {
		t.Error("expected token comparison with 'blah' to be invalid")
	}

	if ok, err := other.ValidateCafeToken(token); err != nil || !ok {
		t.Error("expected token comparison to be valid")
	}
}

func TestTextile_CafeRegistration(t *testing.T) {
	// register w/ wrong credentials
	if _, err := node.RegisterCafe("http://127.0.0.1:5000", "blah"); err == nil {
		t.Error("register node w/ other should have failed")
		return
	}

	// register cafe
	if _, err := node.RegisterCafe("http://127.0.0.1:5000", token); err != nil {
		t.Errorf("register node w/ other failed: %s", err)
		return
	}

	// get sessions
	sessions := node.CafeSessions()
	if len(sessions.Items) > 0 {
		session = sessions.Items[0]
	} else {
		t.Errorf("no active sessions")
	}
}

func TestTextile_AddContact(t *testing.T) {
	if err := node.AddContact(contact); err != nil {
		t.Errorf("add contact failed: %s", err)
	}
}

func TestTextile_AddContactAgain(t *testing.T) {
	if err := node.AddContact(contact); err != nil {
		t.Errorf("adding duplicate contact should not throw error")
	}
}

func TestTextile_GetMedia(t *testing.T) {
	f, err := os.Open("../mill/testdata/image.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	media, err := node.GetMedia(f, &mill.ImageResize{})
	if err != nil {
		t.Fatal(err)
	}
	if media != "image/jpeg" {
		t.Errorf("wrong media type: %s", media)
	}
}

func TestTextile_AddSchema(t *testing.T) {
	file, err := node.AddSchema(textile.Media, "test")
	if err != nil {
		t.Fatal(err)
	}
	schemaHash = file.Hash
}

func TestTextile_AddThread(t *testing.T) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	config := pb.AddThreadConfig{
		Key:     ksuid.New().String(),
		Name:    "test",
		Schema:  &pb.AddThreadConfig_Schema{Id: schemaHash},
		Type:    pb.Thread_Open,
		Sharing: pb.Thread_Shared,
		Members: []string{},
	}
	thrd, err := node.AddThread(config, sk, node.Account().Address(), true)
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	if thrd == nil {
		t.Error("add thread didn't return thread")
	}

	// add again w/ same key
	sk2, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	if _, err := node.AddThread(pb.AddThreadConfig{
		Key:     config.Key,
		Name:    "test2",
		Type:    pb.Thread_Public,
		Sharing: pb.Thread_NotShared,
		Members: []string{},
	}, sk2, node.Account().Address(), true); err == nil {
		t.Error("add thread with same key should fail")
		return
	}

	// add again w/ same key but force true
	sk3, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	forced, err := node.AddThread(pb.AddThreadConfig{
		Key:     config.Key,
		Force:   true,
		Name:    "test3",
		Type:    pb.Thread_Public,
		Sharing: pb.Thread_NotShared,
		Members: []string{},
	}, sk3, node.Account().Address(), true)
	if err != nil {
		t.Errorf("add thread with same key and force should not fail: %s", err)
		return
	}
	if forced.Key != config.Key+"_1" {
		t.Errorf("add thread with same key and force resulted in bad key")
	}
}

func TestTextile_AddFile(t *testing.T) {
	f, err := os.Open("../mill/testdata/image.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	m := &mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   "200",
			Quality: "75",
		},
	}
	conf := AddFileConfig{
		Input: data,
		Name:  "image.jpeg",
		Media: "image/jpeg",
	}

	file, err := node.AddFileIndex(m, conf)
	if err != nil {
		t.Errorf("add file failed: %s", err)
		return
	}

	if file.Mill != "/image/resize" {
		t.Error("wrong mill")
	}
	if file.Checksum != "3LLsrJ4zcF66d9r5pMnip243y2zZQpkKShhYVHn4Hk2j" {
		t.Error("wrong checksum")
	}
}

func TestTextile_RemoveCafeToken(t *testing.T) {
	if err := other.RemoveCafeToken(token); err != nil {
		t.Error("expected be remove token cleanly")
	}

	tokens, _ := other.CafeTokens()
	if len(tokens) > 0 {
		t.Error("token database not updated (should be zero length)")
	}
}

func TestTextile_Stop(t *testing.T) {
	if err := node.Stop(); err != nil {
		t.Errorf("stop node failed: %s", err)
	}
	if err := other.Stop(); err != nil {
		t.Errorf("stop other failed: %s", err)
	}
}

func TestTextile_StartedAgain(t *testing.T) {
	if node.Started() {
		t.Errorf("node should report stopped")
	}
	if other.Started() {
		t.Errorf("other should report stopped")
	}
}

func TestTextile_OnlineAgain(t *testing.T) {
	if node.Online() {
		t.Errorf("node should report offline")
	}
	if other.Online() {
		t.Errorf("other should report offline")
	}
}

func TestTextile_Teardown(t *testing.T) {
	node = nil
	os.RemoveAll(repoPath)
	other = nil
	os.RemoveAll(otherPath)
}
