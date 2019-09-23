package gateway_test

import (
	"os"
	"testing"

	"github.com/textileio/go-textile/core"
	. "github.com/textileio/go-textile/gateway"
	"github.com/textileio/go-textile/keypair"
)

var initConfig = core.InitConfig{
	BaseRepoPath: "testdata/.textile",
	GatewayAddr:  "127.0.0.1:9998",
}

func TestGateway_Creation(t *testing.T) {
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

	Host = &Gateway{Node: node}
	Host.Start(node.Config().Addresses.Gateway)
}

func TestGateway_Addr(t *testing.T) {
	if len(Host.Addr()) == 0 {
		t.Error("get gateway address failed")
		return
	}
}

func TestGateway_Stop(t *testing.T) {
	err := Host.Stop()
	if err != nil {
		t.Errorf("stop gateway failed: %s", err)
	}
	_ = os.RemoveAll(initConfig.RepoPath())
}
