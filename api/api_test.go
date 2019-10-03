package api_test

import (
	"os"
	"testing"

	. "github.com/textileio/go-textile/api"
	"github.com/textileio/go-textile/bots"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
)

var initConfig = core.InitConfig{
	BaseRepoPath: "testdata/.textile",
	ApiAddr:      "127.0.0.1:9998",
}

func TestApi_Creation(t *testing.T) {
	initConfig.Account = keypair.Random()

	_ = os.RemoveAll(initConfig.RepoPath())

	err := core.InitRepo(initConfig)
	if err != nil {
		t.Errorf("init node failed: %s", err)
		return
	}

	node, err := core.NewTextile(core.RunConfig{
		RepoPath: initConfig.RepoPath(),
	})
	if err != nil {
		t.Errorf("create node failed: %s", err)
		return
	}

	bots := &bots.NewService(node)
	bots.RunAll(initConfig.BaseRepoPath, []string{})

	Host = &Api{
		Node:     node,
		Bots:     bots,
		PinCode:  "",
		RepoPath: initConfig.BaseRepoPath,
	}

	Host.Start(node.Config().Addresses.API, false)
}

func TestApi_Addr(t *testing.T) {
	if len(Host.Addr()) == 0 {
		t.Error("get gateway address failed")
		return
	}
}

func TestApi_Stop(t *testing.T) {
	err := Host.Stop()
	if err != nil {
		t.Errorf("stop gateway failed: %s", err)
	}
	_ = os.RemoveAll(initConfig.RepoPath())
}
