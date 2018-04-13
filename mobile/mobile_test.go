package mobile_test

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/textileio/textile-go/core"
	. "github.com/textileio/textile-go/mobile"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

var wrapper *Wrapper
var hash string

func TestNewTextile(t *testing.T) {
	var err error
	wrapper, err = NewNode("testdata/.ipfs")
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestWrapper_Start(t *testing.T) {
	err := wrapper.Start()
	if err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestWrapper_StartAgain(t *testing.T) {
	err := wrapper.Start()
	if err != nil {
		t.Errorf("attempt to start a running node failed: %s", err)
	}
}

func TestWrapper_IsDatastoreConfigured(t *testing.T) {
	conf := wrapper.IsDatastoreConfigured()
	if conf {
		t.Error("check if datastore is configured bad result")
	}
}

func TestWrapper_ConfigureDatastore(t *testing.T) {
	err := wrapper.ConfigureDatastore("")
	if err != nil {
		t.Errorf("configure datastore on mobile node failed: %s", err)
	}
}

func TestWrapper_IsDatastoreConfigured2(t *testing.T) {
	conf := wrapper.IsDatastoreConfigured()
	if !conf {
		t.Error("check if datastore is configured (2) bad result")
	}
}

func TestWrapper_AddPhoto(t *testing.T) {
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

func TestWrapper_GetPhotos(t *testing.T) {
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

func TestWrapper_GetFileBase64(t *testing.T) {
	res, err := wrapper.GetFileBase64(hash + "/thumb")
	if err != nil {
		t.Errorf("get photo base64 string failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo base64 string bad result")
	}
}

func TestWrapper_GetRecoveryPhrase(t *testing.T) {
	mnemonic, err := wrapper.GetRecoveryPhrase()
	if err != nil {
		t.Errorf("failed to create a new recovery phrase: %s", err)
		return
	}
	list := strings.Split(mnemonic, " ")
	if len(list) != 24 {
		t.Errorf("got bad mnemonic length: %c", len(list))
	}
}

func TestWrapper_GetPeerID(t *testing.T) {
	_, err := wrapper.GetPeerID()
	if err != nil {
		t.Errorf("get peer id failed: %s", err)
	}
}

func TestWrapper_PairDesktop(t *testing.T) {
	_, pk, err := libp2p.GenerateKeyPair(libp2p.RSA, 4096)
	if err != nil {
		t.Errorf("create rsa keypair failed: %s", err)
	}
	pb, err := pk.Bytes()
	if err != nil {
		t.Errorf("get rsa keypair bytes: %s", err)
	}
	ps := base64.StdEncoding.EncodeToString(pb)

	_, err = wrapper.PairDesktop(ps)
	if err != nil {
		t.Errorf("pair desktop failed: %s", err)
	}
}

func TestWrapper_Stop(t *testing.T) {
	err := wrapper.Stop()
	if err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(wrapper.RepoPath)
}
