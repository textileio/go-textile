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

var mobile *Mobile
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
	mobile, err = NewNode(config, &TestMessenger{})
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestMobile_Start(t *testing.T) {
	if err := mobile.Start(); err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestMobile_StartAgain(t *testing.T) {
	if err := mobile.Start(); err != nil {
		t.Errorf("attempt to start a running node failed: %s", err)
	}
}

func TestMobile_SignUpWithEmail(t *testing.T) {
	_, ref, err := util.CreateReferral(util.RefKey, 1, 1, "TestMobile_SignUpWithEmail")
	if err != nil {
		t.Errorf("create referral for signup failed: %s", err)
		return
	}
	if len(ref.RefCodes) == 0 {
		t.Error("create referral for signup got no codes")
		return
	}
	err = mobile.SignUpWithEmail(cusername, cpassword, cemail, ref.RefCodes[0])
	if err != nil {
		t.Errorf("signup failed: %s", err)
		return
	}
}

func TestMobile_SignIn(t *testing.T) {
	if err := mobile.SignIn(cusername, cpassword); err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func TestMobile_IsSignedIn(t *testing.T) {
	if !mobile.IsSignedIn() {
		t.Errorf("is signed in check failed should be true")
		return
	}
}

func TestMobile_GetId(t *testing.T) {
	id, err := mobile.GetId()
	if err != nil {
		t.Errorf("get id failed: %s", err)
		return
	}
	if id == "" {
		t.Error("got bad id")
	}
}

func TestMobile_GetUsername(t *testing.T) {
	un, err := mobile.GetUsername()
	if err != nil {
		t.Errorf("get username failed: %s", err)
		return
	}
	if un != cusername {
		t.Errorf("got bad username: %s", un)
	}
}

func TestMobile_GetAccessToken(t *testing.T) {
	_, err := mobile.GetAccessToken()
	if err != nil {
		t.Errorf("get access token failed: %s", err)
		return
	}
}

func TestMobile_AddThread(t *testing.T) {
	if err := mobile.AddThread("default", ""); err != nil {
		t.Errorf("add thread failed: %s", err)
	}
}

func TestMobile_AddThreadAgain(t *testing.T) {
	if err := mobile.AddThread("default", ""); err != nil {
		t.Errorf("add thread again failed: %s", err)
	}
}

func TestMobile_AddPhoto(t *testing.T) {
	mr, err := mobile.AddPhoto("testdata/image.jpg", "default", "howdy")
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

func TestMobile_SharePhoto(t *testing.T) {
	err := mobile.AddThread("test", "")
	if err != nil {
		t.Errorf("add test thread failed: %s", err)
		return
	}
	caption := "rasputin's eyes"
	sharedBlockId, err = mobile.SharePhoto(addedPhotoId, "test", caption)
	if err != nil {
		t.Errorf("share photo failed: %s", err)
		return
	}
	if len(sharedBlockId) == 0 {
		t.Errorf("share photo got bad id")
	}
}

func TestMobile_GetPhotoBlocks(t *testing.T) {
	res, err := mobile.GetPhotoBlocks("", -1, "default")
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

func TestMobile_GetPhotosBadThread(t *testing.T) {
	_, err := mobile.GetPhotoBlocks("", -1, "empty")
	if err == nil {
		t.Errorf("get photo blocks from bad thread should fail: %s", err)
		return
	}
}

func TestMobile_GetBlockData(t *testing.T) {
	res, err := mobile.GetBlockData(sharedBlockId, "caption")
	if err != nil {
		t.Errorf("get block data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get block data bad result")
	}
}

func TestMobile_GetFileData(t *testing.T) {
	res, err := mobile.GetFileData(addedPhotoId, "thumb")
	if err != nil {
		t.Errorf("get file data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get file data bad result")
	}
}

func TestMobile_SignOut(t *testing.T) {
	if err := mobile.SignOut(); err != nil {
		t.Errorf("signout failed: %s", err)
		return
	}
}

func TestMobile_IsSignedInAgain(t *testing.T) {
	if mobile.IsSignedIn() {
		t.Errorf("is signed in check failed should be false")
		return
	}
}

func TestMobile_Stop(t *testing.T) {
	if err := mobile.Stop(); err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}

func TestMobile_StopAgain(t *testing.T) {
	if err := mobile.Stop(); err != nil {
		t.Errorf("stop mobile node again should not return error: %s", err)
	}
}

// test signin in stopped state, should re-connect to db
func TestMobile_SignInAgain(t *testing.T) {
	if err := mobile.SignIn(cusername, cpassword); err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(mobile.RepoPath)
}
