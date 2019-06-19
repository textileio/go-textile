package mobile

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

var testVars = struct {
	repoPath1 string
	repoPath2 string

	recovery string
	seed     string

	mobile1 *Mobile
	mobile2 *Mobile

	thrdId     string
	dir        []byte
	filesBlock *pb.Block
	files      []*pb.Files
	invite     *pb.ExternalInvite
	avatar     string
}{
	repoPath1: "testdata/.textile1",
	repoPath2: "testdata/.textile2",
}

type testHandler struct{}

func (th *testHandler) Flush() {
	fmt.Println("=== MOBILE FLUSH CALLED")
}

type testMessenger struct{}

func (tm *testMessenger) Notify(event *Event) {
	etype := pb.MobileEventType(event.Type)
	fmt.Println(fmt.Sprintf("+++ MOBILE EVENT: %s", event.Name))

	switch etype {
	case pb.MobileEventType_NODE_START:
	case pb.MobileEventType_NODE_ONLINE:
	case pb.MobileEventType_NODE_STOP:
	case pb.MobileEventType_ACCOUNT_UPDATE:
	case pb.MobileEventType_THREAD_UPDATE:
	case pb.MobileEventType_NOTIFICATION:
	case pb.MobileEventType_QUERY_RESPONSE:
		res := new(pb.MobileQueryEvent)
		err := proto.Unmarshal(event.Data, res)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(fmt.Sprintf("+++ MOBILE QUERY EVENT: %s", res.Type.String()))

		switch res.Type {
		case pb.MobileQueryEvent_DATA:
			switch res.Data.Value.TypeUrl {
			case "/CafeClientThread":
				val := new(pb.CafeClientThread)
				err := ptypes.UnmarshalAny(res.Data.Value, val)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(fmt.Sprintf("+++ FOUND CLIENT THREAD (qid=%s): %s", res.Id, val.Id))

			case "/Contact":
				val := new(pb.Contact)
				err := ptypes.UnmarshalAny(res.Data.Value, val)
				if err != nil {
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
	case pb.MobileEventType_CAFE_SYNC_GROUP_UPDATE,
		pb.MobileEventType_CAFE_SYNC_GROUP_COMPLETE,
		pb.MobileEventType_CAFE_SYNC_GROUP_FAILED:
		status := new(pb.CafeSyncGroupStatus)
		err := proto.Unmarshal(event.Data, status)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		printSyncGroupStatus(status)
	}
}

func TestNewWallet(t *testing.T) {
	mnemonic, err := NewWallet(12)
	if err != nil {
		t.Fatalf("new mobile wallet failed: %s", err)
	}
	testVars.recovery = mnemonic
}

func TestWalletAccountAt(t *testing.T) {
	res, err := WalletAccountAt(testVars.recovery, 0, "")
	if err != nil {
		t.Fatalf("get mobile wallet account at failed: %s", err)
	}
	accnt := new(pb.MobileWalletAccount)
	err = proto.Unmarshal(res, accnt)
	if err != nil {
		t.Fatal(err)
	}
	testVars.seed = accnt.Seed
}

func TestInitRepo(t *testing.T) {
	_ = os.RemoveAll(testVars.repoPath1)
	err := InitRepo(&InitConfig{
		Seed:     testVars.seed,
		RepoPath: testVars.repoPath1,
	})
	if err != nil {
		t.Fatalf("init mobile repo failed: %s", err)
	}
}

func TestMigrateRepo(t *testing.T) {
	err := MigrateRepo(&MigrateConfig{
		RepoPath: testVars.repoPath1,
	})
	if err != nil {
		t.Fatalf("migrate mobile repo failed: %s", err)
	}
}

func TestNewTextile(t *testing.T) {
	config := &RunConfig{
		RepoPath:          testVars.repoPath1,
		CafeOutboxHandler: &testHandler{},
	}
	var err error
	testVars.mobile1, err = NewTextile(config, &testMessenger{})
	if err != nil {
		t.Fatalf("create mobile node failed: %s", err)
	}
}

func TestNewTextileAgain(t *testing.T) {
	config := &RunConfig{
		RepoPath:          testVars.repoPath1,
		CafeOutboxHandler: &testHandler{},
	}
	_, err := NewTextile(config, &testMessenger{})
	if err != nil {
		t.Fatalf("create mobile node failed: %s", err)
	}
}

func TestSetLogLevels(t *testing.T) {
	logLevel, err := proto.Marshal(&pb.LogLevel{
		Systems: map[string]pb.LogLevel_Level{
			"tex-core":      pb.LogLevel_DEBUG,
			"tex-datastore": pb.LogLevel_DEBUG,
		},
	})
	if err != nil {
		t.Fatalf("unable to marshal test map")
	}
	err = testVars.mobile1.SetLogLevel(logLevel)
	if err != nil {
		t.Fatalf("attempt to set log level failed: %s", err)
	}
}

func TestMobile_Start(t *testing.T) {
	err := testVars.mobile1.Start()
	if err != nil {
		t.Fatalf("start mobile node failed: %s", err)
	}
}

func TestMobile_StopAndStart(t *testing.T) {
	err := testVars.mobile1.Start()
	if err != nil {
		t.Fatalf("attempt to start a running node failed: %s", err)
	}
	err = testVars.mobile1.Stop()
	if err != nil {
		t.Fatalf("stop mobile node failed: %s", err)
	}
	err = testVars.mobile1.Start()
	if err != nil {
		t.Fatalf("start mobile node again failed: %s", err)
	}
}

func TestMobile_Address(t *testing.T) {
	if testVars.mobile1.Address() == "" {
		t.Fatal("got bad address")
	}
}

func TestMobile_Seed(t *testing.T) {
	if testVars.mobile1.Seed() == "" {
		t.Fatal("got bad seed")
	}
}

func TestMobile_AccountThread(t *testing.T) {
	res, err := testVars.mobile1.AccountThread()
	if err != nil {
		t.Fatalf("error getting account thread: %s", err)
	}
	thrd := new(pb.Thread)
	err = proto.Unmarshal(res, thrd)
	if err != nil {
		t.Fatal(err)
	}
	if thrd.Id == "" {
		t.Fatal("missing account thread")
	}
}

func TestMobile_AddThread(t *testing.T) {
	thrd, err := addTestThread(testVars.mobile1, &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_MEDIA,
		},
		Type:    pb.Thread_OPEN,
		Sharing: pb.Thread_SHARED,
	})
	if err != nil {
		t.Fatalf("add thread failed: %s", err)
	}
	testVars.thrdId = thrd.Id
}

func TestMobile_AddThreadWithSchemaJson(t *testing.T) {
	_, err := addTestThread(testVars.mobile1, &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test",
		Schema: &pb.AddThreadConfig_Schema{
			Json: util.TestLogSchema,
		},
		Type:    pb.Thread_READ_ONLY,
		Sharing: pb.Thread_INVITE_ONLY,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_Threads(t *testing.T) {
	res, err := testVars.mobile1.Threads()
	if err != nil {
		t.Fatalf("get threads failed: %s", err)
	}
	list := new(pb.ThreadList)
	err = proto.Unmarshal(res, list)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 2 {
		t.Fatal("get threads bad result")
	}
}

func TestMobile_RemoveThread(t *testing.T) {
	thrd, err := addTestThread(testVars.mobile1, &pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "another",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_CAMERA_ROLL,
		},
		Type:    pb.Thread_PRIVATE,
		Sharing: pb.Thread_NOT_SHARED,
	})
	if err != nil {
		t.Fatalf("add thread failed: %s", err)
	}
	res, err := testVars.mobile1.RemoveThread(thrd.Id)
	if err != nil {
		t.Fatal(err)
	}
	if res == "" {
		t.Fatal("failed to remove thread")
	}
	// try to remove it again
	res2, err := testVars.mobile1.RemoveThread(thrd.Id)
	if err == nil {
		t.Fatal(err)
	}
	if res2 != "" {
		t.Fatal("bad result removing thread again")
	}
}

func TestMobile_AddMessage(t *testing.T) {
	_, err := testVars.mobile1.AddMessage(testVars.thrdId, "ping pong")
	if err != nil {
		t.Fatalf("add thread message failed: %s", err)
	}
}

func TestMobile_Messages(t *testing.T) {
	res, err := testVars.mobile1.Messages("", -1, testVars.thrdId)
	if err != nil {
		t.Fatalf("thread messages failed: %s", err)
	}
	list := new(pb.TextList)
	err = proto.Unmarshal(res, list)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 {
		t.Fatal("wrong number of messages")
	}
}

func TestMobile_AddData(t *testing.T) {
	input := "howdy"

	conf, err := proto.Marshal(&pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "what",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_BLOB,
		},
		Type:    pb.Thread_OPEN,
		Sharing: pb.Thread_SHARED,
	})
	if err != nil {
		t.Fatal(err)
	}
	res, err := testVars.mobile1.AddThread(conf)
	if err != nil {
		t.Fatalf("add thread failed: %s", err)
	}
	thrd := new(pb.Thread)
	err = proto.Unmarshal(res, thrd)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testVars.mobile1.addData([]byte(input), thrd.Id, "caption")
	if err != nil {
		t.Fatalf("add data failed: %s", err)
	}

	res3, err := testVars.mobile1.Files(thrd.Id, "", -1)
	if err != nil {
		t.Fatalf("get thread files failed: %s", err)
	}
	list := new(pb.FilesList)
	err = proto.Unmarshal(res3, list)
	if err != nil {
		t.Fatal(err)
	}

	res4, err := testVars.mobile1.FileContent(list.Items[0].Files[0].File.Hash)
	if err != nil {
		t.Fatal(err)
	}

	res5 := util.SplitString(res4, ",")
	res6, err := base64.StdEncoding.DecodeString(res5[1])
	if err != nil {
		t.Fatal(err)
	}
	output := string(res6)

	if output != input {
		t.Fatal("file output does not match input")
	}
}

func TestMobile_AddFiles(t *testing.T) {
	hash, err := testVars.mobile1.addFiles([]string{"../mill/testdata/image.jpeg"}, testVars.thrdId, "caption")
	if err != nil {
		t.Fatalf("prepare files failed: %s", err)
	}

	block, err := testVars.mobile1.node.BlockView(hash.B58String())
	if err != nil {
		t.Fatal(err)
	}
	testVars.filesBlock = block
}

func TestMobile_ShareFiles(t *testing.T) {
	_, err := testVars.mobile1.shareFiles(testVars.filesBlock.Data, testVars.thrdId, "hello")
	if err != nil {
		t.Fatalf("share files failed: %s", err)
	}
}

func TestMobile_AddComment(t *testing.T) {
	_, err := testVars.mobile1.AddComment(testVars.filesBlock.Id, "yeah")
	if err != nil {
		t.Fatalf("add thread comment failed: %s", err)
	}
}

func TestMobile_AddLike(t *testing.T) {
	_, err := testVars.mobile1.AddLike(testVars.filesBlock.Id)
	if err != nil {
		t.Fatalf("add thread like failed: %s", err)
	}
}

func TestMobile_Files(t *testing.T) {
	res, err := testVars.mobile1.Files(testVars.thrdId, "", -1)
	if err != nil {
		t.Fatalf("get thread files failed: %s", err)
	}
	list := new(pb.FilesList)
	err = proto.Unmarshal(res, list)
	if err != nil {
		t.Fatal(err)
	}
	testVars.files = list.Items
	if len(testVars.files) != 2 {
		t.Fatalf("get thread files bad result")
	}
	if len(testVars.files[1].Comments) != 1 {
		t.Fatalf("file comments bad result")
	}
	if len(testVars.files[1].Likes) != 1 {
		t.Fatalf("file likes bad result")
	}
}

func TestMobile_FilesBadThread(t *testing.T) {
	if _, err := testVars.mobile1.Files("empty", "", -1); err == nil {
		t.Fatal("get thread files from bad thread should fail")
	}
}

func TestMobile_FileData(t *testing.T) {
	res, err := testVars.mobile1.FileContent(testVars.files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Fatalf("get file data failed: %s", err)
	}
	if len(res) == 0 {
		t.Fatalf("get file data bad result")
	}
}

func TestMobile_AddIgnore(t *testing.T) {
	_, err := testVars.mobile1.AddIgnore(testVars.filesBlock.Id)
	if err != nil {
		t.Fatalf("add thread ignore failed: %s", err)
	}
	res, err := testVars.mobile1.Files(testVars.thrdId, "", -1)
	if err != nil {
		t.Fatalf("get thread files failed: %s", err)
	}
	list := new(pb.FilesList)
	err = proto.Unmarshal(res, list)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("thread ignore bad result")
	}
}

func TestMobile_Feed(t *testing.T) {
	req, err := proto.Marshal(&pb.FeedRequest{
		Thread: testVars.thrdId,
		Limit:  20,
		Mode:   pb.FeedRequest_STACKS,
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := testVars.mobile1.Feed(req)
	if err != nil {
		t.Fatalf("get thread feed failed: %s", err)
	}
	list := new(pb.FeedItemList)
	err = proto.Unmarshal(res, list)
	if err != nil {
		t.Fatal(err)
	}
	if list.Count != 3 {
		t.Fatalf("get thread feed bad result")
	}
}

func TestMobile_ImageFileDataForMinWidth(t *testing.T) {
	large, err := testVars.mobile1.FileContent(testVars.files[0].Files[0].Links["large"].Hash)
	if err != nil {
		t.Fatal(err)
	}
	small, err := testVars.mobile1.FileContent(testVars.files[0].Files[0].Links["small"].Hash)
	if err != nil {
		t.Fatal(err)
	}
	thumb, err := testVars.mobile1.FileContent(testVars.files[0].Files[0].Links["thumb"].Hash)
	if err != nil {
		t.Fatal(err)
	}

	pth := testVars.files[0].Data + "/0"

	d1, err := testVars.mobile1.ImageFileContentForMinWidth(pth, 2000)
	if err != nil {
		t.Fatal(err)
	}
	if d1 != large {
		t.Fatalf("expected large result")
	}

	d2, err := testVars.mobile1.ImageFileContentForMinWidth(pth, 600)
	if err != nil {
		t.Fatal(err)
	}
	if d2 != large {
		t.Fatalf("expected large result")
	}

	d3, err := testVars.mobile1.ImageFileContentForMinWidth(pth, 320)
	if err != nil {
		t.Fatal(err)
	}
	if d3 != small {
		t.Fatalf("expected small result")
	}

	d4, err := testVars.mobile1.ImageFileContentForMinWidth(pth, 80)
	if err != nil {
		t.Fatal(err)
	}
	if d4 != thumb {
		t.Fatalf("expected thumb result")
	}
}

func TestMobile_Summary(t *testing.T) {
	res, err := testVars.mobile1.Summary()
	if err != nil {
		t.Fatalf("get summary failed: %s", err)
	}
	summary := new(pb.Summary)
	err = proto.Unmarshal(res, summary)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_SetUsername(t *testing.T) {
	<-testVars.mobile1.onlineCh()
	err := testVars.mobile1.SetName("boomer")
	if err != nil {
		t.Fatalf("set username failed: %s", err)
	}
}

func TestMobile_SetAvatar(t *testing.T) {
	hash1, err := testVars.mobile1.Avatar()
	if err != nil {
		t.Fatal(err)
	}

	_, err = testVars.mobile1.setAvatar("../mill/testdata/image.jpeg")
	if err != nil {
		t.Fatalf("set avatar failed: %s", err)
	}

	testVars.avatar, err = testVars.mobile1.Avatar()
	if err != nil {
		t.Fatal(err)
	}

	if testVars.avatar == hash1 {
		t.Fatal("avatar was not updated")
	}
}

func TestMobile_Profile(t *testing.T) {
	profs, err := testVars.mobile1.Profile()
	if err != nil {
		t.Fatalf("get profile failed: %s", err)
	}
	prof := new(pb.Peer)
	err = proto.Unmarshal(profs, prof)
	if err != nil {
		t.Fatal(err)
	}

	if prof.Avatar != testVars.avatar {
		t.Fatal("incorrect profile avatar")
	}
}

func TestMobile_AddContact(t *testing.T) {
	payload, err := proto.Marshal(util.TestContact)
	if err != nil {
		t.Fatal(err)
	}
	err = testVars.mobile1.AddContact(payload)
	if err != nil {
		t.Fatalf("add contact failed: %s", err)
	}
}

func TestMobile_AddContactAgain(t *testing.T) {
	payload, err := proto.Marshal(util.TestContact)
	if err != nil {
		t.Fatal(err)
	}
	err = testVars.mobile1.AddContact(payload)
	if err != nil {
		t.Fatalf("adding duplicate contact should not throw error")
	}
}

func TestMobile_Contact(t *testing.T) {
	self, err := testVars.mobile1.Contact(testVars.mobile1.Address())
	if err != nil {
		t.Fatalf("get own contact failed: %s", err)
	}
	contact := new(pb.Contact)
	err = proto.Unmarshal(self, contact)
	if err != nil {
		t.Fatal(err)
	}

	if contact.Avatar != testVars.avatar {
		t.Fatal("incorrect self contact avatar")
	}
}

func TestMobile_AddInvite(t *testing.T) {
	var err error
	testVars.mobile2, err = createAndStartPeer(InitConfig{
		RepoPath: testVars.repoPath2,
		Debug:    true,
	}, true, &testHandler{}, &testMessenger{})
	if err != nil {
		t.Fatal(err)
	}

	conf, err := proto.Marshal(&pb.AddThreadConfig{
		Key:  ksuid.New().String(),
		Name: "test2",
		Schema: &pb.AddThreadConfig_Schema{
			Preset: pb.AddThreadConfig_Schema_MEDIA,
		},
		Type:    pb.Thread_OPEN,
		Sharing: pb.Thread_SHARED,
	})
	if err != nil {
		t.Fatal(err)
	}
	res, err := testVars.mobile2.AddThread(conf)
	if err != nil {
		t.Fatal(err)
	}
	thrd := new(pb.Thread)
	err = proto.Unmarshal([]byte(res), thrd)
	if err != nil {
		t.Fatal(err)
	}

	contact1, err := testVars.mobile1.Contact(testVars.mobile1.Address())
	if err != nil {
		t.Fatal(err)
	}

	err = testVars.mobile2.AddContact(contact1)
	if err != nil {
		t.Fatal(err)
	}

	err = testVars.mobile2.AddInvite(thrd.Id, testVars.mobile1.Address())
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_AddExternalInvite(t *testing.T) {
	res, err := testVars.mobile1.AddExternalInvite(testVars.thrdId)
	if err != nil {
		t.Fatal(err)
	}
	testVars.invite = new(pb.ExternalInvite)
	err = proto.Unmarshal(res, testVars.invite)
	if err != nil {
		t.Fatal(err)
	}
	if testVars.invite.Key == "" {
		t.Fatalf("bad invite result: %s", res)
	}
}

func TestMobile_AcceptExternalInvite(t *testing.T) {
	_, err := testVars.mobile1.AcceptExternalInvite(testVars.invite.Id, testVars.invite.Key)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_Notifications(t *testing.T) {
	res, err := testVars.mobile1.Notifications("", -1)
	if err != nil {
		t.Fatal(err)
	}
	notes := new(pb.NotificationList)
	err = proto.Unmarshal(res, notes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMobile_CountUnreadNotifications(t *testing.T) {
	testVars.mobile1.CountUnreadNotifications()
}

func TestMobile_ReadAllNotifications(t *testing.T) {
	err := testVars.mobile1.ReadAllNotifications()
	if err != nil {
		t.Fatal(err)
	}
	if testVars.mobile1.CountUnreadNotifications() != 0 {
		t.Fatal("read all notifications bad result")
	}
}

func TestMobile_SearchContacts(t *testing.T) {
	query, err := proto.Marshal(&pb.ContactQuery{Address: testVars.mobile2.Address()})
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
	handle, err := testVars.mobile1.SearchContacts(query, opts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("query ID: %s", handle.Id))

	timer := time.NewTimer(3 * time.Second)
	<-timer.C

	handle.Cancel()
}

func TestMobile_Stop(t *testing.T) {
	err := testVars.mobile1.Stop()
	if err != nil {
		t.Fatalf("stop mobile node failed: %s", err)
	}
}

func TestMobile_StopAgain(t *testing.T) {
	err := testVars.mobile1.Stop()
	if err != nil {
		t.Fatalf("stop mobile node again should not return error: %s", err)
	}
}

func TestMobile_Teardown(t *testing.T) {
	testVars.mobile1 = nil
	_ = testVars.mobile2.Stop()
	testVars.mobile2 = nil
	_ = os.RemoveAll(testVars.repoPath1)
	_ = os.RemoveAll(testVars.repoPath2)
}
