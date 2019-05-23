package mobile

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

var cafesTestVars = struct {
	cafePath    string
	cafe        *core.Textile
	cafeApiPort string
	mobilePath  string
	mobile      *Mobile
}{
	cafePath:    "./testdata/.textile3",
	cafeApiPort: "5000",
	mobilePath:  "./testdata/.textile4",
}

type testCafesHandler struct {
	mux sync.Mutex
}

/*
Handle the request queue.
  1. List some groups
  2. List get the HTTP request list for each of those groups
  3. Handle them
  4. Delete failed (reties not handled here)
  5. Mark successful as complete
*/
func (th *testCafesHandler) Flush() {
	th.mux.Lock()
	defer th.mux.Unlock()

	if cafesTestVars.mobile == nil {
		return
	}
	fmt.Println(">>> FLUSHING <<<")

	res, err := cafesTestVars.mobile.CafeRequests(10)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	groups := new(pb.Strings)
	err = proto.Unmarshal(res, groups)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// print status to start
	//printGroupStatus(groups)

	// write the reqs for each group
	reqs := make(map[string]*pb.CafeHTTPRequestList)
	for _, g := range groups.Values {
		res, err = cafesTestVars.mobile.WriteCafeHTTPRequests(g)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		list := new(pb.CafeHTTPRequestList)
		err = proto.Unmarshal(res, list)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		reqs[g] = list
	}

	// handle each (doing this in a new loop for test clarity)
	failed := map[string]struct{}{}
	for g, list := range reqs {
		for _, req := range list.Items {
			res, err := handleReq(req)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if res.StatusCode >= 300 {
				fmt.Printf("got bad status: %d\n", res.StatusCode)
				failed[g] = struct{}{}
			} else {
				err = cafesTestVars.mobile.SetCafeRequestComplete(g)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			res.Body.Close()
			//printGroupStatus(groups)
		}
	}

	// delete failed groups
	for group := range failed {
		err = cafesTestVars.mobile.SetCafeRequestFailed(group)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func TestMobile_SetupCafes(t *testing.T) {
	var err error
	cafesTestVars.mobile, err = createAndStartMobile(
		cafesTestVars.mobilePath, false, &testCafesHandler{}, &testMessenger{})
	if err != nil {
		t.Fatal(err)
	}

	// start a cafe
	_ = os.RemoveAll(cafesTestVars.cafePath)
	err = core.InitRepo(core.InitConfig{
		Account:     keypair.Random(),
		RepoPath:    cafesTestVars.cafePath,
		CafeApiAddr: "0.0.0.0:" + cafesTestVars.cafeApiPort,
		CafeOpen:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	cafesTestVars.cafe, err = core.NewTextile(core.RunConfig{
		RepoPath: cafesTestVars.cafePath,
		Debug:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = cafesTestVars.cafe.Start()
	if err != nil {
		t.Fatal(err)
	}

	<-cafesTestVars.mobile.OnlineCh()
	<-cafesTestVars.cafe.OnlineCh()
}

func TestMobile_RegisterCafe(t *testing.T) {
	// create a token
	token, err := cafesTestVars.cafe.CreateCafeToken("", true)
	if err != nil {
		t.Fatal(err)
	}

	// register with cafe
	url := "http://127.0.0.1:" + cafesTestVars.cafeApiPort
	err = cafesTestVars.mobile.RegisterCafe(url, token)
	if err != nil {
		t.Fatal(err)
	}

	// add some data
	err = addTestData(cafesTestVars.mobile)
	if err != nil {
		t.Fatal(err)
	}
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

func printGroupStatus(list *pb.Strings) {
	for _, g := range list.Values {
		res, err := cafesTestVars.mobile.CafeRequestGroupStatus(g)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		status := new(pb.CafeRequestSyncGroupStatus)
		err = proto.Unmarshal(res, status)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(">>> " + g)
		fmt.Println(fmt.Sprintf("num. pending: %d", status.NumPending))
		fmt.Println(fmt.Sprintf("num. complete: %d", status.NumComplete))
		fmt.Println(fmt.Sprintf("num. total: %d", status.NumTotal))
		fmt.Println(fmt.Sprintf("size pending: %d", status.SizePending))
		fmt.Println(fmt.Sprintf("size complete: %d", status.SizeComplete))
		fmt.Println(fmt.Sprintf("size total: %d", status.SizeTotal))
		fmt.Println("<<<")
	}
}
