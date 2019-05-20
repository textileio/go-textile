package core_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	. "github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

var repoPath1 = "testdata/.textile1"
var node1 *Textile
var repoPath2 = "testdata/.textile2"
var node2 *Textile

var session *pb.CafeSession
var blockHash = "QmbQ4K3vXNJ3DjCNdG2urCXs7BuHqWQG1iSjZ8fbnF8NMs"
var photoHash = "QmSUnsZi9rGvPZLWy2v5N7fNxUWVNnA5nmppoM96FbLqLp"

type pinResponse struct {
	Id    string `json:"id"`
	Error string `json:"error"`
}

func TestCafeApi_Setup(t *testing.T) {
	// start one node
	_ = os.RemoveAll(repoPath1)
	accnt1 := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  accnt1,
		RepoPath: repoPath1,
	}); err != nil {
		t.Errorf("init node1 failed: %s", err)
		return
	}
	var err error
	node1, err = NewTextile(RunConfig{
		RepoPath: repoPath1,
	})
	if err != nil {
		t.Errorf("create node1 failed: %s", err)
		return
	}
	err = node1.Start()
	if err != nil {
		t.Fatal(err)
	}

	// start another
	_ = os.RemoveAll(repoPath2)
	accnt2 := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:     accnt2,
		RepoPath:    repoPath2,
		CafeApiAddr: "127.0.0.1:5000",
		CafeURL:     "http://127.0.0.1:5000",
		CafeOpen:    true,
	}); err != nil {
		t.Errorf("init node2 failed: %s", err)
		return
	}
	node2, err = NewTextile(RunConfig{
		RepoPath: repoPath2,
	})
	if err != nil {
		t.Errorf("create node2 failed: %s", err)
		return
	}
	err = node2.Start()
	if err != nil {
		t.Fatal(err)
	}

	// wait for cafe to be online
	<-node2.OnlineCh()

	// create token on cafe
	token, err := node2.CreateCafeToken("", true)
	if err != nil {
		t.Error(fmt.Errorf("error creating cafe token: %s", err))
		return
	}

	ok, err := node2.ValidateCafeToken(token)
	if !ok || err != nil {
		t.Error(fmt.Errorf("error checking token: %s", err))
		return
	}

	// register with cafe
	_, err = node1.RegisterCafe("http://127.0.0.1:5000", token)
	if err != nil {
		t.Errorf("register node1 w/ node2 failed: %s", err)
		return
	}

	// get sessions
	sessions := node1.CafeSessions()
	if len(sessions.Items) > 0 {
		session = sessions.Items[0]
	} else {
		t.Errorf("no active sessions")
	}
}

func TestCafeApi_Pin(t *testing.T) {
	block, err := os.Open("testdata/" + blockHash)
	if err != nil {
		t.Error(err)
		return
	}
	defer block.Close()
	res, err := pin(block, "application/octet-stream", session.Access, session.Cafe.Url)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &pinResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Id == "" {
		t.Error("response should contain id")
		return
	}
	if resp.Id != blockHash {
		t.Errorf("hashes do not match: %s, %s", resp.Id, blockHash)
	}
}

func TestCafeApi_PinArchive(t *testing.T) {
	archive, err := os.Open("testdata/" + photoHash + ".tar.gz")
	if err != nil {
		t.Error(err)
		return
	}
	defer archive.Close()
	res, err := pin(archive, "application/gzip", session.Access, session.Cafe.Url)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &pinResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Id == "" {
		t.Error("response should contain id")
		return
	}
	if resp.Id != photoHash {
		t.Errorf("hashes do not match: %s, %s", resp.Id, photoHash)
	}
}

func TestCafeApi_Teardown(t *testing.T) {
	_ = node1.Stop()
	_ = node2.Stop()
	node1 = nil
	node2 = nil
	_ = os.RemoveAll(repoPath1)
	_ = os.RemoveAll(repoPath2)
}

func pin(reader io.Reader, cType string, token string, addr string) (*http.Response, error) {
	url := fmt.Sprintf("%s/cafe/v0/pin", addr)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	return client.Do(req)
}

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
