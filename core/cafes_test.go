package core

import (
	"fmt"
	"strings"
	"testing"
	"time"

	icid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema/textile"
)

var cafeVars = struct {
	nodeInitConfig InitConfig
	cafeInitConfig InitConfig

	node *Textile
	cafe *Textile

	token string
}{
	nodeInitConfig: InitConfig{
		BaseRepoPath: "./testdata/.textile3",
		Debug:        true,
	},
	cafeInitConfig: InitConfig{
		BaseRepoPath: "./testdata/.textile4",
		Debug:        true,
		SwarmPorts:   "4001",
		CafeApiAddr:  "0.0.0.0:5000",
		CafeURL:      "http://127.0.0.1:5000",
		CafeOpen:     true,
	},
}

func TestCore_SetupCafes(t *testing.T) {
	var err error
	cafeVars.node, err = CreateAndStartPeer(cafeVars.nodeInitConfig, true)
	if err != nil {
		t.Fatal(err)
	}

	cafeVars.cafe, err = CreateAndStartPeer(cafeVars.cafeInitConfig, true)
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

	// register with cafe
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

	n := cafeVars.node
	c := cafeVars.cafe

	// ensure all requests have been deleted
	cnt := n.datastore.CafeRequests().Count(-1)
	ncnt := n.datastore.CafeRequests().Count(0)
	if ncnt != 0 {
		t.Fatalf("expected all requests to be handled, got %d total, %d new", cnt, ncnt)
	}

	// check if blocks are pinned
	var blocks []string
	var datas []string
	list := n.Blocks("", -1, "")
	for _, b := range list.Items {
		blocks = append(blocks, b.Id)
		if b.Type == pb.Block_FILES {
			datas = append(datas, b.Data)
		}
	}
	missingBlockPins, err := ipfs.NotPinned(c.Ipfs(), blocks)
	if err != nil {
		t.Fatal(err)
	}
	if len(missingBlockPins) != 0 {
		var strs []string
		for _, id := range missingBlockPins {
			strs = append(strs, id.Hash().B58String())
		}
		t.Fatalf("blocks not pinned: %s", strings.Join(strs, ", "))
	}

	// check if datas are pinned
	missingDataPins, err := ipfs.NotPinned(c.Ipfs(), datas)
	if err != nil {
		t.Fatal(err)
	}
	if len(missingDataPins) != 0 {
		var strs []string
		for _, id := range missingDataPins {
			strs = append(strs, id.Hash().B58String())
		}
		t.Fatalf("datas not pinned: %s", strings.Join(strs, ", "))
	}

	// try unpinning data
	if len(datas) > 0 {
		dec, err := icid.Decode(datas[0])
		if err != nil {
			t.Fatal(err)
		}
		err = ipfs.UnpinCid(c.Ipfs(), dec, true)
		if err != nil {
			t.Fatal(err)
		}
		not, err := ipfs.NotPinned(c.Ipfs(), []string{datas[0]})
		if err != nil {
			t.Fatal(err)
		}
		if len(not) == 0 || not[0].Hash().B58String() != datas[0] {
			t.Fatal("data was not recursively unpinned")
		}
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
	_, err = thrd.AddMessage("", "hi")
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
	_, err = thrd.AddMessage("", "bye")
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
