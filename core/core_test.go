package core_test

import (
	"fmt"
	"github.com/op/go-logging"
	. "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/wallet"
	"os"
	"testing"
)

var repo = "testdata/.textile"

var node *TextileNode

func TestNewNode(t *testing.T) {
	os.RemoveAll(repo)
	cfg := NodeConfig{
		LogLevel: logging.DEBUG,
		LogFiles: false,
		WalletConfig: wallet.Config{
			RepoPath: repo,
			IsMobile: false,
		},
	}
	var err error
	node, _, err = NewNode(cfg)
	if err != nil {
		t.Errorf("create node failed: %s", err)
	}
}

func TestTextileNode_StartWallet(t *testing.T) {
	if err := node.StartWallet(); err != nil {
		t.Errorf("start node failed: %s", err)
	}
	<-node.Wallet.Online()
}

func TestTextileNode_StartAgain(t *testing.T) {
	if err := node.StartWallet(); err != wallet.ErrStarted {
		t.Errorf("start node again reported wrong error: %s", err)
	}
}

func TestTextileNode_StartServer(t *testing.T) {
	node.StartGateway(fmt.Sprintf("127.0.0.1:%d", config.GetRandomPort()))
}

func TestTextileNode_GetGatewayAddr(t *testing.T) {
	if len(node.GetGatewayAddr()) == 0 {
		t.Error("get server address failed")
	}
}

func TestTextileNode_StopGateway(t *testing.T) {
	err := node.StopGateway()
	if err != nil {
		t.Errorf("stop server failed: %s", err)
	}
}

func TestTextileNode_Stop(t *testing.T) {
	err := node.StopWallet()
	if err != nil {
		t.Errorf("stop node failed: %s", err)
	}
	if node.Wallet.Started() {
		t.Errorf("should not report started")
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(node.Wallet.GetRepoPath())
}
