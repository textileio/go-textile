package cafe

import (
	"github.com/textileio/textile-go/repo"
	"os"
	"testing"
)

var session *repo.CafeSession
var blockHash = "QmbQ4K3vXNJ3DjCNdG2urCXs7BuHqWQG1iSjZ8fbnF8NMs"
var photoHash = "QmSUnsZi9rGvPZLWy2v5N7fNxUWVNnA5nmppoM96FbLqLp"

func TestPin_Setup(t *testing.T) {
	// TODO get token from cafe service
}

func TestPin_Pin(t *testing.T) {
	block, err := os.Open("testdata/" + blockHash)
	if err != nil {
		t.Error(err)
		return
	}
	defer block.Close()
	res, err := pin(block, session.Access, "application/octet-stream")
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
	if resp.Id == nil {
		t.Error("response should contain id")
		return
	}
	if *resp.Id != blockHash {
		t.Errorf("hashes do not match: %s, %s", *resp.Id, blockHash)
	}
}

func TestPin_PinArchive(t *testing.T) {
	archive, err := os.Open("testdata/" + photoHash + ".tar.gz")
	if err != nil {
		t.Error(err)
		return
	}
	defer archive.Close()
	res, err := pin(archive, session.Access, "application/gzip")
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
	if resp.Id == nil {
		t.Error("response should contain id")
		return
	}
	if *resp.Id != photoHash {
		t.Errorf("hashes do not match: %s, %s", *resp.Id, photoHash)
	}
}
