package core

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/util"

	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema/textile"
)

var vars = struct {
	initConfig InitConfig

	node   *Textile
	thread *Thread

	token string

	schemaHash string
}{
	initConfig: InitConfig{
		BaseRepoPath: "testdata/.textile1",
		Account:      keypair.Random(),
		ApiAddr:      fmt.Sprintf("127.0.0.1:9999"),
		Debug:        true,
	},
}

func TestRepoPath(t *testing.T) {
	conf := InitConfig{
		RepoPath: "path",
	}
	value, err := conf.Repo()
	if err != nil {
		t.Fatal(err)
	}
	if conf.RepoPath != value {
		t.Fatalf("repo path incorrect")
	}
}

func TestBaseRepoPath(t *testing.T) {
	conf := InitConfig{
		BaseRepoPath: "base",
		Account:      keypair.Random(),
	}
	target := path.Join(conf.BaseRepoPath, conf.Account.Address())
	value, err := conf.Repo()
	if err != nil {
		t.Fatal(err)
	}
	if target != value {
		t.Fatalf("repo path incorrect")
	}
}

func TestNoRepoExists(t *testing.T) {
	repo, err := vars.initConfig.Repo()
	if err != nil {
		t.Fatal(err)
	}
	_ = os.RemoveAll(repo)
	exists, err := vars.initConfig.RepoExists()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("repo should not exist via InitConfig.RepoExists, but it does")
	}
	exists = RepoExists(repo)
	if exists {
		t.Fatalf("repo should not exist via RepoExists, but it does")
	}
	exists = AccountRepoExists(vars.initConfig.BaseRepoPath, vars.initConfig.Account.Address())
	if exists {
		t.Fatalf("repo should not exist via AccountRepoExists, but it does")
	}
}

func TestInitRepo(t *testing.T) {
	if err := InitRepo(vars.initConfig); err != nil {
		t.Fatalf("init node failed: %s", err)
	}
}

func TestRepoExists(t *testing.T) {
	repo, err := vars.initConfig.Repo()
	if err != nil {
		t.Fatal(err)
	}
	exists, err := vars.initConfig.RepoExists()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("repo should exist via InitConfig.RepoExists, but it doesn't")
	}
	exists = RepoExists(repo)
	if !exists {
		t.Fatalf("repo should exist via RepoExists, but it doesn't")
	}
	exists = AccountRepoExists(vars.initConfig.BaseRepoPath, vars.initConfig.Account.Address())
	if !exists {
		t.Fatalf("repo should exist via AccountRepoExists, but it doesn't")
	}
}

func TestNewTextile(t *testing.T) {
	repo, err := vars.initConfig.Repo()
	if err != nil {
		t.Fatal(err)
	}
	vars.node, err = NewTextile(RunConfig{
		RepoPath: repo,
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("create node failed: %s", err)
	}
}

func TestSetLogLevel(t *testing.T) {
	logLevel := &pb.LogLevel{Systems: map[string]pb.LogLevel_Level{
		"tex-core":      pb.LogLevel_DEBUG,
		"tex-datastore": pb.LogLevel_INFO,
		"tex-service":   pb.LogLevel_DEBUG,
	}}
	if err := vars.node.SetLogLevel(logLevel, true); err != nil {
		t.Fatalf("set log levels failed: %s", err)
	}
}

func TestTextile_Start(t *testing.T) {
	if err := vars.node.Start(); err != nil {
		t.Fatalf("start node failed: %s", err)
	}
	<-vars.node.OnlineCh()
}

func TestTextile_Started(t *testing.T) {
	if !vars.node.Started() {
		t.Fatal("should report node started")
	}
}

func TestTextile_Online(t *testing.T) {
	if !vars.node.Online() {
		t.Fatal("should report node online")
	}
}

func TestTextile_AddContact(t *testing.T) {
	if err := vars.node.AddContact(util.TestContact); err != nil {
		t.Fatalf("add contact failed: %s", err)
	}
}

func TestTextile_AddContactAgain(t *testing.T) {
	if err := vars.node.AddContact(util.TestContact); err != nil {
		t.Fatal("adding duplicate contact should not throw error")
	}
}

func TestTextile_GetMedia(t *testing.T) {
	f, err := os.Open("../mill/testdata/image.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	media, err := vars.node.GetMillMedia(f, &mill.ImageResize{})
	if err != nil {
		t.Fatal(err)
	}
	if media != "image/jpeg" {
		t.Fatalf("wrong media type: %s", media)
	}
}

func TestTextile_AddSchema(t *testing.T) {
	file, err := vars.node.AddSchema(textile.Blob, "test")
	if err != nil {
		t.Fatal(err)
	}
	vars.schemaHash = file.Hash
}

func TestTextile_AddThread(t *testing.T) {
	var err error
	vars.thread, err = addTestThread(vars.node, &pb.AddThreadConfig{
		Key:       ksuid.New().String(),
		Name:      "test",
		Schema:    &pb.AddThreadConfig_Schema{Id: vars.schemaHash},
		Type:      pb.Thread_OPEN,
		Sharing:   pb.Thread_SHARED,
		Whitelist: []string{},
	})
	if err != nil {
		t.Fatalf("add thread failed: %s", err)
	}

	// add again w/ same key
	_, err = addTestThread(vars.node, &pb.AddThreadConfig{
		Key:       vars.thread.Key,
		Name:      "test2",
		Type:      pb.Thread_PUBLIC,
		Sharing:   pb.Thread_NOT_SHARED,
		Whitelist: []string{},
	})
	if err == nil {
		t.Fatal("add thread with same key should fail")
	}

	// add again w/ same key but force true
	forced, err := addTestThread(vars.node, &pb.AddThreadConfig{
		Key:       vars.thread.Key,
		Force:     true,
		Name:      "test3",
		Type:      pb.Thread_PUBLIC,
		Sharing:   pb.Thread_NOT_SHARED,
		Whitelist: []string{},
	})
	if err != nil {
		t.Fatalf("add thread with same key and force should not fail: %s", err)
	}
	if forced.Key != vars.thread.Key+"_1" {
		t.Fatal("add thread with same key and force resulted in bad key")
	}
}

func TestTextile_RenameThread(t *testing.T) {
	err := vars.node.RenameThread(vars.thread.Id, "new name")
	if err != nil {
		t.Fatalf("error renaming thread: %s", err)
	}

	thrd := vars.node.Thread(vars.thread.Id)
	if thrd.Name != "new name" {
		t.Fatal("error renaming thread")
	}
}

func TestTextile_AddFile(t *testing.T) {
	files, err := addData(vars.node, []string{"../mill/testdata/image.jpeg"}, vars.thread, "oi!")
	if err != nil {
		t.Fatal(err)
	}

	if files.Files[0].File.Checksum != "9sjWaHS2qRdjnaFGa394EjRCfJfZifNR3mwNysBxWTAX" {
		t.Fatal("wrong checksum")
	}
}

func TestTextile_Stop(t *testing.T) {
	err := vars.node.Stop()
	if err != nil {
		t.Fatalf("stop node failed: %s", err)
	}
}

func TestTextile_StartedAgain(t *testing.T) {
	if vars.node.Started() {
		t.Fatal("node should report stopped")
	}
}

func TestTextile_OnlineAgain(t *testing.T) {
	if vars.node.Online() {
		t.Fatal("node should report offline")
	}
}

func TestTextile_Teardown(t *testing.T) {
	vars.node = nil
}
