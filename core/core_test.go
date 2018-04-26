package core_test

import (
	"github.com/op/go-logging"
	. "github.com/textileio/textile-go/core"
	"os"
	"testing"
)

var node *TextileNode
var hash string

func TestNewNode(t *testing.T) {
	os.RemoveAll("testdata/.ipfs")
	var err error
	node, err = NewNode("testdata/.ipfs", false, logging.DEBUG)
	if err != nil {
		t.Errorf("create node failed: %s", err)
	}
}

func TestTextileNode_Start(t *testing.T) {
	err := node.Start()
	if err != nil {
		t.Errorf("start node failed: %s", err)
	}
}

func TestTextileNode_StartServices(t *testing.T) {
	_, err := node.StartServices()
	if err != nil {
		t.Errorf("start services failed: %s", err)
	}
}

func TestTextileNode_CreateAlbum(t *testing.T) {
	err := node.CreateAlbum("", "test")
	if err != nil {
		t.Errorf("create album failed: %s", err)
		return
	}
}

func TestTextileNode_AddPhoto(t *testing.T) {
	mr, err := node.AddPhoto("testdata/photo.jpg", "testdata/thumb.jpg", "default")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	if len(mr.Boundary) == 0 {
		t.Errorf("add photo got bad hash")
	}
	err = os.Remove("testdata/" + mr.Boundary)
	if err != nil {
		t.Errorf("error unlinking test multipart file: %s", err)
	}
}

func TestTextileNode_GetPhotos(t *testing.T) {
	list := node.GetPhotos("", -1, "default")
	if len(list.Hashes) == 0 {
		t.Errorf("get photos bad result")
	}
	hash = list.Hashes[0]
}

func TestTextileNode_GetFile(t *testing.T) {
	res, err := node.GetFile("/ipfs/"+hash+"/thumb", nil)
	if err != nil {
		t.Errorf("get photo failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo bad result")
	}
}

func TestTextileNode_GetPublicPeerKeyString(t *testing.T) {
	_, err := node.GetPublicPeerKeyString()
	if err != nil {
		t.Errorf("get peer public key as base 64 string failed: %s", err)
		return
	}
}

func TestTextileNode_Stop(t *testing.T) {
	err := node.Stop()
	if err != nil {
		t.Errorf("stop node failed: %s", err)
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(node.RepoPath)
}
