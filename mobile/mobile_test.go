package mobile_test

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/proto"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/core"
	. "github.com/textileio/textile-go/mobile"
	"github.com/textileio/textile-go/pb"
)

type TestMessenger struct{}

func (tm *TestMessenger) Notify(event *Event) {}

type TestCallback struct{}

func (tc *TestCallback) Call(payload []byte, err error) {
	if err != nil {
		fmt.Println(fmt.Errorf("callback error: %s", err))
		return
	}
	pre := new(pb.MobilePreparedFiles)
	if err := proto.Unmarshal(payload, pre); err != nil {
		fmt.Println(fmt.Errorf("callback unmarshal error: %s", err))
	}
}

var repoPath1 = "testdata/.textile1"
var repoPath2 = "testdata/.textile2"

var recovery string
var seed string

var mobile1 *Mobile
var mobile2 *Mobile

var thrdId string
var dir []byte
var filesBlock core.BlockInfo
var files []core.ThreadFilesInfo
var invite ExternalInvite

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
	os.RemoveAll(repoPath1)
	if err := InitRepo(&InitConfig{
		Seed:     seed,
		RepoPath: repoPath1,
	}); err != nil {
		t.Errorf("init mobile repo failed: %s", err)
	}
}

func TestMigrateRepo(t *testing.T) {
	if err := MigrateRepo(&MigrateConfig{
		RepoPath: repoPath1,
	}); err != nil {
		t.Errorf("migrate mobile repo failed: %s", err)
	}
}

func TestNewTextile(t *testing.T) {
	config := &RunConfig{
		RepoPath: repoPath1,
		LogLevels: `{
			"tex-core":   "debug",
			"tex-mobile": "debug"
		}`,
	}
	var err error
	mobile1, err = NewTextile(config, &TestMessenger{})
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestNewTextileAgain(t *testing.T) {
	logLevels, err := json.Marshal(map[string]string{
		"tex-core":   "debug",
		"tex-mobile": "debug",
	})
	if err != nil {
		t.Errorf("unable to marshal test map")
	}
	config := &RunConfig{
		RepoPath:  repoPath1,
		LogLevels: string(logLevels),
	}
	if _, err := NewTextile(config, &TestMessenger{}); err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestMobile_Start(t *testing.T) {
	if err := mobile1.Start(); err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestMobile_StartAgain(t *testing.T) {
	if err := mobile1.Start(); err != nil {
		t.Errorf("attempt to start a running node failed: %s", err)
	}
}

func TestMobile_Address(t *testing.T) {
	if mobile1.Address() == "" {
		t.Error("got bad address")
	}
}

func TestMobile_Seed(t *testing.T) {
	if mobile1.Seed() == "" {
		t.Error("got bad seed")
	}
}

func TestMobile_AddThread(t *testing.T) {
	res, err := mobile1.AddThread(ksuid.New().String(), "test")
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	var thrd *core.ThreadInfo
	if err := json.Unmarshal([]byte(res), &thrd); err != nil {
		t.Error(err)
		return
	}
	thrdId = thrd.Id
}

func TestMobile_AddPeerToThread(t *testing.T) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		t.Fatal(err)
	}

	if err := mobile1.AddPeerToThread(id.Pretty(), thrdId); err != nil {
		t.Errorf("add peer to thread failed: %s", err)
		return
	}
}

func TestMobile_Threads(t *testing.T) {
	res, err := mobile1.Threads()
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

func TestMobile_RemoveThread(t *testing.T) {
	res, err := mobile1.AddThread(ksuid.New().String(), "another")
	if err != nil {
		t.Errorf("remove thread failed: %s", err)
		return
	}
	var thrd *core.ThreadInfo
	if err := json.Unmarshal([]byte(res), &thrd); err != nil {
		t.Error(err)
		return
	}
	res2, err := mobile1.RemoveThread(thrd.Id)
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
	res, err := mobile1.PrepareFiles("../mill/testdata/image.jpeg", thrdId)
	if err != nil {
		t.Errorf("prepare files failed: %s", err)
		return
	}
	pre := new(pb.MobilePreparedFiles)
	if err := proto.Unmarshal(res, pre); err != nil {
		t.Error(err)
		return
	}
	if len(pre.Dir.Files) != 6 {
		t.Error("wrong number of files")
	}
	dir, err = proto.Marshal(pre.Dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_PrepareFilesAsync(t *testing.T) {
	mobile1.PrepareFilesAsync("../mill/testdata/image.jpeg", thrdId, &TestCallback{})
}

func TestMobile_AddThreadFiles(t *testing.T) {
	res, err := mobile1.AddThreadFiles(dir, thrdId, "hello")
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
	res, err := mobile1.AddThreadFilesByTarget(filesBlock.Target, thrdId, "hello again")
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
	if _, err := mobile1.AddThreadComment(filesBlock.Id, "hell yeah"); err != nil {
		t.Errorf("add thread comment failed: %s", err)
	}
}

func TestMobile_AddThreadLike(t *testing.T) {
	if _, err := mobile1.AddThreadLike(filesBlock.Id); err != nil {
		t.Errorf("add thread like failed: %s", err)
	}
}

func TestMobile_ThreadFiles(t *testing.T) {
	res, err := mobile1.ThreadFiles("", -1, thrdId)
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
	if _, err := mobile1.ThreadFiles("", -1, "empty"); err == nil {
		t.Error("get thread files from bad thread should fail")
	}
}

func TestMobile_FileData(t *testing.T) {
	res, err := mobile1.FileData(files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Errorf("get file data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get file data bad result")
	}
}

func TestMobile_AddThreadIgnore(t *testing.T) {
	if _, err := mobile1.AddThreadIgnore(filesBlock.Id); err != nil {
		t.Errorf("add thread ignore failed: %s", err)
	}
	res, err := mobile1.ThreadFiles("", -1, thrdId)
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
	large, err := mobile1.FileData(files[0].Files[0].Links["large"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	medium, err := mobile1.FileData(files[0].Files[0].Links["medium"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	small, err := mobile1.FileData(files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	thumb, err := mobile1.FileData(files[0].Files[0].Links["thumb"].Hash)
	if err != nil {
		t.Error(err)
		return
	}

	pth := files[0].Target + "/0"

	d1, err := mobile1.ImageFileDataForMinWidth(pth, 2000)
	if err != nil {
		t.Error(err)
		return
	}
	if d1 != large {
		t.Errorf("expected large result")
		return
	}

	d2, err := mobile1.ImageFileDataForMinWidth(pth, 600)
	if err != nil {
		t.Error(err)
		return
	}
	if d2 != medium {
		t.Errorf("expected medium result")
		return
	}

	d3, err := mobile1.ImageFileDataForMinWidth(pth, 320)
	if err != nil {
		t.Error(err)
		return
	}
	if d3 != small {
		t.Errorf("expected small result")
		return
	}

	d4, err := mobile1.ImageFileDataForMinWidth(pth, 80)
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
	res, err := mobile1.Overview()
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
	<-mobile1.OnlineCh()
	if err := mobile1.SetUsername("boomer"); err != nil {
		t.Errorf("set username failed: %s", err)
		return
	}
}

func TestMobile_SetAvatar(t *testing.T) {
	if err := mobile1.SetAvatar(files[0].Files[0].Links["large"].Hash); err != nil {
		t.Errorf("set avatar failed: %s", err)
		return
	}
}

func TestMobile_Profile(t *testing.T) {
	profs, err := mobile1.Profile()
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

func TestMobile_AddContact(t *testing.T) {
	if err := mobile1.AddContact("Qm123", "Pabc", "joe"); err != nil {
		t.Errorf("add contact failed: %s", err)
		return
	}
}

func TestMobile_AddThreadInvite(t *testing.T) {
	var err error
	mobile2, err = createAndStartMobile(repoPath2, true)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := mobile2.AddThread(ksuid.New().String(), "test2")
	if err != nil {
		t.Error(err)
		return
	}
	var thrd *core.ThreadInfo
	if err := json.Unmarshal([]byte(res), &thrd); err != nil {
		t.Error(err)
		return
	}

	pid, err := mobile1.PeerId()
	if err != nil {
		t.Error(err)
		return
	}

	hash, err := mobile2.AddThreadInvite(thrd.Id, pid)
	if err != nil {
		t.Error(err)
		return
	}

	if hash == "" {
		t.Errorf("bad invite result: %s", hash)
	}
}

func TestMobile_AddExternalThreadInvite(t *testing.T) {
	res, err := mobile1.AddExternalThreadInvite(thrdId)
	if err != nil {
		t.Error(err)
		return
	}
	if err := json.Unmarshal([]byte(res), &invite); err != nil {
		t.Error(err)
		return
	}
	if invite.Key == "" {
		t.Errorf("bad invite result: %s", res)
	}
}

func TestMobile_AcceptExternalThreadInvite(t *testing.T) {
	hash, err := mobile2.AcceptExternalThreadInvite(invite.Id, invite.Key)
	if err != nil {
		t.Error(err)
		return
	}

	if hash == "" {
		t.Errorf("bad accept external invite result: %s", hash)
	}
}

func TestMobile_Notifications(t *testing.T) {
	res, err := mobile1.Notifications("", -1)
	if err != nil {
		t.Error(err)
		return
	}
	var notes []core.NotificationInfo
	if err := json.Unmarshal([]byte(res), &notes); err != nil {
		t.Error(err)
		return
	}
	if len(notes) != 1 {
		t.Error("get notifications bad result")
		return
	}
}

func TestMobile_CountUnreadNotifications(t *testing.T) {
	if mobile1.CountUnreadNotifications() != 1 {
		t.Error("count unread notifications bad result")
	}
}

func TestMobile_ReadAllNotifications(t *testing.T) {
	if err := mobile1.ReadAllNotifications(); err != nil {
		t.Error(err)
	}
	if mobile1.CountUnreadNotifications() != 0 {
		t.Error("read all notifications bad result")
	}
}

func TestMobile_Stop(t *testing.T) {
	if err := mobile1.Stop(); err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}

func TestMobile_StopAgain(t *testing.T) {
	if err := mobile1.Stop(); err != nil {
		t.Errorf("stop mobile node again should not return error: %s", err)
	}
}

func TestMobile_Teardown(t *testing.T) {
	mobile1 = nil
	mobile2.Stop()
	mobile2 = nil
	os.RemoveAll(repoPath1)
	os.RemoveAll(repoPath2)
}

func createAndStartMobile(repoPath string, waitForOnline bool) (*Mobile, error) {
	os.RemoveAll(repoPath)

	recovery, err := NewWallet(12)
	if err != nil {
		return nil, err
	}

	res, err := WalletAccountAt(recovery, 0, "")
	if err != nil {
		return nil, err
	}
	accnt := WalletAccount{}
	if err := json.Unmarshal([]byte(res), &accnt); err != nil {
		return nil, err
	}

	if err := InitRepo(&InitConfig{
		Seed:     accnt.Seed,
		RepoPath: repoPath,
	}); err != nil {
		return nil, err
	}

	mobile, err := NewTextile(&RunConfig{RepoPath: repoPath}, &TestMessenger{})
	if err != nil {
		return nil, err
	}

	if err := mobile.Start(); err != nil {
		return nil, err
	}

	if waitForOnline {
		<-mobile.OnlineCh()
	}

	return mobile, nil
}
