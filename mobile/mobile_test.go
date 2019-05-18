package mobile_test

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	. "github.com/textileio/go-textile/mobile"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

type TestHandler struct{}

func (th *TestHandler) Flush() {
	fmt.Println("=== MOBILE FLUSH CALLED")
}

type TestMessenger struct{}

func (tm *TestMessenger) Notify(event *Event) {
	etype := pb.MobileEventType(event.Type)
	fmt.Println(fmt.Sprintf("+++ MOBILE EVENT: %s", event.Name))

	switch etype {
	case pb.MobileEventType_NODE_START:
	case pb.MobileEventType_NODE_ONLINE:
	case pb.MobileEventType_NODE_STOP:
	case pb.MobileEventType_WALLET_UPDATE:
	case pb.MobileEventType_THREAD_UPDATE:
	case pb.MobileEventType_NOTIFICATION:
	case pb.MobileEventType_QUERY_RESPONSE:
		res := new(pb.MobileQueryEvent)
		if err := proto.Unmarshal(event.Data, res); err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(fmt.Sprintf("+++ MOBILE QUERY EVENT: %s", res.Type.String()))

		switch res.Type {
		case pb.MobileQueryEvent_DATA:
			switch res.Data.Value.TypeUrl {
			case "/CafeClientThread":
				val := new(pb.CafeClientThread)
				if err := ptypes.UnmarshalAny(res.Data.Value, val); err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(fmt.Sprintf("+++ FOUND CLIENT THREAD (qid=%s): %s", res.Id, val.Id))

			case "/Contact":
				val := new(pb.Contact)
				if err := ptypes.UnmarshalAny(res.Data.Value, val); err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(fmt.Sprintf("+++ FOUND CONTACT (qid=%s): %s", res.Id, val.Address))
			}
		case pb.MobileQueryEvent_DONE:
			fmt.Println(fmt.Sprintf("+++ DONE (qid=%s)", res.Id))
		case pb.MobileQueryEvent_ERROR:
			fmt.Println(fmt.Sprintf("+++ ERROR (%d) (qid=%s): %s", res.Error.Code, res.Id, res.Error.Message))
		}
	}
}

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
var filesBlock *pb.Block
var files []*pb.Files
var invite *pb.ExternalInvite

var contact = &pb.Contact{
	Address: "address1",
	Name:    "joe",
	Avatar:  "Qm123",
	Peers: []*pb.Peer{
		{
			Id:      "abcde",
			Address: "address1",
			Name:    "joe",
			Avatar:  "Qm123",
			Inboxes: []*pb.Cafe{{
				Peer:     "peer",
				Address:  "address",
				Api:      "v0",
				Protocol: "/textile/cafe/1.0.0",
				Node:     "v1.0.0",
				Url:      "https://mycafe.com",
			}},
		},
	},
}

var schema = `
{
  "pin": true,
  "mill": "/json",
  "json_schema": {
    "$schema": "http://json-schema.org/draft-04/schema#",
    "$ref": "#/definitions/Log",
    "definitions": {
      "Log": {
        "required": [
          "priority",
          "timestamp",
          "hostname",
          "application",
          "pid",
          "message"
        ],
        "properties": {
          "application": {
            "type": "string"
          },
          "hostname": {
            "type": "string"
          },
          "message": {
            "type": "string"
          },
          "pid": {
            "type": "integer"
          },
          "priority": {
            "type": "integer"
          },
          "timestamp": {
            "type": "string"
          }
        },
        "additionalProperties": false,
        "type": "object"
      }
    }
  }
}
`

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
	accnt := new(pb.MobileWalletAccount)
	if err := proto.Unmarshal(res, accnt); err != nil {
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
		RepoPath:          repoPath1,
		CafeOutboxHandler: &TestHandler{},
	}
	var err error
	mobile1, err = NewTextile(config, &TestMessenger{})
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestNewTextileAgain(t *testing.T) {
	config := &RunConfig{
		RepoPath:          repoPath1,
		CafeOutboxHandler: &TestHandler{},
	}
	if _, err := NewTextile(config, &TestMessenger{}); err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestSetLogLevels(t *testing.T) {
	logLevel, err := proto.Marshal(&pb.LogLevel{
		Systems: map[string]pb.LogLevel_Level{
			"tex-core":      pb.LogLevel_DEBUG,
			"tex-datastore": pb.LogLevel_INFO,
		},
	})
	if err != nil {
		t.Errorf("unable to marshal test map")
		return
	}
	if err := mobile1.SetLogLevel(logLevel); err != nil {
		t.Errorf("attempt to set log level failed: %s", err)
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

func TestMobile_AccountThread(t *testing.T) {
	res, err := mobile1.AccountThread()
	if err != nil {
		t.Errorf("error getting account thread: %s", err)
	}
	thrd := new(pb.Thread)
	if err := proto.Unmarshal(res, thrd); err != nil {
		t.Error(err)
	}
	if thrd.Id == "" {
		t.Error("missing account thread")
	}
}

func TestMobile_AddThread(t *testing.T) {
	conf := &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_MEDIA,
		},
		Type:    pb.Thread_OPEN,
		Sharing: pb.Thread_SHARED,
	}
	mconf, err := proto.Marshal(conf)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := mobile1.AddThread(mconf)
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	thrd := new(pb.Thread)
	if err := proto.Unmarshal(res, thrd); err != nil {
		t.Error(err)
	}
	thrdId = thrd.Id
}

func TestMobile_AddThreadWithSchemaJson(t *testing.T) {
	conf := &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test",
		Schema: &pb.AddThreadConfig_Schema{
			Json: schema,
		},
		Type:    pb.Thread_READ_ONLY,
		Sharing: pb.Thread_INVITE_ONLY,
	}
	mconf, err := proto.Marshal(conf)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := mobile1.AddThread(mconf)
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	thrd := new(pb.Thread)
	if err := proto.Unmarshal(res, thrd); err != nil {
		t.Error(err)
		return
	}
	res2, err := mobile1.RemoveThread(thrd.Id)
	if err != nil {
		t.Error(err)
		return
	}
	if res2 == "" {
		t.Errorf("remove thread bad result: %s", err)
	}
}

func TestMobile_Threads(t *testing.T) {
	res, err := mobile1.Threads()
	if err != nil {
		t.Errorf("get threads failed: %s", err)
		return
	}
	list := new(pb.ThreadList)
	if err := proto.Unmarshal(res, list); err != nil {
		t.Error(err)
		return
	}
	if len(list.Items) != 1 {
		t.Error("get threads bad result")
	}
}

func TestMobile_RemoveThread(t *testing.T) {
	conf := &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "another",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_CAMERA_ROLL,
		},
		Type:    pb.Thread_PRIVATE,
		Sharing: pb.Thread_NOT_SHARED,
	}
	mconf, err := proto.Marshal(conf)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := mobile1.AddThread(mconf)
	if err != nil {
		t.Errorf("remove thread failed: %s", err)
		return
	}
	thrd := new(pb.Thread)
	if err := proto.Unmarshal(res, thrd); err != nil {
		t.Error(err)
		return
	}
	res2, err := mobile1.RemoveThread(thrd.Id)
	if err != nil {
		t.Error(err)
		return
	}
	if res2 == "" {
		t.Errorf("remove thread bad result: %s", err)
	}
}

func TestMobile_AddMessage(t *testing.T) {
	if _, err := mobile1.AddMessage(thrdId, "ping pong"); err != nil {
		t.Errorf("add thread message failed: %s", err)
	}
}

func TestMobile_Messages(t *testing.T) {
	res, err := mobile1.Messages("", -1, thrdId)
	if err != nil {
		t.Errorf("thread messages failed: %s", err)
		return
	}
	list := new(pb.TextList)
	if err := proto.Unmarshal(res, list); err != nil {
		t.Error(err)
		return
	}
	if len(list.Items) != 1 {
		t.Error("wrong number of messages")
	}
}

func TestMobile_PrepareFilesSync(t *testing.T) {
	input := "howdy"

	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	conf := &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "what",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_BLOB,
		},
		Type:    pb.Thread_OPEN,
		Sharing: pb.Thread_SHARED,
	}
	mconf, err := proto.Marshal(conf)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := mobile1.AddThread(mconf)
	if err != nil {
		t.Errorf("remove thread failed: %s", err)
		return
	}
	thrd := new(pb.Thread)
	if err := proto.Unmarshal(res, thrd); err != nil {
		t.Error(err)
		return
	}

	res2, err := mobile1.PrepareFilesSync(encoded, thrd.Id)
	if err != nil {
		t.Errorf("prepare files failed: %s", err)
		return
	}
	pre := new(pb.MobilePreparedFiles)
	if err := proto.Unmarshal(res2, pre); err != nil {
		t.Error(err)
		return
	}

	dir, err := proto.Marshal(pre.Dir)
	if err != nil {
		t.Error(err)
		return
	}

	res3, err := mobile1.AddFiles(dir, thrd.Id, "")
	if err != nil {
		t.Errorf("add thread files failed: %s", err)
		return
	}
	block := new(pb.Block)
	if err := proto.Unmarshal(res3, block); err != nil {
		t.Error(err)
		return
	}

	res4, err := mobile1.Files(thrd.Id, "", -1)
	if err != nil {
		t.Errorf("get thread files failed: %s", err)
		return
	}
	list := new(pb.FilesList)
	if err := proto.Unmarshal(res4, list); err != nil {
		t.Error(err)
		return
	}

	res5, err := mobile1.FileContent(list.Items[0].Files[0].File.Hash)
	if err != nil {
		t.Error(err)
		return
	}

	res6 := util.SplitString(res5, ",")
	res7, err := base64.StdEncoding.DecodeString(res6[1])
	if err != nil {
		t.Error(err)
		return
	}
	output := string(res7)

	if output != input {
		t.Error("file output does not match input")
	}
}

func TestMobile_PrepareFiles(t *testing.T) {
	mobile1.PrepareFiles("hello", thrdId, &TestCallback{})
}

func TestMobile_PrepareFilesByPathSync(t *testing.T) {
	res, err := mobile1.PrepareFilesByPathSync("../mill/testdata/image.jpeg", thrdId)
	if err != nil {
		t.Errorf("prepare files failed: %s", err)
		return
	}
	pre := new(pb.MobilePreparedFiles)
	if err := proto.Unmarshal(res, pre); err != nil {
		t.Error(err)
		return
	}
	if len(pre.Dir.Files) != 3 {
		t.Error("wrong number of files")
	}
	dir, err = proto.Marshal(pre.Dir)
	if err != nil {
		t.Error(err)
		return
	}

	res2, err := mobile1.PrepareFilesByPathSync(pre.Dir.Files["large"].Hash, thrdId)
	if err != nil {
		t.Errorf("prepare files by existing hash failed: %s", err)
		return
	}
	pre2 := new(pb.MobilePreparedFiles)
	if err := proto.Unmarshal(res2, pre2); err != nil {
		t.Error(err)
		return
	}
	if len(pre2.Dir.Files) != 3 {
		t.Error("wrong number of files")
	}
}

func TestMobile_PrepareFilesByPath(t *testing.T) {
	mobile1.PrepareFilesByPath("../mill/testdata/image.png", thrdId, &TestCallback{})
}

func TestMobile_AddFiles(t *testing.T) {
	res, err := mobile1.AddFiles(dir, thrdId, "hello")
	if err != nil {
		t.Errorf("add thread files failed: %s", err)
		return
	}
	block := new(pb.Block)
	if err := proto.Unmarshal(res, block); err != nil {
		t.Error(err)
		return
	}
	filesBlock = block
	time.Sleep(time.Second)
}

func TestMobile_AddFilesByTarget(t *testing.T) {
	res, err := mobile1.AddFilesByTarget(filesBlock.Target, thrdId, "hello again")
	if err != nil {
		t.Errorf("add thread files by target failed: %s", err)
		return
	}
	block := new(pb.Block)
	if err := proto.Unmarshal(res, block); err != nil {
		t.Error(err)
	}
}

func TestMobile_AddComment(t *testing.T) {
	if _, err := mobile1.AddComment(filesBlock.Id, "hell yeah"); err != nil {
		t.Errorf("add thread comment failed: %s", err)
	}
}

func TestMobile_AddLike(t *testing.T) {
	if _, err := mobile1.AddLike(filesBlock.Id); err != nil {
		t.Errorf("add thread like failed: %s", err)
	}
}

func TestMobile_Files(t *testing.T) {
	res, err := mobile1.Files(thrdId, "", -1)
	if err != nil {
		t.Errorf("get thread files failed: %s", err)
		return
	}
	list := new(pb.FilesList)
	if err := proto.Unmarshal(res, list); err != nil {
		t.Error(err)
		return
	}
	files = list.Items
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

func TestMobile_FilesBadThread(t *testing.T) {
	if _, err := mobile1.Files("empty", "", -1); err == nil {
		t.Error("get thread files from bad thread should fail")
	}
}

func TestMobile_FileContent(t *testing.T) {
	res, err := mobile1.FileContent(files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Errorf("get file data failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get file data bad result")
	}
}

func TestMobile_AddIgnore(t *testing.T) {
	if _, err := mobile1.AddIgnore(filesBlock.Id); err != nil {
		t.Errorf("add thread ignore failed: %s", err)
		return
	}
	res, err := mobile1.Files(thrdId, "", -1)
	if err != nil {
		t.Errorf("get thread files failed: %s", err)
		return
	}
	list := new(pb.FilesList)
	if err := proto.Unmarshal(res, list); err != nil {
		t.Error(err)
		return
	}
	if len(list.Items) != 1 {
		t.Errorf("thread ignore bad result")
	}
}

func TestMobile_Feed(t *testing.T) {
	req, err := proto.Marshal(&pb.FeedRequest{
		Thread: thrdId,
		Limit:  20,
		Mode:   pb.FeedRequest_STACKS,
	})
	if err != nil {
		t.Error(err)
		return
	}

	res, err := mobile1.Feed(req)
	if err != nil {
		t.Errorf("get thread feed failed: %s", err)
		return
	}
	list := new(pb.FeedItemList)
	if err := proto.Unmarshal(res, list); err != nil {
		t.Error(err)
		return
	}
	if list.Count != 3 {
		t.Errorf("get thread feed bad result")
	}
}

func TestMobile_ImageFileContentForMinWidth(t *testing.T) {
	large, err := mobile1.FileContent(files[0].Files[0].Links["large"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	small, err := mobile1.FileContent(files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Error(err)
		return
	}
	thumb, err := mobile1.FileContent(files[0].Files[0].Links["thumb"].Hash)
	if err != nil {
		t.Error(err)
		return
	}

	pth := files[0].Target + "/0"

	d1, err := mobile1.ImageFileContentForMinWidth(pth, 2000)
	if err != nil {
		t.Error(err)
		return
	}
	if d1 != large {
		t.Errorf("expected large result")
		return
	}

	d2, err := mobile1.ImageFileContentForMinWidth(pth, 600)
	if err != nil {
		t.Error(err)
		return
	}
	if d2 != large {
		t.Errorf("expected large result")
		return
	}

	d3, err := mobile1.ImageFileContentForMinWidth(pth, 320)
	if err != nil {
		t.Error(err)
		return
	}
	if d3 != small {
		t.Errorf("expected small result")
		return
	}

	d4, err := mobile1.ImageFileContentForMinWidth(pth, 80)
	if err != nil {
		t.Error(err)
		return
	}
	if d4 != thumb {
		t.Errorf("expected thumb result")
	}
}

func TestMobile_Summary(t *testing.T) {
	res, err := mobile1.Summary()
	if err != nil {
		t.Errorf("get summary failed: %s", err)
		return
	}
	summary := new(pb.Summary)
	if err := proto.Unmarshal(res, summary); err != nil {
		t.Error(err)
	}
}

func TestMobile_SetUsername(t *testing.T) {
	<-mobile1.OnlineCh()
	if err := mobile1.SetName("boomer"); err != nil {
		t.Errorf("set username failed: %s", err)
	}
}

func TestMobile_Profile(t *testing.T) {
	profs, err := mobile1.Profile()
	if err != nil {
		t.Errorf("get profile failed: %s", err)
		return
	}
	prof := new(pb.Peer)
	if err := proto.Unmarshal(profs, prof); err != nil {
		t.Error(err)
	}
}

func TestMobile_AddContact(t *testing.T) {
	payload, err := proto.Marshal(contact)
	if err != nil {
		t.Error(err)
		return
	}
	if err := mobile1.AddContact(payload); err != nil {
		t.Errorf("add contact failed: %s", err)
	}
}

func TestMobile_AddContactAgain(t *testing.T) {
	payload, err := proto.Marshal(contact)
	if err != nil {
		t.Error(err)
		return
	}
	if err := mobile1.AddContact(payload); err != nil {
		t.Errorf("adding duplicate contact should not throw error")
	}
}

func TestMobile_Contact(t *testing.T) {
	self, err := mobile1.Contact(mobile1.Address())
	if err != nil {
		t.Errorf("get own contact failed: %s", err)
		return
	}
	contact := new(pb.Contact)
	if err := proto.Unmarshal(self, contact); err != nil {
		t.Error(err)
	}
}

func TestMobile_AddInvite(t *testing.T) {
	var err error
	mobile2, err = createAndStartMobile(repoPath2, true)
	if err != nil {
		t.Error(err)
		return
	}

	conf := &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test2",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_MEDIA,
		},
		Type:    pb.Thread_OPEN,
		Sharing: pb.Thread_SHARED,
	}
	mconf, err := proto.Marshal(conf)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := mobile2.AddThread(mconf)
	if err != nil {
		t.Error(err)
		return
	}
	thrd := new(pb.Thread)
	if err := proto.Unmarshal([]byte(res), thrd); err != nil {
		t.Error(err)
		return
	}

	contact1, err := mobile1.Contact(mobile1.Address())
	if err != nil {
		t.Error(err)
		return
	}

	if err := mobile2.AddContact(contact1); err != nil {
		t.Error(err)
		return
	}

	if err := mobile2.AddInvite(thrd.Id, mobile1.Address()); err != nil {
		t.Error(err)
		return
	}
}

func TestMobile_AddExternalInvite(t *testing.T) {
	res, err := mobile1.AddExternalInvite(thrdId)
	if err != nil {
		t.Error(err)
		return
	}
	invite = new(pb.ExternalInvite)
	if err := proto.Unmarshal(res, invite); err != nil {
		t.Error(err)
		return
	}
	if invite.Key == "" {
		t.Errorf("bad invite result: %s", res)
	}
}

func TestMobile_AcceptExternalInvite(t *testing.T) {
	_, err := mobile1.AcceptExternalInvite(invite.Id, invite.Key)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMobile_Notifications(t *testing.T) {
	res, err := mobile1.Notifications("", -1)
	if err != nil {
		t.Error(err)
		return
	}
	notes := new(pb.NotificationList)
	if err := proto.Unmarshal(res, notes); err != nil {
		t.Error(err)
	}
}

func TestMobile_CountUnreadNotifications(t *testing.T) {
	mobile1.CountUnreadNotifications()
}

func TestMobile_ReadAllNotifications(t *testing.T) {
	if err := mobile1.ReadAllNotifications(); err != nil {
		t.Error(err)
		return
	}
	if mobile1.CountUnreadNotifications() != 0 {
		t.Error("read all notifications bad result")
	}
}

func TestMobile_SearchContacts(t *testing.T) {
	query, err := proto.Marshal(&pb.ContactQuery{Address: mobile2.Address()})
	if err != nil {
		t.Fatal(err)
	}
	opts, err := proto.Marshal(&pb.QueryOptions{
		Wait:  10,
		Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	handle, err := mobile1.SearchContacts(query, opts)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(fmt.Sprintf("query ID: %s", handle.Id))

	timer := time.NewTimer(3 * time.Second)
	<-timer.C

	handle.Cancel()
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
	accnt := new(pb.MobileWalletAccount)
	if err := proto.Unmarshal(res, accnt); err != nil {
		return nil, err
	}

	if err := InitRepo(&InitConfig{
		Seed:     accnt.Seed,
		RepoPath: repoPath,
	}); err != nil {
		return nil, err
	}

	mobile, err := NewTextile(&RunConfig{
		RepoPath:          repoPath,
		CafeOutboxHandler: &TestHandler{},
	}, &TestMessenger{})
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
