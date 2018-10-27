package core_test

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	. "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/repo"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var repoPath1 = "testdata/.textile1"
var node1 *Textile
var repoPath2 = "testdata/.textile2"
var node2 *Textile

var session *repo.CafeSession
var blockHash = "QmbQ4K3vXNJ3DjCNdG2urCXs7BuHqWQG1iSjZ8fbnF8NMs"
var photoHash = "QmSUnsZi9rGvPZLWy2v5N7fNxUWVNnA5nmppoM96FbLqLp"

var client = &http.Client{}

func pin(reader io.Reader, cType string, token string, addr string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v0/pin", addr)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return client.Do(req)
}

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

func TestPin_Setup(t *testing.T) {
	// start one node
	os.RemoveAll(repoPath1)
	accnt1 := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  *accnt1,
		RepoPath: repoPath1,
		LogLevel: logging.DEBUG,
	}); err != nil {
		t.Errorf("init node1 failed: %s", err)
		return
	}
	var err error
	node1, err = NewTextile(RunConfig{
		RepoPath: repoPath1,
		LogLevel: logging.DEBUG,
	})
	if err != nil {
		t.Errorf("create node1 failed: %s", err)
		return
	}
	node1.Start()

	// start another
	os.RemoveAll(repoPath2)
	accnt2 := keypair.Random()
	if err := InitRepo(InitConfig{
		Account:  *accnt2,
		RepoPath: repoPath2,
		LogLevel: logging.DEBUG,
	}); err != nil {
		t.Errorf("init node2 failed: %s", err)
		return
	}
	node2, err = NewTextile(RunConfig{
		RepoPath: repoPath2,
		LogLevel: logging.DEBUG,
	})
	if err != nil {
		t.Errorf("create node2 failed: %s", err)
		return
	}
	node2.StartHttpApi("0.0.0.0:5000")
	node2.Start()

	// wait for both
	<-node1.OnlineCh()
	<-node2.OnlineCh()

	// register cafe
	peerId2, err := node2.PeerId()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if err := node1.RegisterCafe(peerId2.Pretty()); err != nil {
		t.Errorf("register node1 w/ node2 failed: %s", err)
		return
	}

	// get sessions
	sessions, err := node1.ListCafeSessions()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if len(sessions) > 0 {
		session = &sessions[0]
	} else {
		t.Errorf("no active sessions")
	}
}

func TestPin_Pin(t *testing.T) {
	block, err := os.Open("testdata/" + blockHash)
	if err != nil {
		t.Error(err)
		return
	}
	defer block.Close()
	addr := "http://" + session.HttpAddr
	res, err := pin(block, "application/octet-stream", session.Access, addr)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &PinResponse{}
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

func TestPin_PinArchive(t *testing.T) {
	archive, err := os.Open("testdata/" + photoHash + ".tar.gz")
	if err != nil {
		t.Error(err)
		return
	}
	defer archive.Close()
	addr := "http://" + session.HttpAddr
	res, err := pin(archive, "application/gzip", session.Access, addr)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &PinResponse{}
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
