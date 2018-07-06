package core_test

import (
	"github.com/op/go-logging"
	. "github.com/textileio/textile-go/core"
	util "github.com/textileio/textile-go/util/testing"
	"github.com/textileio/textile-go/wallet"
	"os"
	"testing"
)

var repo = "testdata/.textile"

var node *TextileNode

func TestNewNode(t *testing.T) {
	os.RemoveAll(repo)
	config := NodeConfig{
		LogLevel: logging.DEBUG,
		LogFiles: false,
		WalletConfig: wallet.Config{
			RepoPath:   repo,
			CentralAPI: util.CentralApiURL,
			IsMobile:   false,
		},
	}
	var err error
	node, _, err = NewNode(config)
	if err != nil {
		t.Errorf("create node failed: %s", err)
	}
}

func TestTextileNode_StartWallet(t *testing.T) {
	online, err := node.StartWallet()
	if err != nil {
		t.Errorf("start node failed: %s", err)
	}
	<-online
}

func TestTextileNode_StartAgain(t *testing.T) {
	_, err := node.StartWallet()
	if err != wallet.ErrStarted {
		t.Errorf("start node again reported wrong error: %s", err)
	}
}

func TestTextileNode_StartGateway(t *testing.T) {
	err := node.StartGateway()
	if err != nil {
		t.Errorf("start gateway failed: %s", err)
	}
}

func TestTextileNode_GetGatewayAddress(t *testing.T) {
	if len(node.GetGatewayAddress()) == 0 {
		t.Error("get gateway address failed")
	}
}

func TestTextileNode_StopGateway(t *testing.T) {
	err := node.StopGateway()
	if err != nil {
		t.Errorf("stop gateway failed: %s", err)
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
