package mobile_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/textileio/textile-go/core"
	. "github.com/textileio/textile-go/mobile"
)

type TestMessenger struct {
	Messenger
}

func (tm *TestMessenger) Notify(event *Event) {}

var repoPath = "testdata/.textile"

var recovery string
var seed string

var mobile *Mobile
var defaultThreadId string
var prepared string
var filesBlock core.BlockInfo
var files []core.ThreadFilesInfo

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
	os.RemoveAll(repoPath)
	if err := InitRepo(&InitConfig{
		Seed:     seed,
		RepoPath: repoPath,
	}); err != nil {
		t.Errorf("init mobile repo failed: %s", err)
	}
}

func TestMigrateRepo(t *testing.T) {
	if err := MigrateRepo(&MigrateConfig{
		RepoPath: repoPath,
	}); err != nil {
		t.Errorf("migrate mobile repo failed: %s", err)
	}
}

func TestNewTextile(t *testing.T) {
	config := &RunConfig{
		RepoPath: repoPath,
	}
	var err error
	mobile, err = NewTextile(config, &TestMessenger{})
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestNewTextileAgain(t *testing.T) {
	config := &RunConfig{
		RepoPath: repoPath,
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

func TestMobile_CheckAccountThread(t *testing.T) {
	res, err := mobile.Threads()
	if err != nil {
		t.Errorf("get threads failed: %s", err)
		return
	}
	var threads []core.ThreadInfo
	if err := json.Unmarshal([]byte(res), &threads); err != nil {
		t.Error(err)
		return
	}
	if len(threads) != 1 {
		t.Error("get threads bad result")
	}
}

func TestMobile_AddThread(t *testing.T) {
	res, err := mobile.AddThread(ksuid.New().String(), "default")
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	var thrd *core.ThreadInfo
	if err := json.Unmarshal([]byte(res), &thrd); err != nil {
		t.Error(err)
		return
	}
	defaultThreadId = thrd.Id
}

func TestMobile_Threads(t *testing.T) {
	res, err := mobile.Threads()
	if err != nil {
		t.Errorf("get threads failed: %s", err)
		return
	}
	var threads []core.ThreadInfo
	if err := json.Unmarshal([]byte(res), &threads); err != nil {
		t.Error(err)
		return
	}
	if len(threads) != 2 {
		t.Error("get threads bad result")
	}
}

func TestMobile_RemoveThread(t *testing.T) {
	res, err := mobile.AddThread(ksuid.New().String(), "another")
	if err != nil {
		t.Errorf("remove thread failed: %s", err)
		return
	}
	var thrd *core.ThreadInfo
	if err := json.Unmarshal([]byte(res), &thrd); err != nil {
		t.Error(err)
		return
	}
	res2, err := mobile.RemoveThread(thrd.Id)
	if err != nil {
		t.Error(err)
		return
	}
	if err != nil {
		t.Errorf("remove thread failed: %s", err)
	}
	if res2 == "" {
		t.Errorf("remove thread bad result: %s", err)
	}
}

func TestMobile_PrepareFiles(t *testing.T) {
	res, err := mobile.PrepareFiles("../mill/testdata/image.jpeg", defaultThreadId)
	if err != nil {
		t.Errorf("prepare files failed: %s", err)
		return
	}
	dir := core.Directory{}
	if err := json.Unmarshal([]byte(res), &dir); err != nil {
		t.Error(err)
		return
	}
	if len(dir) != 6 {
		t.Error("wrong number of files")
	}
	prepared = res
}

func TestMobile_AddThreadFiles(t *testing.T) {
	res, err := mobile.AddThreadFiles(prepared, defaultThreadId, "hello")
	if err != nil {
		t.Errorf("add thread files failed: %s", err)
		return
	}
	info := core.BlockInfo{}
	if err := json.Unmarshal([]byte(res), &info); err != nil {
		t.Error(err)
	}
	filesBlock = info
	time.Sleep(time.Second)
}

func TestMobile_AddThreadFilesByTarget(t *testing.T) {
	res, err := mobile.AddThreadFilesByTarget(filesBlock.Target, defaultThreadId, "hello again")
	if err != nil {
		t.Errorf("add thread files by target failed: %s", err)
		return
	}
	info := &core.BlockInfo{}
	if err := json.Unmarshal([]byte(res), &info); err != nil {
		t.Error(err)
	}
}

func TestMobile_AddThreadComment(t *testing.T) {
	if _, err := mobile.AddThreadComment(filesBlock.Id, "hell yeah"); err != nil {
		t.Errorf("add thread comment failed: %s", err)
	}
}

func TestMobile_AddThreadLike(t *testing.T) {
	if _, err := mobile.AddThreadLike(filesBlock.Id); err != nil {
		t.Errorf("add thread like failed: %s", err)
	}
}

func TestMobile_ThreadFiles(t *testing.T) {
	res, err := mobile.ThreadFiles("", -1, defaultThreadId)
	if err != nil {
		t.Errorf("get thread files failed: %s", err)
		return
	}
	if err := json.Unmarshal([]byte(res), &files); err != nil {
		t.Error(err)
		return
	}
	if len(files) != 2 {
		t.Errorf("get thread files bad result")
	}
	if len(files[1].Comments) != 1 {
		t.Errorf("file comments bad result")
	}
	if len(files[1].Likes) != 1 {
		t.Errorf("file likes bad result")
	}
}

func TestMobile_ThreadFilesBadThread(t *testing.T) {
	if _, err := mobile.ThreadFiles("", -1, "empty"); err == nil {
		t.Error("get thread files from bad thread should fail")
	}
}

func TestMobile_FileData(t *testing.T) {
	res, err := mobile.FileData(files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Errorf("get file data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get file data bad result")
	}
}

func TestMobile_AddThreadIgnore(t *testing.T) {
	if _, err := mobile.AddThreadIgnore(filesBlock.Id); err != nil {
		t.Errorf("add thread ignore failed: %s", err)
	}
	res, err := mobile.ThreadFiles("", -1, defaultThreadId)
	if err != nil {
		t.Errorf("get thread files failed: %s", err)
		return
	}
	var files []core.ThreadFilesInfo
	if err := json.Unmarshal([]byte(res), &files); err != nil {
		t.Error(err)
		return
	}
	if len(files) != 1 {
		t.Errorf("thread ignore bad result")
	}
}

func TestMobile_PhotoDataForMinWidth(t *testing.T) {
	large, err := mobile.FileData(files[0].Files[0].Links["large"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	medium, err := mobile.FileData(files[0].Files[0].Links["medium"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	small, err := mobile.FileData(files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	thumb, err := mobile.FileData(files[0].Files[0].Links["thumb"].Hash)
	if err != nil {
		t.Error(err)
		return
	}

	pth := files[0].Target + "/0"

	d1, err := mobile.ImageFileDataForMinWidth(pth, 2000)
	if err != nil {
		t.Error(err)
		return
	}
	if d1 != large {
		t.Errorf("expected large result")
		return
	}

	d2, err := mobile.ImageFileDataForMinWidth(pth, 600)
	if err != nil {
		t.Error(err)
		return
	}
	if d2 != medium {
		t.Errorf("expected medium result")
		return
	}

	d3, err := mobile.ImageFileDataForMinWidth(pth, 320)
	if err != nil {
		t.Error(err)
		return
	}
	if d3 != small {
		t.Errorf("expected small result")
		return
	}

	d4, err := mobile.ImageFileDataForMinWidth(pth, 80)
	if err != nil {
		t.Error(err)
		return
	}
	if d4 != thumb {
		t.Errorf("expected thumb result")
		return
	}
}

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
	if err := mobile.SetAvatar(files[0].Files[0].Links["large"].Hash); err != nil {
		t.Errorf("set avatar failed: %s", err)
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

func TestMobile_Teardown(t *testing.T) {
	mobile = nil
	os.RemoveAll(repoPath)
}
