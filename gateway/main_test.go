package gateway_test

import (
	"fmt"
	"github.com/textileio/go-textile/util"
	"net/http"
	"os"
	"testing"

	"github.com/textileio/go-textile/core"
	. "github.com/textileio/go-textile/gateway"
	"github.com/textileio/go-textile/keypair"
)

var repoPath = "testdata/.textile"

func TestGateway_Creation(t *testing.T) {
	os.RemoveAll(repoPath)

	err := core.InitRepo(core.InitConfig{
		Account:     keypair.Random(),
		RepoPath:    repoPath,
		GatewayAddr: fmt.Sprintf("127.0.0.1:%s", core.GetRandomPort()),
	})
	if err != nil {
		t.Errorf("init node failed: %s", err)
		return
	}

	node, err := core.NewTextile(core.RunConfig{
		RepoPath: repoPath,
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

func TestGateway_Health(t *testing.T) {
	// prepare the URL
	addr := "http://" + Host.Addr() + "/health"

	// test the request
	util.TestURL(t, addr, http.MethodGet, http.StatusNoContent)
}

func TestGateway_Stop(t *testing.T) {
	err := Host.Stop()
	if err != nil {
		t.Errorf("stop gateway failed: %s", err)
	}
	os.RemoveAll(repoPath)
}
