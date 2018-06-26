package mobile_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/segmentio/ksuid"
	. "github.com/textileio/textile-go/mobile"
	util "github.com/textileio/textile-go/util/testing"
)

type TestMessenger struct {
	Messenger
}

func (tm *TestMessenger) Notify(event *Event) {}

var wrapper *Wrapper
var addedPhotoId string
var sharedBlockId string

var cusername = ksuid.New().String()
var cpassword = ksuid.New().String()
var cemail = ksuid.New().String() + "@textile.io"

func TestNewTextile(t *testing.T) {
	os.RemoveAll("testdata/.ipfs")
	config := &NodeConfig{
		RepoPath:      "testdata/.ipfs",
		CentralApiURL: util.CentralApiURL,
		LogLevel:      "DEBUG",
	}
	var err error
	wrapper, err = NewNode(config, &TestMessenger{})
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestWrapper_Start(t *testing.T) {
	if err := wrapper.Start(); err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestWrapper_StartAgain(t *testing.T) {
	if err := wrapper.Start(); err != nil {
		t.Errorf("attempt to start a running node failed: %s", err)
	}
}

func TestWrapper_SignUpWithEmail(t *testing.T) {
	_, ref, err := util.CreateReferral(util.RefKey, 1, 1)
	if err != nil {
		t.Errorf("create referral for signup failed: %s", err)
		return
	}
	if len(ref.RefCodes) == 0 {
		t.Error("create referral for signup got no codes")
		return
	}
	err = wrapper.SignUpWithEmail(cusername, cpassword, cemail, ref.RefCodes[0])
	if err != nil {
		t.Errorf("signup failed: %s", err)
		return
	}
}

func TestWrapper_SignIn(t *testing.T) {
	if err := wrapper.SignIn(cusername, cpassword); err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func TestWrapper_IsSignedIn(t *testing.T) {
	if !wrapper.IsSignedIn() {
		t.Errorf("is signed in check failed should be true")
		return
	}
}

func TestWrapper_GetId(t *testing.T) {
	id, err := wrapper.GetId()
	if err != nil {
		t.Errorf("get id failed: %s", err)
		return
	}
	if id == "" {
		t.Error("got bad id")
	}
}

func TestWrapper_GetIPFSPeerId(t *testing.T) {
	_, err := wrapper.GetIPFSPeerId()
	if err != nil {
		t.Errorf("get peer id failed: %s", err)
	}
}

func TestWrapper_GetUsername(t *testing.T) {
	un, err := wrapper.GetUsername()
	if err != nil {
		t.Errorf("get username failed: %s", err)
		return
	}
	if un != cusername {
		t.Errorf("got bad username: %s", un)
	}
}

func TestWrapper_GetAccessToken(t *testing.T) {
	_, err := wrapper.GetAccessToken()
	if err != nil {
		t.Errorf("get access token failed: %s", err)
		return
	}
}

func TestWrapper_AddThread(t *testing.T) {
	if err := wrapper.AddThread("default", ""); err != nil {
		t.Errorf("add thread failed: %s", err)
	}
}

func TestWrapper_AddThreadAgain(t *testing.T) {
	if err := wrapper.AddThread("default", ""); err != nil {
		t.Errorf("add thread again failed: %s", err)
	}
}

func TestWrapper_AddPhoto(t *testing.T) {
	mr, err := wrapper.AddPhoto("testdata/image.jpg", "default", "howdy")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	if len(mr.Boundary) == 0 {
		t.Errorf("add photo got bad hash")
	}
	addedPhotoId = mr.Boundary
	err = os.Remove("testdata/.ipfs/tmp/" + mr.Boundary)
	if err != nil {
		t.Errorf("error unlinking test multipart file: %s", err)
	}
}

func TestWrapper_SharePhoto(t *testing.T) {
	err := wrapper.AddThread("test", "")
	if err != nil {
		t.Errorf("add test thread failed: %s", err)
		return
	}
	caption := "rasputin's eyes"
	sharedBlockId, err = wrapper.SharePhoto(addedPhotoId, "test", caption)
	if err != nil {
		t.Errorf("share photo failed: %s", err)
		return
	}
	if len(sharedBlockId) == 0 {
		t.Errorf("share photo got bad id")
	}
}

func TestWrapper_GetPhotoBlocks(t *testing.T) {
	res, err := wrapper.GetPhotoBlocks("", -1, "default")
	if err != nil {
		t.Errorf("get photo blocks failed: %s", err)
		return
	}
	blocks := Blocks{}
	json.Unmarshal([]byte(res), &blocks)
	if len(blocks.Items) == 0 {
		t.Errorf("get photo blocks bad result")
	}
}

func TestWrapper_GetPhotosBadThread(t *testing.T) {
	_, err := wrapper.GetPhotoBlocks("", -1, "empty")
	if err == nil {
		t.Errorf("get photo blocks from bad thread should fail: %s", err)
		return
	}
}

func TestWrapper_GetBlockData(t *testing.T) {
	res, err := wrapper.GetBlockData(sharedBlockId, "caption")
	if err != nil {
		t.Errorf("get block data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get block data bad result")
	}
}

func TestWrapper_GetFileData(t *testing.T) {
	res, err := wrapper.GetFileData(addedPhotoId, "thumb")
	if err != nil {
		t.Errorf("get file data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get file data bad result")
	}
}

//func TestWrapper_PairDevice(t *testing.T) {
//	_, pk, err := libp2p.GenerateKeyPair(libp2p.Ed25519, 1024)
//	if err != nil {
//		t.Errorf("create keypair failed: %s", err)
//	}
//	pb, err := pk.Bytes()
//	if err != nil {
//		t.Errorf("get keypair bytes: %s", err)
//	}
//	ps := base64.StdEncoding.EncodeToString(pb)
//
//	_, err = wrapper.PairDevice(ps)
//	if err != nil {
//		t.Errorf("pair device failed: %s", err)
//	}
//}

func TestWrapper_SignOut(t *testing.T) {
	if err := wrapper.SignOut(); err != nil {
		t.Errorf("signout failed: %s", err)
		return
	}
}

func TestWrapper_IsSignedInAgain(t *testing.T) {
	if wrapper.IsSignedIn() {
		t.Errorf("is signed in check failed should be false")
		return
	}
}

func TestWrapper_Stop(t *testing.T) {
	if err := wrapper.Stop(); err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}

func TestWrapper_StopAgain(t *testing.T) {
	if err := wrapper.Stop(); err != nil {
		t.Errorf("stop mobile node again should not return error: %s", err)
	}
}

// test signin in stopped state, should re-connect to db
func TestWrapper_SignInAgain(t *testing.T) {
	if err := wrapper.SignIn(cusername, cpassword); err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(wrapper.RepoPath)
}
