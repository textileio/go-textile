package mobile_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	. "github.com/textileio/textile-go/mobile"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"image/jpeg"
	"os"
	"strings"
	"testing"
)

type TestMessenger struct {
	Messenger
}

func (tm *TestMessenger) Notify(event *Event) {}

var repo = "testdata/.textile"

var mobile *Mobile
var defaultThreadId string
var threadId, threadId2 string
var addedPhotoId, addedBlockId string
var sharedBlockId string
var addedPhotoKey string
var deviceId string
var noteId string

func TestNewTextile(t *testing.T) {
	os.RemoveAll(repo)
	config := &NodeConfig{
		Account:  keypair.Random().Seed(),
		RepoPath: repo,
		CafeAddr: os.Getenv("CAFE_ADDR"),
		LogLevel: "DEBUG",
	}
	var err error
	mobile, err = NewNode(config, &TestMessenger{})
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestNewTextileAgain(t *testing.T) {
	config := &NodeConfig{
		RepoPath: repo,
		CafeAddr: os.Getenv("CAFE_ADDR"),
		LogLevel: "DEBUG",
	}
	if _, err := NewNode(config, &TestMessenger{}); err != nil {
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

func TestMobile_CafeRegister(t *testing.T) {
	req := &models.ReferralRequest{
		Key:         os.Getenv("CAFE_REFERRAL_KEY"),
		Count:       1,
		Limit:       1,
		RequestedBy: "test",
	}
	res, err := core.Node.CreateCafeReferral(req)
	if err != nil {
		t.Errorf("create referral for registration failed: %s", err)
		return
	}
	if len(res.RefCodes) == 0 {
		t.Error("create referral for registration got no codes")
		return
	}
	if err := mobile.CafeRegister(res.RefCodes[0]); err != nil {
		t.Errorf("register failed: %s", err)
	}
}

func TestMobile_CafeLogin(t *testing.T) {
	if err := mobile.CafeLogin(); err != nil {
		t.Errorf("login failed: %s", err)
	}
}

func TestMobile_CafeLoggedIn(t *testing.T) {
	if !mobile.CafeLoggedIn() {
		t.Errorf("check logged in failed, should be true")
	}
}

func TestMobile_GetID(t *testing.T) {
	id, err := mobile.GetID()
	if err != nil {
		t.Errorf("get id failed: %s", err)
		return
	}
	if id == "" {
		t.Error("got bad id")
	}
}

func TestMobile_GetAddress(t *testing.T) {
	id, err := mobile.GetAddress()
	if err != nil {
		t.Errorf("get address failed: %s", err)
		return
	}
	if id == "" {
		t.Error("got bad address")
	}
}

func TestMobile_GetSeed(t *testing.T) {
	id, err := mobile.GetSeed()
	if err != nil {
		t.Errorf("get seed failed: %s", err)
		return
	}
	if id == "" {
		t.Error("got bad seed")
	}
}

// TODO: set username
//func TestMobile_GetUsername(t *testing.T) {
//	un, err := mobile.GetUsername()
//	if err != nil {
//		t.Errorf("get username failed: %s", err)
//		return
//	}
//	if un != cusername {
//		t.Errorf("got bad username: %s", un)
//	}
//}

func TestMobile_GetCafeTokens(t *testing.T) {
	if _, err := mobile.GetCafeTokens(false); err != nil {
		t.Errorf("get cafe tokens failed: %s", err)
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
	itemStr, err := mobile.AddThread("default")
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
	itemStr, err := mobile.AddThread("another")
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
	<-core.Node.Online()
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
	resStr, err := mobile.AddPhoto("../photo/testdata/image.jpg")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	res := core.AddDataResult{}
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
	id, err := mobile.AddPhotoToThread(addedPhotoId, addedPhotoKey, threadId, "")
	if err != nil {
		t.Errorf("add photo to thread failed: %s", err)
		return
	}
	addedBlockId = id
}

func TestMobile_SharePhotoToThread(t *testing.T) {
	itemStr, err := mobile.AddThread("test")
	if err != nil {
		t.Errorf("add test thread failed: %s", err)
		return
	}
	item := Thread{}
	if err := json.Unmarshal([]byte(itemStr), &item); err != nil {
		t.Error(err)
		return
	}
	threadId2 = item.Id
	sharedBlockId, err = mobile.SharePhotoToThread(addedPhotoId, item.Id, "howdy")
	if err != nil {
		t.Errorf("share photo to thread failed: %s", err)
	}
}

func TestMobile_IgnorePhoto(t *testing.T) {
	if _, err := mobile.IgnorePhoto(sharedBlockId); err != nil {
		t.Errorf("ignore photo failed: %s", err)
		return
	}
	res, err := mobile.GetPhotos("", -1, threadId2)
	if err != nil {
		t.Errorf("get photos failed: %s", err)
		return
	}
	photos := Photos{}
	if err := json.Unmarshal([]byte(res), &photos); err != nil {
		t.Error(err)
		return
	}
	if len(photos.Items) != 0 {
		t.Errorf("ignore photo bad result")
	}
}

func TestMobile_AddPhotoComment(t *testing.T) {
	if _, err := mobile.AddPhotoComment(addedBlockId, "well, well, well"); err != nil {
		t.Errorf("add photo comment failed: %s", err)
	}
}

func TestMobile_AddPhotoLike(t *testing.T) {
	if _, err := mobile.AddPhotoLike(addedBlockId); err != nil {
		t.Errorf("add photo like failed: %s", err)
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
	if len(photos.Items[0].Comments) != 1 {
		t.Errorf("get photo comments bad result")
	}
	if len(photos.Items[0].Likes) != 1 {
		t.Errorf("get photo likes bad result")
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

func TestMobile_GetPhotoDataForMinWidth(t *testing.T) {
	// test photo
	res, err := mobile.GetPhotoDataForMinWidth(addedPhotoId, 2000)
	if err != nil {
		t.Errorf("get photo data for min width failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo data for min width bad result")
		return
	}
	width, err := getWidthDataUrl(res)
	if err != nil {
		t.Errorf("get width failed: %s", err)
		return
	}
	if width != 1600 {
		t.Errorf("get photo data for min width bad result")
	}

	// test medium
	res, err = mobile.GetPhotoDataForMinWidth(addedPhotoId, 600)
	if err != nil {
		t.Errorf("get photo data for min width failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo data for min width bad result")
		return
	}
	width, err = getWidthDataUrl(res)
	if err != nil {
		t.Errorf("get width failed: %s", err)
		return
	}
	if width != 800 {
		t.Errorf("get photo data for min width bad result")
	}

	// test small
	res, err = mobile.GetPhotoDataForMinWidth(addedPhotoId, 320)
	if err != nil {
		t.Errorf("get photo data for min width failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo data for min width bad result")
		return
	}
	width, err = getWidthDataUrl(res)
	if err != nil {
		t.Errorf("get width failed: %s", err)
		return
	}
	if width != 320 {
		t.Errorf("get photo data for min width bad result")
	}

	// test photo
	res, err = mobile.GetPhotoDataForMinWidth(addedPhotoId, 80)
	if err != nil {
		t.Errorf("get photo data for min width failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo data for min width bad result")
		return
	}
	width, err = getWidthDataUrl(res)
	if err != nil {
		t.Errorf("get width failed: %s", err)
		return
	}
	if width != 100 {
		t.Errorf("get photo data for min width bad result")
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

func TestMobile_SetAvatar(t *testing.T) {
	if err := mobile.SetAvatar(addedPhotoId); err != nil {
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
	prof := core.AccountProfile{}
	if err := json.Unmarshal([]byte(profs), &prof); err != nil {
		t.Error(err)
		return
	}
	//if prof.Username != cusername {
	//	t.Errorf("get profile bad username result")
	//}
	//if len(prof.AvatarId) == 0 {
	//	t.Errorf("get profile bad avatar result")
	//}
}

func TestMobile_Overview(t *testing.T) {
	res, err := mobile.Overview()
	if err != nil {
		t.Errorf("get overview failed: %s", err)
		return
	}
	stats := core.Overview{}
	if err := json.Unmarshal([]byte(res), &stats); err != nil {
		t.Error(err)
		return
	}
}

func TestMobile_GetNotifications(t *testing.T) {
	res, err := mobile.GetNotifications("", -1)
	if err != nil {
		t.Error(err)
		return
	}
	notes := Notifications{}
	if err := json.Unmarshal([]byte(res), &notes); err != nil {
		t.Error(err)
		return
	}
	if len(notes.Items) != 1 {
		t.Error("get notifications bad result")
		return
	}
	noteId = notes.Items[0].Id
}

func TestMobile_CountUnreadNotifications(t *testing.T) {
	if mobile.CountUnreadNotifications() != 1 {
		t.Error("count unread notifications bad result")
	}
}

func TestMobile_ReadNotification(t *testing.T) {
	if err := mobile.ReadNotification(noteId); err != nil {
		t.Error(err)
	}
	if mobile.CountUnreadNotifications() != 0 {
		t.Error("read notification bad result")
	}
}

func TestMobile_ReadAllNotifications(t *testing.T) {
	if err := mobile.ReadAllNotifications(); err != nil {
		t.Error(err)
	}
	if mobile.CountUnreadNotifications() != 0 {
		t.Error("read all notifications bad result")
	}
}

func TestMobile_CafeLogout(t *testing.T) {
	if err := mobile.CafeLogout(); err != nil {
		t.Errorf("logout failed: %s", err)
	}
}

func TestMobile_CafeLoggedInAgain(t *testing.T) {
	if mobile.CafeLoggedIn() {
		t.Errorf("check logged in failed, should be false")
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

// test login in stopped state, should re-connect to db
func TestMobile_CafeLoginAgain(t *testing.T) {
	if err := mobile.CafeLogin(); err != nil {
		t.Errorf("login again failed: %s", err)
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(mobile.RepoPath)
}

func getWidthDataUrl(res string) (int, error) {
	var img *ImageData
	if err := json.Unmarshal([]byte(res), &img); err != nil {
		return 0, err
	}
	url := strings.Replace(img.Url, "data:image/jpeg;base64,", "", 1)
	data, err := libp2pc.ConfigDecodeKey(url)
	if err != nil {
		return 0, err
	}
	conf, err := jpeg.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	return conf.Width, nil
}
