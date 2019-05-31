package mobile

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	icid "github.com/ipfs/go-cid"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

var cafesTestVars = struct {
	mobilePath  string
	mobile      *Mobile
	cafePath    string
	cafe        *core.Textile
	cafeApiPort string
}{
	mobilePath:  "./testdata/.textile3",
	cafePath:    "./testdata/.textile4",
	cafeApiPort: "5000",
}

func TestMobile_SetupCafes(t *testing.T) {
	var err error
	cafesTestVars.mobile, err = createAndStartPeer(InitConfig{
		RepoPath: cafesTestVars.mobilePath,
		Debug:    true,
	}, true, &testHandler{}, &testMessenger{})
	if err != nil {
		t.Fatal(err)
	}

	cafesTestVars.cafe, err = core.CreateAndStartPeer(core.InitConfig{
		RepoPath:    cafesTestVars.cafePath,
		Debug:       true,
		SwarmPorts:  "4001",
		CafeApiAddr: "0.0.0.0:" + cafesTestVars.cafeApiPort,
		CafeURL:     "http://127.0.0.1:" + cafesTestVars.cafeApiPort,
		CafeOpen:    true,
	}, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_RegisterCafe(t *testing.T) {
	token, err := cafesTestVars.cafe.CreateCafeToken("", true)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cafesTestVars.cafe.ValidateCafeToken(token)
	if !ok || err != nil {
		t.Fatal(err)
	}

	// because local discovery almost _always_ fails initially, a backoff is
	// set and we fail to register until it's removed... this cheats around that.
	cafeID := cafesTestVars.cafe.Ipfs().Identity
	cafesTestVars.mobile.node.Ipfs().Peerstore.AddAddrs(
		cafeID, cafesTestVars.cafe.Ipfs().PeerHost.Addrs(), peerstore.PermanentAddrTTL)

	// register with cafe
	err = cafesTestVars.mobile.RegisterCafe(cafeID.Pretty(), token)
	if err != nil {
		t.Fatal(err)
	}

	// add some data
	err = addTestData(cafesTestVars.mobile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_HandleCafeRequests(t *testing.T) {
	// manually flush until empty
	err := flush(10)
	if err != nil {
		t.Fatal(err)
	}

	m := cafesTestVars.mobile
	c := cafesTestVars.cafe

	err = m.node.Datastore().CafeRequests().DeleteCompleteSyncGroups()
	if err != nil {
		t.Fatal(err)
	}

	// ensure all requests have been deleted
	cnt := m.node.Datastore().CafeRequests().Count(-1)
	if cnt != 0 {
		t.Fatalf("expected all requests to be handled, got %d", cnt)
	}

	// check if blocks are pinned
	var blocks []string
	var targets []string
	list := m.node.Blocks("", -1, "")
	printBlocks(list)
	for _, b := range list.Items {
		blocks = append(blocks, b.Id)
		if b.Type == pb.Block_FILES {
			targets = append(targets, b.Target)
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

	// check if targets are pinned
	missingTargetPins, err := ipfs.NotPinned(c.Ipfs(), targets)
	if err != nil {
		t.Fatal(err)
	}
	if len(missingTargetPins) != 0 {
		var strs []string
		for _, id := range missingTargetPins {
			strs = append(strs, id.Hash().B58String())
		}
		t.Fatalf("targets not pinned: %s", strings.Join(strs, ", "))
	}

	// try unpinning a target
	if len(targets) > 0 {
		dec, err := icid.Decode(targets[0])
		if err != nil {
			t.Fatal(err)
		}
		err = ipfs.UnpinCid(c.Ipfs(), dec, true)
		if err != nil {
			t.Fatal(err)
		}
		not, err := ipfs.NotPinned(c.Ipfs(), []string{targets[0]})
		if err != nil {
			t.Fatal(err)
		}
		if len(not) == 0 || not[0].Hash().B58String() != targets[0] {
			t.Fatal("target was not recursively unpinned")
		}
	}
}

func TestMobile_TeardownCafes(t *testing.T) {
	_ = cafesTestVars.mobile.Stop()
	_ = cafesTestVars.cafe.Stop()
	cafesTestVars.mobile = nil
	cafesTestVars.cafe = nil
}

func addTestData(m *Mobile) error {
	thrd, err := addTestThread(m, &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_MEDIA,
		},
		Type:    pb.Thread_PRIVATE,
		Sharing: pb.Thread_INVITE_ONLY,
	})
	if err != nil {
		return err
	}

	_, err = m.addFiles([]string{"../mill/testdata/image.jpeg"}, thrd.Id, "hi")
	if err != nil {
		return err
	}
	_, err = m.AddMessage(thrd.Id, "hi")
	if err != nil {
		return err
	}
	hash, err := m.addFiles([]string{"../mill/testdata/image.png"}, thrd.Id, "hi")
	if err != nil {
		return err
	}
	_, err = m.AddComment(hash.B58String(), "nice")
	if err != nil {
		return err
	}
	hash, err = m.addFiles([]string{"../mill/testdata/image.jpeg", "../mill/testdata/image.png"}, thrd.Id, "hi")
	if err != nil {
		return err
	}
	_, err = m.AddLike(hash.B58String())
	if err != nil {
		return err
	}
	_, err = m.AddMessage(thrd.Id, "bye")
	if err != nil {
		return err
	}

	return nil
}

/*
Handle the request queue.
  1. List some requests
  2. Write the HTTP request for each
  3. Handle them (set to pending, send to cafe)
  4. Delete failed (reties not handled here)
  5. Set successful to complete
*/
func flushCafeRequests(limit int) (int, error) {
	var count int
	res, err := cafesTestVars.mobile.CafeRequests(limit)
	if err != nil {
		return count, err
	}
	groups := new(pb.Strings)
	err = proto.Unmarshal(res, groups)
	if err != nil {
		return count, err
	}
	count = len(groups.Values)

	// write the req for each group
	reqs := make(map[string]*pb.CafeHTTPRequest)
	for _, g := range groups.Values {
		res, err = cafesTestVars.mobile.WriteCafeRequest(g)
		if err != nil {
			return count, err
		}
		req := new(pb.CafeHTTPRequest)
		err = proto.Unmarshal(res, req)
		if err != nil {
			return count, err
		}
		reqs[g] = req
	}

	// mark each as pending (new loops for clarity)
	for g := range reqs {
		err = cafesTestVars.mobile.CafeRequestPending(g)
		if err != nil {
			return count, err
		}
	}

	// handle each
	for g, req := range reqs {
		res, err := handleReq(req)
		if err != nil {
			return count, err
		}
		if res.StatusCode >= 300 {
			fmt.Printf("got bad status: %d\n", res.StatusCode)
			err = cafesTestVars.mobile.FailCafeRequest(g)
		} else {
			err = cafesTestVars.mobile.CompleteCafeRequest(g)
		}
		if err != nil {
			return count, err
		}
		res.Body.Close()
	}
	return count, nil
}

func flush(batchSize int) error {
	count, err := flushCafeRequests(batchSize)
	if err != nil {
		return err
	}
	if count > 0 {
		return flush(batchSize)
	}
	return nil
}

func printSyncGroupStatus(status *pb.CafeSyncGroupStatus) {
	fmt.Println(">>> " + status.Id)
	fmt.Println(fmt.Sprintf("num. pending: %d", status.NumPending))
	fmt.Println(fmt.Sprintf("num. complete: %d", status.NumComplete))
	fmt.Println(fmt.Sprintf("num. total: %d", status.NumTotal))
	fmt.Println(fmt.Sprintf("size pending: %d", status.SizePending))
	fmt.Println(fmt.Sprintf("size complete: %d", status.SizeComplete))
	fmt.Println(fmt.Sprintf("size total: %d", status.SizeTotal))
	fmt.Println("<<<")
}

func handleReq(r *pb.CafeHTTPRequest) (*http.Response, error) {
	f, err := os.Open(r.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	req, err := http.NewRequest(r.Type.String(), r.Url, f)
	if err != nil {
		return nil, err
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	return client.Do(req)
}

func printBlocks(msg proto.Message) {
	marshaler := jsonpb.Marshaler{
		OrigName: true,
	}
	str, _ := marshaler.MarshalToString(msg)
	fmt.Println(str)
}
