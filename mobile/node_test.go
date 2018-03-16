package mobile_test

import (
	"testing"

	. "github.com/textileio/textile-go/mobile"
)

var textile *Node

func TestNewTextile(t *testing.T) {
	textile = NewTextile("testdata/.ipfs")
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
		t.Errorf("pin photo on mobile node failed: %s", err)
	}
	if hash != "QmNnKzbJzAX8mUu1uuvyGzxe7p3z75D6UvUmS8LD5tc5ek" {
		t.Errorf("pin photo on mobile node bad hash: %s", err)
	}
}

func TestNode_Stop(t *testing.T) {
	err := textile.Stop()
	if err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}