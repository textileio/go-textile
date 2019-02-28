package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

var notificationStore repo.NotificationStore

func init() {
	setupNotificationDB()
}

func setupNotificationDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	notificationStore = NewNotificationStore(conn, new(sync.Mutex))
}

func TestNotificationDB_Add(t *testing.T) {
	err := notificationStore.Add(&pb.Notification{
		Id:          "abcde",
		Date:        ptypes.TimestampNow(),
		Actor:       ksuid.New().String(),
		SubjectDesc: "test",
		Subject:     ksuid.New().String(),
		Block:       ksuid.New().String(),
		Type:        pb.Notification_INVITE_RECEIVED,
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != "abcde" {
		t.Errorf(`expected "abcde" got %s`, id)
	}
}

func TestNotificationDB_Get(t *testing.T) {
	notif := notificationStore.Get("abcde")
	if notif == nil {
		t.Error("could not get notification")
	}
}

func TestNotificationDB_Read(t *testing.T) {
	err := notificationStore.Read("abcde")
	if err != nil {
		t.Error(err)
		return
	}
	notifs := notificationStore.List("", 1)
	if len(notifs.Items) == 0 || !notifs.Items[0].Read {
		t.Error("notification read bad result")
	}
}

func TestNotificationDB_ReadAll(t *testing.T) {
	setupNotificationDB()
	err := notificationStore.Add(&pb.Notification{
		Id:          "abcde",
		Date:        ptypes.TimestampNow(),
		Actor:       ksuid.New().String(),
		SubjectDesc: "test",
		Subject:     ksuid.New().String(),
		Block:       ksuid.New().String(),
		Type:        pb.Notification_INVITE_RECEIVED,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&pb.Notification{
		Id:          "abcdef",
		Date:        ptypes.TimestampNow(),
		Actor:       ksuid.New().String(),
		SubjectDesc: "test",
		Subject:     ksuid.New().String(),
		Block:       ksuid.New().String(),
		Type:        pb.Notification_PEER_JOINED,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.ReadAll()
	if err != nil {
		t.Error(err)
		return
	}
	notifs := notificationStore.List("", -1)
	if len(notifs.Items) != 2 || !notifs.Items[0].Read || !notifs.Items[1].Read {
		t.Error("notification read all bad result")
	}
}

func TestNotificationDB_List(t *testing.T) {
	setupNotificationDB()
	err := notificationStore.Add(&pb.Notification{
		Id:          "abc",
		Date:        ptypes.TimestampNow(),
		Actor:       "actor1",
		SubjectDesc: "test",
		Subject:     ksuid.New().String(),
		Block:       "block1",
		Type:        pb.Notification_INVITE_RECEIVED,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&pb.Notification{
		Id:          "def",
		Date:        util.ProtoTs(time.Now().Add(time.Minute).UnixNano()),
		Actor:       "actor1",
		SubjectDesc: "test",
		Subject:     ksuid.New().String(),
		Block:       "block2",
		Type:        pb.Notification_PEER_JOINED,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&pb.Notification{
		Id:          "ghi",
		Date:        util.ProtoTs(time.Now().Add(time.Minute * 2).UnixNano()),
		Actor:       "actor2",
		SubjectDesc: "test",
		Subject:     ksuid.New().String(),
		Block:       "block2",
		Type:        pb.Notification_COMMENT_ADDED,
	})
	if err != nil {
		t.Error(err)
	}
	err = notificationStore.Add(&pb.Notification{
		Id:          "jkl",
		Date:        util.ProtoTs(time.Now().Add(time.Minute * 3).UnixNano()),
		Actor:       "actor3",
		SubjectDesc: "test",
		Subject:     "subject1",
		Block:       "block3",
		Target:      "target",
		Type:        pb.Notification_FILES_ADDED,
	})
	if err != nil {
		t.Error(err)
	}
	all := notificationStore.List("", -1)
	if len(all.Items) != 4 {
		t.Error("returned incorrect number of notifications")
		return
	}
	limited := notificationStore.List("", 1)
	if len(limited.Items) != 1 {
		t.Error("returned incorrect number of notifications")
		return
	}
	offset := notificationStore.List(limited.Items[0].Id, -1)
	if len(offset.Items) != 3 {
		t.Error("returned incorrect number of notifications")
		return
	}
}

func TestNotificationDB_CountUnread(t *testing.T) {
	cnt := notificationStore.CountUnread()
	if cnt != 4 {
		t.Error("returned incorrect count of unread notifications")
	}
}

func TestNotificationDB_Delete(t *testing.T) {
	err := notificationStore.Delete("abc")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abc").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestNotificationDB_DeleteByActor(t *testing.T) {
	err := notificationStore.DeleteByActor("actor1")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("def").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestNotificationDB_DeleteBySubject(t *testing.T) {
	err := notificationStore.DeleteBySubject("subject1")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("jkl").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestNotificationDB_DeleteByBlock(t *testing.T) {
	err := notificationStore.DeleteByBlock("block2")
	if err != nil {
		t.Error(err)
	}
	stmt, err := notificationStore.PrepareQuery("select id from notifications where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("ghi").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
