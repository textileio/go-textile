package core_test

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	libp2pc "github.com/libp2p/go-libp2p-crypto"
	"github.com/segmentio/ksuid"
	. "github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema/textile"
	"github.com/textileio/go-textile/util"
)

var repoPath = "testdata/.textile"
var otherPath = "testdata/.textile2"

var node *Textile
var other *Textile

var testThread *Thread

var token string

var contact = &pb.Contact{
	Address: "address1",
	Name:    "joe",
	Avatar:  "Qm123",
	Peers: []*pb.Peer{
		{
			Id:      "abcde",
			Address: "address1",
			Name:    "joe",
			Avatar:  "Qm123",
			Inboxes: []*pb.Cafe{{
				Peer:     "peer",
				Address:  "address",
				Api:      "v0",
				Protocol: "/textile/cafe/1.0.0",
				Node:     "v1.0.0",
				Url:      "https://mycafe.com",
			}},
		},
	},
}

var schemaHash string

func TestInitRepo(t *testing.T) {
	_ = os.RemoveAll(repoPath)
	accnt := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  accnt,
		RepoPath: repoPath,
		ApiAddr:  fmt.Sprintf("127.0.0.1:%s", GetRandomPort()),
	}); err != nil {
		t.Fatalf("init node failed: %s", err)
	}
}

func TestNewTextile(t *testing.T) {
	var err error
	node, err = NewTextile(RunConfig{
		RepoPath: repoPath,
	})
	if err != nil {
		t.Fatalf("create node failed: %s", err)
	}
}

func TestSetLogLevel(t *testing.T) {
	logLevel := &pb.LogLevel{Systems: map[string]pb.LogLevel_Level{
		"tex-core":      pb.LogLevel_DEBUG,
		"tex-datastore": pb.LogLevel_INFO,
	}}
	if err := node.SetLogLevel(logLevel); err != nil {
		t.Fatalf("set log levels failed: %s", err)
	}
}

func TestTextile_Start(t *testing.T) {
	if err := node.Start(); err != nil {
		t.Fatalf("start node failed: %s", err)
	}
	<-node.OnlineCh()
}

func TestTextile_API_Start(t *testing.T) {
	node.StartApi(node.Config().Addresses.API, false)
}

func TestTextile_API_Addr(t *testing.T) {
	if len(node.ApiAddr()) == 0 {
		t.Error("get api address failed")
		return
	}
}

func TestTextile_API_Health(t *testing.T) {
	// prepare the URL
	addr := "http://" + node.ApiAddr() + "/health"

	// test the request
	util.TestURL(t, addr, http.MethodGet, http.StatusNoContent)
}

func TestTextile_API_Stop(t *testing.T) {
	if err := node.StopApi(); err != nil {
		t.Errorf("stop api failed: %s", err)
		return
	}
}

func TestTextile_CafeSetup(t *testing.T) {
	// start another
	_ = os.RemoveAll(otherPath)
	accnt := keypair.Random()
	err := InitRepo(InitConfig{
		Account:     accnt,
		RepoPath:    otherPath,
		CafeApiAddr: "127.0.0.1:5000",
		CafeURL:     "http://127.0.0.1:5000",
		CafeOpen:    true,
	})
	if err != nil {
		t.Fatalf("init other failed: %s", err)
	}
	other, err = NewTextile(RunConfig{
		RepoPath: otherPath,
	})
	if err != nil {
		t.Fatalf("create other failed: %s", err)
	}
	err = other.Start()
	if err != nil {
		t.Fatalf("start other failed: %s", err)
	}

	// wait for cafe to be online
	<-other.OnlineCh()
}

func TestTextile_Started(t *testing.T) {
	if !node.Started() {
		t.Fatal("should report node started")
	}
	if !other.Started() {
		t.Fatal("should report other started")
	}
}

func TestTextile_Online(t *testing.T) {
	if !node.Online() {
		t.Fatal("should report node online")
	}
	if !other.Online() {
		t.Fatal("should report other online")
	}
}

func TestTextile_CafeTokens(t *testing.T) {
	var err error
	token, err = other.CreateCafeToken("", true)
	if err != nil {
		t.Fatalf("error creating cafe token: %s", err)
	}
	if len(token) == 0 {
		t.Fatal("invalid token created")
	}

	tokens, _ := other.CafeTokens()
	if len(tokens) < 1 {
		t.Fatal("token database not updated (should be length 1)")
	}

	ok, err := other.ValidateCafeToken("blah")
	if err == nil || ok {
		t.Fatal("expected token comparison with 'blah' to be invalid")
	}

	ok, err = other.ValidateCafeToken(token)
	if err != nil || !ok {
		t.Fatal("expected token comparison to be valid")
	}
}

func TestTextile_CafeRegistration(t *testing.T) {
	// register w/ wrong credentials
	_, err := node.RegisterCafe("http://127.0.0.1:5000", "blah")
	if err == nil {
		t.Fatal("register node w/ other should have failed")
	}

	// register cafe
	_, err = node.RegisterCafe("http://127.0.0.1:5000", token)
	if err != nil {
		t.Fatalf("register node w/ other failed: %s", err)
	}

	// get sessions
	sessions := node.CafeSessions()
	if len(sessions.Items) > 0 {
		session = sessions.Items[0]
	} else {
		t.Fatal("no active sessions")
	}
}

func TestTextile_AddContact(t *testing.T) {
	if err := node.AddContact(contact); err != nil {
		t.Fatalf("add contact failed: %s", err)
	}
}

func TestTextile_AddContactAgain(t *testing.T) {
	if err := node.AddContact(contact); err != nil {
		t.Fatal("adding duplicate contact should not throw error")
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
		t.Fatalf("wrong media type: %s", media)
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
		t.Fatal(err)
	}
	config := pb.AddThreadConfig{
		Key:       ksuid.New().String(),
		Name:      "test",
		Schema:    &pb.AddThreadConfig_Schema{Id: schemaHash},
		Type:      pb.Thread_OPEN,
		Sharing:   pb.Thread_SHARED,
		Whitelist: []string{},
	}
	thrd, err := node.AddThread(config, sk, node.Account().Address(), true, true)
	if err != nil {
		t.Fatalf("add thread failed: %s", err)
	}
	if thrd == nil {
		t.Fatal("add thread didn't return thread")
	}
	testThread = thrd

	// add again w/ same key
	sk2, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, err = node.AddThread(pb.AddThreadConfig{
		Key:       config.Key,
		Name:      "test2",
		Type:      pb.Thread_PUBLIC,
		Sharing:   pb.Thread_NOT_SHARED,
		Whitelist: []string{},
	}, sk2, node.Account().Address(), true, true)
	if err == nil {
		t.Fatal("add thread with same key should fail")
	}

	// add again w/ same key but force true
	sk3, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	forced, err := node.AddThread(pb.AddThreadConfig{
		Key:       config.Key,
		Force:     true,
		Name:      "test3",
		Type:      pb.Thread_PUBLIC,
		Sharing:   pb.Thread_NOT_SHARED,
		Whitelist: []string{},
	}, sk3, node.Account().Address(), true, true)
	if err != nil {
		t.Fatalf("add thread with same key and force should not fail: %s", err)
	}
	if forced.Key != config.Key+"_1" {
		t.Fatal("add thread with same key and force resulted in bad key")
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
		t.Fatalf("add file failed: %s", err)
	}

	if file.Mill != "/image/resize" {
		t.Fatal("wrong mill")
	}
	if file.Checksum != "EWiMoePQAUrY9GWHBYu71upb9Z5dj1q9D3bS9Xtfp5fe" {
		t.Fatal("wrong checksum")
	}
}

func TestTextile_RenameThread(t *testing.T) {
	err := node.RenameThread(testThread.Id, "new name")
	if err != nil {
		t.Fatalf("error renaming thread: %s", err)
	}

	thrd := node.Thread(testThread.Id)
	if thrd.Name != "new name" {
		t.Fatal("error renaming thread")
	}
}

func TestTextile_RemoveCafeToken(t *testing.T) {
	err := other.RemoveCafeToken(token)
	if err != nil {
		t.Fatal("expected be remove token cleanly")
	}

	tokens, _ := other.CafeTokens()
	if len(tokens) > 0 {
		t.Fatal("token database not updated (should be zero length)")
	}
}

func TestTextile_Stop(t *testing.T) {
	err := node.Stop()
	if err != nil {
		t.Fatalf("stop node failed: %s", err)
	}
	err = other.Stop()
	if err != nil {
		t.Fatalf("stop other failed: %s", err)
	}
}

func TestTextile_StartedAgain(t *testing.T) {
	if node.Started() {
		t.Fatal("node should report stopped")
	}
	if other.Started() {
		t.Fatal("other should report stopped")
	}
}

func TestTextile_OnlineAgain(t *testing.T) {
	if node.Online() {
		t.Fatal("node should report offline")
	}
	if other.Online() {
		t.Fatal("other should report offline")
	}
}

func TestTextile_Teardown(t *testing.T) {
	node = nil
	_ = os.RemoveAll(repoPath)
	other = nil
	_ = os.RemoveAll(otherPath)
}
