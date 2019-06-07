package core

import (
	"fmt"
	"testing"
	"time"

	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema/textile"
)

var cafeVars = struct {
	nodePath string
	cafePath string

	cafeApiPort string

	node *Textile
	cafe *Textile

	token string
}{
	nodePath:    "./testdata/.textile3",
	cafePath:    "./testdata/.textile4",
	cafeApiPort: "5000",
}

func TestCore_SetupCafes(t *testing.T) {
	var err error
	cafeVars.node, err = CreateAndStartPeer(InitConfig{
		RepoPath: cafeVars.nodePath,
		Debug:    true,
	}, true)
	if err != nil {
		t.Fatal(err)
	}

	cafeVars.cafe, err = CreateAndStartPeer(InitConfig{
		RepoPath:    cafeVars.cafePath,
		Debug:       true,
		SwarmPorts:  "4001",
		CafeApiAddr: "0.0.0.0:" + cafeVars.cafeApiPort,
		CafeURL:     "http://127.0.0.1:" + cafeVars.cafeApiPort,
		CafeOpen:    true,
	}, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTextile_CafeTokens(t *testing.T) {
	var err error
	cafeVars.token, err = cafeVars.cafe.CreateCafeToken("", true)
	if err != nil {
		t.Fatalf("error creating cafe token: %s", err)
	}
	if len(cafeVars.token) == 0 {
		t.Fatal("invalid token created")
	}

	tokens, _ := cafeVars.cafe.CafeTokens()
	if len(tokens) < 1 {
		t.Fatal("token database not updated (should be length 1)")
	}

	ok, err := cafeVars.cafe.ValidateCafeToken("blah")
	if err == nil || ok {
		t.Fatal("expected token comparison with 'blah' to be invalid")
	}

	ok, err = cafeVars.cafe.ValidateCafeToken(cafeVars.token)
	if err != nil || !ok {
		t.Fatal("expected token comparison to be valid")
	}
}

func TestTextile_RemoveCafeToken(t *testing.T) {
	err := cafeVars.cafe.RemoveCafeToken(cafeVars.token)
	if err != nil {
		t.Fatal("expected be remove token cleanly")
	}

	tokens, _ := cafeVars.cafe.CafeTokens()
	if len(tokens) > 0 {
		t.Fatal("token database not updated (should be zero length)")
	}
}

func TestCore_RegisterCafe(t *testing.T) {
	token, err := cafeVars.cafe.CreateCafeToken("", true)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cafeVars.cafe.ValidateCafeToken(token)
	if !ok || err != nil {
		t.Fatal(err)
	}

	// because local discovery almost _always_ fails initially, a backoff is
	// set and we fail to register until it's removed... this cheats around that.
	cafeID := cafeVars.cafe.Ipfs().Identity
	cafeVars.node.Ipfs().Peerstore.AddAddrs(
		cafeID, cafeVars.cafe.Ipfs().PeerHost.Addrs(), peerstore.PermanentAddrTTL)

	_, err = cafeVars.node.RegisterCafe(cafeID.Pretty(), token)
	if err != nil {
		t.Fatalf("register node1 w/ node2 failed: %s", err)
	}

	// add some data
	err = addTestData(cafeVars.node)
	if err != nil {
		t.Fatal(err)
	}
	cafeVars.node.FlushCafes()
}

func TestCore_HandleCafeRequests(t *testing.T) {
	waitOnRequests(time.Second * 60)

	// ensure all requests have been deleted
	total := cafeVars.node.datastore.CafeRequests().Count(-1)
	neww := cafeVars.node.datastore.CafeRequests().Count(0)
	if neww != 0 {
		t.Fatalf("expected all requests to be handled, got %d total, %d new", total, neww)
	}
}

func TestCore_TeardownCafes(t *testing.T) {
	_ = cafeVars.node.Stop()
	_ = cafeVars.cafe.Stop()
	cafeVars.node = nil
	cafeVars.cafe = nil
}

func addTestData(n *Textile) error {
	thrd, err := addTestThread(n, &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test",
		Schema: &pb.AddThreadConfig_Schema{
			Json: textile.Blob,
		},
		Type:    pb.Thread_PRIVATE,
		Sharing: pb.Thread_INVITE_ONLY,
	})
	if err != nil {
		return err
	}

	_, err = addData(n, []string{"../mill/testdata/image.jpeg"}, thrd, "hi")
	if err != nil {
		return err
	}
	_, err = thrd.AddMessage("hi")
	if err != nil {
		return err
	}
	files, err := addData(n, []string{"../mill/testdata/image.png"}, thrd, "hi")
	if err != nil {
		return err
	}
	_, err = thrd.AddComment(files.Block, "nice")
	if err != nil {
		return err
	}
	files, err = addData(n, []string{"../mill/testdata/image.jpeg", "../mill/testdata/image.png"}, thrd, "hi")
	if err != nil {
		return err
	}
	_, err = thrd.AddLike(files.Block)
	if err != nil {
		return err
	}
	_, err = thrd.AddMessage("bye")
	if err != nil {
		return err
	}

	return nil
}

func waitOnRequests(total time.Duration) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	var waited time.Duration
	for {
		select {
		case <-tick.C:
			cnt := cafeVars.node.datastore.CafeRequests().Count(-1)
			if cnt == 0 {
				return
			} else {
				fmt.Printf("waiting on %d requests to complete\n", cnt)
			}
			waited += time.Second
			if waited >= total {
				return
			}
		}
	}
}
