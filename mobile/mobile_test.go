package mobile_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/textileio/textile-go/core"
	. "github.com/textileio/textile-go/mobile"
)

type TestMessenger struct {
	Messenger
}

func (tm *TestMessenger) Notify(event *Event) {}

var repo = "testdata/.textile"

var recovery string
var seed string

var mobile *Mobile
var defaultThreadId string
var threadId, threadId2 string
var addedPhotoId, addedBlockId string
var sharedBlockId string
var addedPhotoKey string

//var noteId string

func TestNewWallet(t *testing.T) {
	var err error
	recovery, err = NewWallet(12)
	if err != nil {
		t.Errorf("new mobile wallet failed: %s", err)
	}
}

func TestWalletAccountAt(t *testing.T) {
	res, err := WalletAccountAt(recovery, 0, "")
	if err != nil {
		t.Errorf("get mobile wallet account at failed: %s", err)
	}
	accnt := WalletAccount{}
	if err := json.Unmarshal([]byte(res), &accnt); err != nil {
		t.Error(err)
		return
	}
	seed = accnt.Seed
}

func TestInitRepo(t *testing.T) {
	os.RemoveAll(repo)
	if err := InitRepo(&InitConfig{
		Seed:     seed,
		RepoPath: repo,
	}); err != nil {
		t.Errorf("init mobile repo failed: %s", err)
	}
}

func TestMigrateRepo(t *testing.T) {
	if err := MigrateRepo(&MigrateConfig{
		RepoPath: repo,
	}); err != nil {
		t.Errorf("migrate mobile repo failed: %s", err)
	}
}

func TestNewTextile(t *testing.T) {
	config := &RunConfig{
		RepoPath: repo,
	}
	var err error
	mobile, err = NewTextile(config, &TestMessenger{})
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestNewTextileAgain(t *testing.T) {
	config := &RunConfig{
		RepoPath: repo,
	}
	if _, err := NewTextile(config, &TestMessenger{}); err != nil {
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

func TestMobile_Address(t *testing.T) {
	if mobile.Address() == "" {
		t.Error("got bad address")
	}
}

func TestMobile_Seed(t *testing.T) {
	if mobile.Seed() == "" {
		t.Error("got bad seed")
	}
}

func TestMobile_AccountThread(t *testing.T) {
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
	if len(threads.Items) != 1 {
		t.Error("get threads bad result")
	}
}

func TestMobile_AddThread(t *testing.T) {
	itemStr, err := mobile.AddThread(ksuid.New().String(), "default")
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
	itemStr, err := mobile.AddThread(ksuid.New().String(), "another")
	if err != nil {
		t.Errorf("add another thread failed: %s", err)
		return
	}
	item := Thread{}
	if err := json.Unmarshal([]byte(itemStr), &item); err != nil {
		t.Error(err)
		return
	}
	blockId, err := mobile.RemoveThread(item.Id)
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

func TestMobile_AddFile(t *testing.T) {
	resStr, err := mobile.AddFile("../images/testdata/image.jpg", defaultThreadId)
	if err != nil {
		t.Errorf("add file failed: %s", err)
		return
	}
	res := core.Directory{}
	if err := json.Unmarshal([]byte(resStr), &res); err != nil {
		t.Error(err)
		return
	}
	if len(res) != 6 {
		t.Error("wrong number of files")
		return
	}
}

//func TestMobile_AddPhotoToThread(t *testing.T) {
//	id, err := mobile.AddPhotoToThread(addedPhotoId, addedPhotoKey, threadId, "")
//	if err != nil {
//		t.Errorf("add photo to thread failed: %s", err)
//		return
//	}
//	addedBlockId = id
//}
//
//func TestMobile_SharePhotoToThread(t *testing.T) {
//	itemStr, err := mobile.AddThread("test", "test", "TextilePhotos")
//	if err != nil {
//		t.Errorf("add test thread failed: %s", err)
//		return
//	}
//	item := Thread{}
//	if err := json.Unmarshal([]byte(itemStr), &item); err != nil {
//		t.Error(err)
//		return
//	}
//	threadId2 = item.Id
//	sharedBlockId, err = mobile.SharePhotoToThread(addedPhotoId, item.Id, "howdy")
//	if err != nil {
//		t.Errorf("share photo to thread failed: %s", err)
//	}
//}
//
//func TestMobile_IgnorePhoto(t *testing.T) {
//	if _, err := mobile.IgnorePhoto(sharedBlockId); err != nil {
//		t.Errorf("ignore photo failed: %s", err)
//		return
//	}
//	res, err := mobile.Photos("", -1, threadId2)
//	if err != nil {
//		t.Errorf("get photos failed: %s", err)
//		return
//	}
//	photos := Photos{}
//	if err := json.Unmarshal([]byte(res), &photos); err != nil {
//		t.Error(err)
//		return
//	}
//	if len(photos.Items) != 0 {
//		t.Errorf("ignore photo bad result")
//	}
//}
//
//func TestMobile_AddPhotoComment(t *testing.T) {
//	if _, err := mobile.AddPhotoComment(addedBlockId, "well, well, well"); err != nil {
//		t.Errorf("add photo comment failed: %s", err)
//	}
//}
//
//func TestMobile_AddPhotoLike(t *testing.T) {
//	if _, err := mobile.AddPhotoLike(addedBlockId); err != nil {
//		t.Errorf("add photo like failed: %s", err)
//	}
//}
//
//func TestMobile_Photos(t *testing.T) {
//	res, err := mobile.Photos("", -1, threadId)
//	if err != nil {
//		t.Errorf("get photos failed: %s", err)
//		return
//	}
//	photos := Photos{}
//	if err := json.Unmarshal([]byte(res), &photos); err != nil {
//		t.Error(err)
//		return
//	}
//	if len(photos.Items) != 1 {
//		t.Errorf("get photos bad result")
//	}
//	if len(photos.Items[0].Comments) != 1 {
//		t.Errorf("get photo comments bad result")
//	}
//	if len(photos.Items[0].Likes) != 1 {
//		t.Errorf("get photo likes bad result")
//	}
//}
//
//func TestMobile_PhotosBadThread(t *testing.T) {
//	if _, err := mobile.Photos("", -1, "empty"); err == nil {
//		t.Errorf("get photo blocks from bad thread should fail: %s", err)
//	}
//}

//func TestMobile_PhotoThreads(t *testing.T) {
//	res, err := mobile.PhotoThreads(addedPhotoId)
//	if err != nil {
//		t.Errorf("get photo threads failed: %s", err)
//		return
//	}
//	threads := Threads{}
//	if err := json.Unmarshal([]byte(res), &threads); err != nil {
//		t.Error(err)
//		return
//	}
//	if len(threads.Items) != 2 {
//		t.Error("get photo threads bad result")
//	}
//}
//
//func TestMobile_PhotoData(t *testing.T) {
//	res, err := mobile.PhotoData(addedPhotoId, "thumb")
//	if err != nil {
//		t.Errorf("get photo data failed: %s", err)
//		return
//	}
//	if len(res) == 0 {
//		t.Errorf("get photo data bad result")
//	}
//}
//
//func TestMobile_PhotoDataForMinWidth(t *testing.T) {
//	// test photo
//	res, err := mobile.PhotoDataForMinWidth(addedPhotoId, 2000)
//	if err != nil {
//		t.Errorf("get photo data for min width failed: %s", err)
//		return
//	}
//	if len(res) == 0 {
//		t.Errorf("get photo data for min width bad result")
//		return
//	}
//	width, err := getWidthDataUrl(res)
//	if err != nil {
//		t.Errorf("get width failed: %s", err)
//		return
//	}
//	if width != 1600 {
//		t.Errorf("get photo data for min width bad result")
//	}
//
//	// test medium
//	res, err = mobile.PhotoDataForMinWidth(addedPhotoId, 600)
//	if err != nil {
//		t.Errorf("get photo data for min width failed: %s", err)
//		return
//	}
//	if len(res) == 0 {
//		t.Errorf("get photo data for min width bad result")
//		return
//	}
//	width, err = getWidthDataUrl(res)
//	if err != nil {
//		t.Errorf("get width failed: %s", err)
//		return
//	}
//	if width != 800 {
//		t.Errorf("get photo data for min width bad result")
//	}
//
//	// test small
//	res, err = mobile.PhotoDataForMinWidth(addedPhotoId, 320)
//	if err != nil {
//		t.Errorf("get photo data for min width failed: %s", err)
//		return
//	}
//	if len(res) == 0 {
//		t.Errorf("get photo data for min width bad result")
//		return
//	}
//	width, err = getWidthDataUrl(res)
//	if err != nil {
//		t.Errorf("get width failed: %s", err)
//		return
//	}
//	if width != 320 {
//		t.Errorf("get photo data for min width bad result")
//	}
//
//	// test photo
//	res, err = mobile.PhotoDataForMinWidth(addedPhotoId, 80)
//	if err != nil {
//		t.Errorf("get photo data for min width failed: %s", err)
//		return
//	}
//	if len(res) == 0 {
//		t.Errorf("get photo data for min width bad result")
//		return
//	}
//	width, err = getWidthDataUrl(res)
//	if err != nil {
//		t.Errorf("get width failed: %s", err)
//		return
//	}
//	if width != 100 {
//		t.Errorf("get photo data for min width bad result")
//	}
//}
//
//func TestMobile_PhotoMetadata(t *testing.T) {
//	res, err := mobile.PhotoMetadata(addedPhotoId)
//	if err != nil {
//		t.Errorf("get meta data failed: %s", err)
//		return
//	}
//	if len(res) == 0 {
//		t.Errorf("get meta data bad result")
//	}
//}
//
//func TestMobile_PhotoKey(t *testing.T) {
//	res, err := mobile.PhotoKey(addedPhotoId)
//	if err != nil {
//		t.Errorf("get key failed: %s", err)
//		return
//	}
//	if len(res) == 0 {
//		t.Errorf("get key bad result")
//	}
//}

func TestMobile_Overview(t *testing.T) {
	<-mobile.OnlineCh()
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

func TestMobile_SetUsername(t *testing.T) {
	if err := mobile.SetUsername("boomer"); err != nil {
		t.Errorf("set username failed: %s", err)
		return
	}
}

func TestMobile_SetAvatar(t *testing.T) {
	if err := mobile.SetAvatar(addedPhotoId); err != nil {
		t.Errorf("set avatar id failed: %s", err)
		return
	}
}

func TestMobile_Profile(t *testing.T) {
	profs, err := mobile.Profile()
	if err != nil {
		t.Errorf("get profile failed: %s", err)
		return
	}
	prof := core.Profile{}
	if err := json.Unmarshal([]byte(profs), &prof); err != nil {
		t.Error(err)
		return
	}
}

//func TestMobile_Notifications(t *testing.T) {
//	res, err := mobile.Notifications("", -1)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	notes := Notifications{}
//	if err := json.Unmarshal([]byte(res), &notes); err != nil {
//		t.Error(err)
//		return
//	}
//	if len(notes.Items) != 1 {
//		t.Error("get notifications bad result")
//		return
//	}
//	noteId = notes.Items[0].Id
//}
//
//func TestMobile_CountUnreadNotifications(t *testing.T) {
//	if mobile.CountUnreadNotifications() != 1 {
//		t.Error("count unread notifications bad result")
//	}
//}
//
//func TestMobile_ReadNotification(t *testing.T) {
//	if err := mobile.ReadNotification(noteId); err != nil {
//		t.Error(err)
//	}
//	if mobile.CountUnreadNotifications() != 0 {
//		t.Error("read notification bad result")
//	}
//}

func TestMobile_ReadAllNotifications(t *testing.T) {
	if err := mobile.ReadAllNotifications(); err != nil {
		t.Error(err)
	}
	if mobile.CountUnreadNotifications() != 0 {
		t.Error("read all notifications bad result")
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

func Test_Teardown(t *testing.T) {
	mobile = nil
}

//func getWidthDataUrl(res string) (int, error) {
//	var img *ImageData
//	if err := json.Unmarshal([]byte(res), &img); err != nil {
//		return 0, err
//	}
//	url := strings.Replace(img.Url, "data:image/jpeg;base64,", "", 1)
//	data, err := libp2pc.ConfigDecodeKey(url)
//	if err != nil {
//		return 0, err
//	}
//	conf, err := jpeg.DecodeConfig(bytes.NewReader(data))
//	if err != nil {
//		return 0, err
//	}
//	return conf.Width, nil
//}
