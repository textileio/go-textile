package mobile_test

import (
	"testing"

	. "github.com/textileio/textile-go/mobile"
)

var textile *Node

func TestNewTextile(t *testing.T) {
	textile = NewTextile("testdata/.ipfs", "")
}

func TestNode_Start(t *testing.T) {
	err := textile.Start()
	if err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestNode_PinPhoto(t *testing.T) {
	hash, err := textile.PinPhoto("testdata/test.jpg")
	if err != nil {
		t.Errorf("pin photo failed: %s", err)
		return
	}
	if hash != "QmQzKq4hy8mTiiZGVfsW3PT95qbczCj2j5GWfZJV7hH2eu" {
		t.Errorf("pin photo got bad hash: %s", hash)
	}
}

func TestNode_GetPhotos(t *testing.T) {
	res, err := textile.GetPhotos("", -1)
	if err != nil {
		t.Errorf("get photos failed: %s", err)
		return
	}
	if res != "[QmQzKq4hy8mTiiZGVfsW3PT95qbczCj2j5GWfZJV7hH2eu]" {
		t.Errorf("get photos bad result")
	}
}

func TestNode_Stop(t *testing.T) {
	err := textile.Stop()
	if err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}