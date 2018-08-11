package mobile_test

import (
	"crypto/rand"
	"encoding/json"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/core"
	. "github.com/textileio/textile-go/mobile"
	util "github.com/textileio/textile-go/util/testing"
	"github.com/textileio/textile-go/wallet"
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
var deviceId string

var cusername = ksuid.New().String()
var cpassword = ksuid.New().String()
var cemail = ksuid.New().String() + "@textile.io"

func TestNewTextile(t *testing.T) {
	os.RemoveAll(repo)
	config := &NodeConfig{
		RepoPath: repo,
		CafeAddr: util.CafeAddr,
		LogLevel: "DEBUG",
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
	ref, err := util.CreateReferral(util.CafeReferralKey, 1, 1, "test")
	if err != nil {
		t.Errorf("create referral for signup failed: %s", err)
		return
	}
	if len(ref.RefCodes) == 0 {
		t.Error("create referral for signup got no codes")
		return
	}
	if err := mobile.SignUpWithEmail(cemail, cusername, cpassword, ref.RefCodes[0]); err != nil {
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

func TestMobile_GetTokens(t *testing.T) {
	if _, err := mobile.GetTokens(); err != nil {
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
	if err := json.Unmarshal([]byte(res), &threads); err != nil {
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
	if err := json.Unmarshal([]byte(itemStr), &item); err != nil {
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
	if err := json.Unmarshal([]byte(itemStr), &item); err != nil {
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
	if err := json.Unmarshal([]byte(res), &threads); err != nil {
		t.Error(err)
		return
	}
	if len(threads.Items) != 2 {
		t.Error("get threads bad result")
	}
}

func TestMobile_RemoveThread(t *testing.T) {
	<-core.Node.Wallet.Online()
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
	if err := json.Unmarshal([]byte(res), &devices); err != nil {
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
	resStr, err := mobile.AddPhoto("../util/testdata/image.jpg")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	res := wallet.AddDataResult{}
	if err := json.Unmarshal([]byte(resStr), &res); err != nil {
		t.Error(err)
		return
	}
	if res.Archive == nil {
		t.Error("add photo result should have an archive")
		return
	}
	addedPhotoKey = res.Key
	addedPhotoId = res.Id
}

func TestMobile_AddPhotoToThread(t *testing.T) {
	if _, err := mobile.AddPhotoToThread(addedPhotoId, addedPhotoKey, threadId, ""); err != nil {
		t.Errorf("add photo to thread failed: %s", err)
		return
	}
}

func TestMobile_SharePhotoToThread(t *testing.T) {
	itemStr, err := mobile.AddThread("test", "")
	if err != nil {
		t.Errorf("add test thread failed: %s", err)
		return
	}
	item := Thread{}
	if err := json.Unmarshal([]byte(itemStr), &item); err != nil {
		t.Error(err)
		return
	}
	if _, err := mobile.SharePhotoToThread(addedPhotoId, item.Id, "howdy"); err != nil {
		t.Errorf("share photo to thread failed: %s", err)
		return
	}
}

func TestMobile_GetPhotos(t *testing.T) {
	res, err := mobile.GetPhotos("", -1, threadId)
	if err != nil {
		t.Errorf("get photos failed: %s", err)
		return
	}
	photos := Photos{}
	if err := json.Unmarshal([]byte(res), &photos); err != nil {
		t.Error(err)
		return
	}
	if len(photos.Items) != 1 {
		t.Errorf("get photos bad result")
	}
}

func TestMobile_GetPhotosBadThread(t *testing.T) {
	if _, err := mobile.GetPhotos("", -1, "empty"); err == nil {
		t.Errorf("get photo blocks from bad thread should fail: %s", err)
	}
}

func TestMobile_PhotoThreads(t *testing.T) {
	res, err := mobile.PhotoThreads(addedPhotoId)
	if err != nil {
		t.Errorf("get photo threads failed: %s", err)
		return
	}
	threads := Threads{}
	if err := json.Unmarshal([]byte(res), &threads); err != nil {
		t.Error(err)
		return
	}
	if len(threads.Items) != 2 {
		t.Error("get photo threads bad result")
	}
}

func TestMobile_GetPhotoData(t *testing.T) {
	res, err := mobile.GetPhotoData(addedPhotoId, "thumb")
	if err != nil {
		t.Errorf("get photo data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo data bad result")
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

func TestMobile_GetPhotoKey(t *testing.T) {
	res, err := mobile.GetPhotoKey(addedPhotoId)
	if err != nil {
		t.Errorf("get key failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get key bad result")
	}
}

func TestMobile_SetAvatarId(t *testing.T) {
	if err := mobile.SetAvatarId(addedPhotoId); err != nil {
		t.Errorf("set avatar id failed: %s", err)
		return
	}
}

func TestMobile_GetProfile(t *testing.T) {
	profs, err := mobile.GetProfile()
	if err != nil {
		t.Errorf("get profile failed: %s", err)
		return
	}
	prof := wallet.Profile{}
	if err := json.Unmarshal([]byte(profs), &prof); err != nil {
		t.Error(err)
		return
	}
	if prof.Username != cusername {
		t.Errorf("get profile bad username result")
	}
	if len(prof.AvatarId) == 0 {
		t.Errorf("get profile bad avatar result")
	}
}

func TestMobile_GetStats(t *testing.T) {
	res, err := mobile.GetStats()
	if err != nil {
		t.Errorf("get stats failed: %s", err)
		return
	}
	stats := wallet.Stats{}
	if err := json.Unmarshal([]byte(res), &stats); err != nil {
		t.Error(err)
		return
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
