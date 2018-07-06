package mobile_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/segmentio/ksuid"
	. "github.com/textileio/textile-go/mobile"
	util "github.com/textileio/textile-go/util/testing"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

type TestMessenger struct {
	Messenger
}

func (tm *TestMessenger) Notify(event *Event) {}

var repo = "testdata/.textile"

var mobile *Mobile
var addedPhotoId string
var sharedBlockId string

var cusername = ksuid.New().String()
var cpassword = ksuid.New().String()
var cemail = ksuid.New().String() + "@textile.io"

func TestNewTextile(t *testing.T) {
	os.RemoveAll(repo)
	config := &NodeConfig{
		RepoPath:      repo,
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
	_, ref, err := util.CreateReferral(util.RefKey, 1, 1, "test")
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
	}
}

func TestMobile_SignIn(t *testing.T) {
	if err := mobile.SignIn(cusername, cpassword); err != nil {
		t.Errorf("signin failed: %s", err)
	}
}

func TestMobile_IsSignedIn(t *testing.T) {
	if !mobile.IsSignedIn() {
		t.Errorf("is signed in check failed should be true")
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
	if _, err := mobile.GetAccessToken(); err != nil {
		t.Errorf("get access token failed: %s", err)
	}
}

func TestMobile_AddThread(t *testing.T) {
	item, err := mobile.AddThread("default", "")
	if err != nil {
		t.Errorf("add thread failed: %s", err)
	}
	if item == "" {
		t.Error("add thread bad result")
	}
}

func TestMobile_AddThreadAgain(t *testing.T) {
	if _, err := mobile.AddThread("default", ""); err == nil {
		t.Errorf("add thread again should fail: %s", err)
	}
}

func TestMobile_Threads(t *testing.T) {
	if _, err := mobile.AddThread("another", ""); err != nil {
		t.Errorf("add another thread failed: %s", err)
		return
	}
	res, err := mobile.Threads()
	if err != nil {
		t.Errorf("get threads failed: %s", err)
		return
	}
	threads := Threads{}
	json.Unmarshal([]byte(res), &threads)
	if len(threads.Items) != 2 {
		t.Error("get threads bad result")
	}
}

func TestMobile_RemoveThread(t *testing.T) {
	if err := mobile.RemoveThread("another"); err != nil {
		t.Errorf("remove thread failed: %s", err)
	}
}

func TestMobile_AddDevice(t *testing.T) {
	<-mobile.Online
	_, pk, err := libp2pc.GenerateKeyPair(libp2pc.Ed25519, 0)
	if err != nil {
		t.Error(err)
		return
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
		return
	}
	if err := mobile.AddDevice("hello", libp2pc.ConfigEncodeKey(pkb)); err != nil {
		t.Errorf("add device failed: %s", err)
	}
}

func TestMobile_AddDeviceAgain(t *testing.T) {
	_, pk, err := libp2pc.GenerateKeyPair(libp2pc.Ed25519, 0)
	if err != nil {
		t.Error(err)
		return
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
		return
	}
	if err := mobile.AddDevice("hello", libp2pc.ConfigEncodeKey(pkb)); err == nil {
		t.Error("add same device again should fail")
	}
}

func TestMobile_Devices(t *testing.T) {
	_, pk, err := libp2pc.GenerateKeyPair(libp2pc.Ed25519, 0)
	if err != nil {
		t.Error(err)
		return
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
		return
	}
	if err := mobile.AddDevice("another", libp2pc.ConfigEncodeKey(pkb)); err != nil {
		t.Errorf("add another device failed: %s", err)
	}
	res, err := mobile.Devices()
	if err != nil {
		t.Errorf("get devices failed: %s", err)
		return
	}
	devices := Devices{}
	json.Unmarshal([]byte(res), &devices)
	if len(devices.Items) != 2 {
		t.Error("get devices bad result")
	}
}

func TestMobile_RemoveDevice(t *testing.T) {
	if err := mobile.RemoveDevice("another"); err != nil {
		t.Errorf("remove device failed: %s", err)
	}
}

func TestMobile_AddPhoto(t *testing.T) {
	mrs, err := mobile.AddPhoto("testdata/image.jpg", "default", "howdy")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	reqs := PinRequests{}
	json.Unmarshal([]byte(mrs), &reqs)
	if len(reqs.Items) != 2 {
		t.Errorf("add photo got bad pin requests")
		return
	}
	addedPhotoId = reqs.Items[0].Boundary
}

func TestMobile_SharePhoto(t *testing.T) {
	if _, err := mobile.AddThread("test", ""); err != nil {
		t.Errorf("add test thread failed: %s", err)
		return
	}
	caption := "rasputin's eyes"
	mrs, err := mobile.SharePhoto(addedPhotoId, "test", caption)
	if err != nil {
		t.Errorf("share photo failed: %s", err)
		return
	}
	reqs := PinRequests{}
	json.Unmarshal([]byte(mrs), &reqs)
	if len(reqs.Items) != 1 {
		t.Errorf("share photo got bad pin requests")
		return
	}
	sharedBlockId = reqs.Items[0].Boundary
}

func TestMobile_PhotoBlocks(t *testing.T) {
	res, err := mobile.PhotoBlocks("", -1, "default")
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

func TestMobile_PhotosBadThread(t *testing.T) {
	_, err := mobile.PhotoBlocks("", -1, "empty")
	if err == nil {
		t.Errorf("get photo blocks from bad thread should fail: %s", err)
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
	}
}

func TestMobile_IsSignedInAgain(t *testing.T) {
	if mobile.IsSignedIn() {
		t.Errorf("is signed in check failed should be false")
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
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(mobile.RepoPath)
}
