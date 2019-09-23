package mobile

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	icid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

var cafesTestVars = struct {
	mobilePath string
	mobile     *Mobile

	cafePath string
	cafe     *core.Textile
	token    string
}{
	mobilePath: "./testdata/.textile3",
	cafePath:   "./testdata/.textile4",
}

func TestMobile_SetupCafes(t *testing.T) {
	var err error
	cafesTestVars.mobile, err = createAndStartPeer(InitConfig{
		BaseRepoPath: cafesTestVars.mobilePath,
		Debug:        true,
	}, true, &testHandler{}, &testMessenger{})
	if err != nil {
		t.Fatal(err)
	}

	cafesTestVars.cafe, err = core.CreateAndStartPeer(core.InitConfig{
		BaseRepoPath: cafesTestVars.cafePath,
		Debug:        true,
		SwarmPorts:   "4001",
		CafeApiAddr:  "0.0.0.0:5000",
		CafeURL:      "http://127.0.0.1:5000",
		CafeOpen:     true,
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

	// register with cafe
	cafeID := cafesTestVars.cafe.Ipfs().Identity
	cafesTestVars.mobile.node.Ipfs().Peerstore.AddAddrs(
		cafeID, cafesTestVars.cafe.Ipfs().PeerHost.Addrs(), peerstore.PermanentAddrTTL)
	err = cafesTestVars.mobile.registerCafe("http://127.0.0.1:5000", token)
	if err != nil {
		t.Fatal(err)
	}

	cafesTestVars.token = token
}

func TestMobile_DeregisterCafe(t *testing.T) {
	res, err := cafesTestVars.mobile.CafeSessions()
	if err != nil {
		t.Fatal(err)
	}
	sessions := new(pb.CafeSessionList)
	err = proto.Unmarshal(res, sessions)
	if err != nil {
		t.Fatal(err)
	}

	err = cafesTestVars.mobile.deregisterCafe(sessions.Items[0].Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_ReregisterCafe(t *testing.T) {
	err := cafesTestVars.mobile.registerCafe("http://127.0.0.1:5000", cafesTestVars.token)
	if err != nil {
		t.Fatal(err)
	}

	// add some data
	err = addTestData(cafesTestVars.mobile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_RefreshCafeSession(t *testing.T) {
	res, err := cafesTestVars.mobile.CafeSessions()
	if err != nil {
		t.Fatal(err)
	}
	sessions := new(pb.CafeSessionList)
	err = proto.Unmarshal(res, sessions)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cafesTestVars.mobile.refreshCafeSession(sessions.Items[0].Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_CheckCafeMessages(t *testing.T) {
	err := cafesTestVars.mobile.checkCafeMessages()
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
	var datas []string
	list := m.node.Blocks("", -1, "")
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

func TestMobile_TeardownCafes(t *testing.T) {
	_ = cafesTestVars.mobile.stop()
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
	ids := new(pb.Strings)
	err = proto.Unmarshal(res, ids)
	if err != nil {
		return count, err
	}
	count = len(ids.Values)

	// write the req for each group
	reqs := make(map[string]*pb.CafeHTTPRequest)
	for _, g := range ids.Values {
		res, err = cafesTestVars.mobile.writeCafeRequest(g)
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
			reason := fmt.Sprintf("got bad status: %d\n", res.StatusCode)
			fmt.Println(reason)
			err = cafesTestVars.mobile.FailCafeRequest(g, reason)
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
