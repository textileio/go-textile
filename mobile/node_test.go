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
	textile = NewTextile("testdata/.ipfs", "")
	err := textile.Start()
	if err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestNode_PinPhoto(t *testing.T) {
	textile = NewTextile("testdata/.ipfs", "")

	errStart := textile.Start()
	if errStart != nil {
		t.Errorf("start mobile node failed: %s", errStart)
	}

	hash, err := textile.PinPhoto("testdata/test.jpg", "testdata/thumb.jpg")
	if err != nil {
		t.Errorf("pin photo on mobile node failed: %s", err)
		return
	}

	if hash != "QmXK1noVgYCFfAWDPFGRJyMeLSXYUw4K82HNsvmyhxWHH1" {
		t.Errorf("pin photo on mobile node bad hash: %s", hash)
	}
}

func TestNode_Stop(t *testing.T) {
	err := textile.Stop()
	if err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}