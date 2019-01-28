package core_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
)

var tokenNode *Textile
var tokenRepoPath = "testdata/.textile3"

func TestTokens_Setup(t *testing.T) {
	// start node
	os.RemoveAll(tokenRepoPath)
	accnt := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  accnt,
		RepoPath: tokenRepoPath,
	}); err != nil {
		t.Errorf("init node failed: %s", err)
		return
	}
	var err error
	tokenNode, err = NewTextile(RunConfig{
		RepoPath: tokenRepoPath,
	})
	if err != nil {
		t.Errorf("create node failed: %s", err)
		return
	}
	tokenNode.Start()

	// wait for peer to be online
	<-tokenNode.OnlineCh()
}

func TestTokens_TestAll(t *testing.T) {
	token, err := tokenNode.CreateCafeToken()
	if err != nil {
		t.Error(fmt.Errorf("error creating cafe token: %s", err))
		return
	}
	if len(token.Token) == 0 {
		t.Error("invalid token created")
	}

	tokens, _ := tokenNode.CafeDevTokens()
	if len(tokens) < 1 {
		t.Error("token database not updated (should be length 1)")
	}

	if ok, err := tokenNode.CompareCafeDevToken(token.Id, "blah"); err == nil || ok {
		t.Error("expected token comparison with 'blah' to be invalid")
	}

	if ok, err := tokenNode.CompareCafeDevToken(token.Id, token.Token); err != nil || !ok {
		t.Error("expected token comparison to be valid")
	}

	if err = tokenNode.RemoveCafeDevToken(token.Id); err != nil {
		t.Error("expected be remove dev token cleanly")
	}

	tokens, _ = tokenNode.CafeDevTokens()
	if len(tokens) > 0 {
		t.Error("token database not updated (should be zero length)")
	}
}

func TestTokens_Teardown(t *testing.T) {
	tokenNode.Stop()
	tokenNode = nil
	os.RemoveAll(tokenRepoPath)
}
