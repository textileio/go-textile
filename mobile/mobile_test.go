package mobile_test

import (
	"crypto/rand"
	"encoding/json"
	"github.com/segmentio/ksuid"
	. "github.com/textileio/textile-go/mobile"
	"github.com/textileio/textile-go/net/model"
	util "github.com/textileio/textile-go/util/testing"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"os"
	"testing"
)

type TestMessenger struct {
	Messenger
}

func (tm *TestMessenger) Notify(event *Event) {}

var repo = "testdata/.textile"

var mobile *Mobile
var defaultThreadId string
var threadId string
var addedPhotoId string
var addedPhotoKey string
var sharedBlockId string
var deviceId string

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

func TestMobile_EmptyThreads(t *testing.T) {
	res, err := mobile.Threads()
	if err != nil {
		t.Errorf("get threads failed: %s", err)
		return
	}
	threads := Threads{}
	err = json.Unmarshal([]byte(res), &threads)
	if err != nil {
		t.Error(err)
		return
	}
	if len(threads.Items) != 0 {
		t.Error("get threads bad result")
	}
}

func TestMobile_AddThread(t *testing.T) {
	itemStr, err := mobile.AddThread("default", "")
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	item := Thread{}
	err = json.Unmarshal([]byte(itemStr), &item)
	if err != nil {
		t.Error(err)
		return
	}
	defaultThreadId = item.Id
}

func TestMobile_Threads(t *testing.T) {
	itemStr, err := mobile.AddThread("another", "")
	if err != nil {
		t.Errorf("add another thread failed: %s", err)
		return
	}
	item := Thread{}
	err = json.Unmarshal([]byte(itemStr), &item)
	if err != nil {
		t.Error(err)
		return
	}
	threadId = item.Id
	res, err := mobile.Threads()
	if err != nil {
		t.Errorf("get threads failed: %s", err)
		return
	}
	threads := Threads{}
	err = json.Unmarshal([]byte(res), &threads)
	if err != nil {
		t.Error(err)
		return
	}
	if len(threads.Items) != 2 {
		t.Error("get threads bad result")
	}
}

func TestMobile_RemoveThread(t *testing.T) {
	<-mobile.Online
	blockId, err := mobile.RemoveThread(defaultThreadId)
	if err != nil {
		t.Error(err)
		return
	}
	if blockId == "" {
		t.Errorf("remove thread bad result: %s", err)
	}
	if err != nil {
		t.Errorf("remove thread failed: %s", err)
	}
}

func TestMobile_AddDevice(t *testing.T) {
	_, pk, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
		return
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
		return
	}
	deviceId = libp2pc.ConfigEncodeKey(pkb)
	if err := mobile.AddDevice("hello", deviceId); err != nil {
		t.Errorf("add device failed: %s", err)
	}
}

func TestMobile_AddDeviceAgain(t *testing.T) {
	if err := mobile.AddDevice("hello", deviceId); err == nil {
		t.Error("add same device again should fail")
	}
}

func TestMobile_Devices(t *testing.T) {
	_, pk, err := libp2pc.GenerateEd25519Key(rand.Reader)
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
	err = json.Unmarshal([]byte(res), &devices)
	if err != nil {
		t.Error(err)
		return
	}
	if len(devices.Items) != 2 {
		t.Error("get devices bad result")
	}
}

func TestMobile_RemoveDevice(t *testing.T) {
	if err := mobile.RemoveDevice(deviceId); err != nil {
		t.Errorf("remove device failed: %s", err)
	}
}

func TestMobile_AddPhoto(t *testing.T) {
	resStr, err := mobile.AddPhoto("testdata/image.jpg")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	res := model.AddResult{}
	err = json.Unmarshal([]byte(resStr), &res)
	if err != nil {
		t.Error(err)
		return
	}
	addedPhotoKey = res.Key
	addedPhotoId = res.Id
}

func TestMobile_AddPhotoToThread(t *testing.T) {
	blockId, err := mobile.AddPhotoToThread(addedPhotoId, addedPhotoKey, threadId, "")
	if err != nil {
		t.Errorf("add photo to thread failed: %s", err)
		return
	}
	sharedBlockId = blockId
}

func TestMobile_SharePhotoToThread(t *testing.T) {
	itemStr, err := mobile.AddThread("test", "")
	if err != nil {
		t.Errorf("add test thread failed: %s", err)
		return
	}
	item := Thread{}
	err = json.Unmarshal([]byte(itemStr), &item)
	if err != nil {
		t.Error(err)
		return
	}
	id, err := mobile.SharePhotoToThread(addedPhotoId, item.Id, "howdy")
	if err != nil {
		t.Errorf("share photo to thread failed: %s", err)
		return
	}
	sharedBlockId = id
}

func TestMobile_GetPhotos(t *testing.T) {
	res, err := mobile.GetPhotos("", -1, threadId)
	if err != nil {
		t.Errorf("get photos failed: %s", err)
		return
	}
	photos := Photos{}
	err = json.Unmarshal([]byte(res), &photos)
	if err != nil {
		t.Error(err)
		return
	}
	if len(photos.Items) != 1 {
		t.Errorf("get photos bad result")
	}
}

func TestMobile_PhotosBadThread(t *testing.T) {
	_, err := mobile.GetPhotos("", -1, "empty")
	if err == nil {
		t.Errorf("get photo blocks from bad thread should fail: %s", err)
	}
}

func TestMobile_GetPhotoData(t *testing.T) {
	res, err := mobile.GetPhotoData(addedPhotoId)
	if err != nil {
		t.Errorf("get photo data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo data bad result")
	}
}

func TestMobile_GetThumbData(t *testing.T) {
	res, err := mobile.GetPhotoData(addedPhotoId)
	if err != nil {
		t.Errorf("get thumb data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get thumb data bad result")
	}
}

func TestMobile_GetPhotoMetadata(t *testing.T) {
	res, err := mobile.GetPhotoMetadata(addedPhotoId)
	if err != nil {
		t.Errorf("get meta data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get meta data bad result")
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
