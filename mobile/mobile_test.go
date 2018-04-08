package mobile_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/textileio/textile-go/core"
	. "github.com/textileio/textile-go/mobile"
)

var wrapper *Wrapper
var hash string

func TestNewTextile(t *testing.T) {
	wrapper = NewTextile("testdata/.ipfs")
}

func TestMobile_Start(t *testing.T) {
	err := wrapper.Start()
	if err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestMobile_StartAgain(t *testing.T) {
	err := wrapper.Start()
	if err != nil {
		t.Errorf("attempt to start a running node failed: %s", err)
	}
}

func TestMobile_ConfigureDatastore(t *testing.T) {
	err := wrapper.ConfigureDatastore("")
	if err != nil {
		t.Errorf("configure datastore on mobile node failed: %s", err)
	}
}

func TestMobile_AddPhoto(t *testing.T) {
	mr, err := wrapper.AddPhoto("testdata/photo.jpg", "testdata/thumb.jpg")
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

func TestMobile_GetPhotos(t *testing.T) {
	res, err := wrapper.GetPhotos("", -1)
	if err != nil {
		t.Errorf("get photos failed: %s", err)
		return
	}
	list := core.PhotoList{}
	json.Unmarshal([]byte(res), &list)
	if len(list.Hashes) == 0 {
		t.Errorf("get photos bad result")
	}
	hash = list.Hashes[0]
}

func TestMobile_GetPhotoBase64String(t *testing.T) {
	res, err := wrapper.GetFileBase64(hash + "/thumb")
	if err != nil {
		t.Errorf("get photo base64 string failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo base64 string bad result")
	}
}

func TestMobile_GetRecoveryPhrase(t *testing.T) {
	mnemonic, err := wrapper.GetRecoveryPhrase()
	if err != nil {
		t.Errorf("failed to create a new recovery phrase: %s", err)
	}
	list := strings.Split(mnemonic, " ")
	if len(list) != 24 {
		t.Errorf("got bad mnemonic length: %c", len(list))
	}
}

func TestMobile_Stop(t *testing.T) {
	err := wrapper.Stop()
	if err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}
