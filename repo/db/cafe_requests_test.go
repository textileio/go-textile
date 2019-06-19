package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

var cafeRequestStore repo.CafeRequestStore

func init() {
	setupCafeRequestDB()
}

func setupCafeRequestDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	_ = initDatabaseTables(conn, "")
	cafeRequestStore = NewCafeRequestStore(conn, new(sync.Mutex))
}

func TestCafeRequestDB_Add(t *testing.T) {
	err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:               "abcde",
		Peer:             "peer",
		Target:           "zxy",
		Cafe:             testCafe,
		Group:            "group",
		SyncGroup:        "sync_group",
		Type:             pb.CafeRequest_STORE,
		Date:             ptypes.TimestampNow(),
		Size:             1024,
		Status:           pb.CafeRequest_NEW,
		Attempts:         0,
		GroupSize:        0,
		GroupTransferred: 0,
	})
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err != nil {
		t.Error(err)
		return
	}
	if id != "abcde" {
		t.Errorf(`expected "abcde" got %s`, id)
	}
}

func TestCafeRequestDB_Get(t *testing.T) {
	setupCafeRequestDB()
	err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "abcde",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE_THREAD,
		Date:      ptypes.TimestampNow(),
		Status:    pb.CafeRequest_NEW,
	})
	if err != nil {
		t.Error(err)
		return
	}
	req := cafeRequestStore.Get("abcde")
	if req == nil {
		t.Error("get request failed")
	}
}

func TestCafeRequestDB_GetGroup(t *testing.T) {
	reqs := cafeRequestStore.GetGroup("group")
	if len(reqs.Items) != 1 {
		t.Error("get request group failed")
	}
}

func TestCafeRequestDB_GetSyncGroup(t *testing.T) {
	syncGroup := cafeRequestStore.GetSyncGroup("group")
	if syncGroup != "sync_group" {
		t.Error("get request sync group failed")
	}
}

func TestCafeRequestDB_List(t *testing.T) {
	setupCafeRequestDB()
	err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "abcde",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group1",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE_THREAD,
		Date:      ptypes.TimestampNow(),
		Status:    pb.CafeRequest_NEW,
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "abcdef",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group2",
		SyncGroup: "sync_group2",
		Type:      pb.CafeRequest_STORE,
		Date:      util.ProtoTs(time.Now().Add(time.Minute).UnixNano()),
		Size:      1024,
		Status:    pb.CafeRequest_NEW,
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "abcdefg",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group3",
		SyncGroup: "sync_group2",
		Type:      pb.CafeRequest_STORE,
		Date:      util.ProtoTs(time.Now().Add(time.Minute * 2).UnixNano()),
		Size:      1024,
		Status:    pb.CafeRequest_NEW,
	})
	if err != nil {
		t.Error(err)
		return
	}
	all := cafeRequestStore.List("", -1).Items
	if len(all) != 3 {
		t.Error("returned incorrect number of requests")
		return
	}
	limited := cafeRequestStore.List("", 1).Items
	if len(limited) != 1 {
		t.Error("returned incorrect number of requests")
		return
	}
	offset := cafeRequestStore.List(limited[0].Id, -1).Items
	if len(offset) != 2 {
		t.Error("returned incorrect number of requests")
		return
	}
}

func TestCafeRequestDB_Count(t *testing.T) {
	if cafeRequestStore.Count(pb.CafeRequest_NEW) != 3 {
		t.Errorf("returned incorrect count")
	}
}

func TestCafeRequestDB_ListGroups(t *testing.T) {
	all := cafeRequestStore.ListGroups("", -1)
	if len(all) != 3 {
		t.Errorf("returned incorrect number of groups")
	}
	limited := cafeRequestStore.ListGroups("", 1)
	if len(limited) != 1 {
		t.Errorf("returned incorrect number of groups")
	}
	offset := cafeRequestStore.ListGroups(limited[0], -1)
	if len(offset) != 2 {
		t.Errorf("returned incorrect number of groups")
	}
}

func TestCafeRequestDB_SyncGroupComplete(t *testing.T) {
	if cafeRequestStore.SyncGroupComplete("sync_group2") {
		t.Error("sync group complete failed")
	}
}

func TestCafeRequestDB_UpdateStatus(t *testing.T) {
	err := cafeRequestStore.UpdateStatus("abcdef", pb.CafeRequest_PENDING)
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select status from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var status int
	err = stmt.QueryRow("abcdef").Scan(&status)
	if err != nil {
		t.Error(err)
		return
	}
	if status != 1 {
		t.Error("wrong status")
	}
}

func TestCafeRequestDB_SyncGroupStatus(t *testing.T) {
	status := cafeRequestStore.SyncGroupStatus("group3") // sync_group2
	if status.NumTotal != 2 {
		t.Errorf("wrong num total %d", status.NumTotal)
	}
	if status.NumPending != 1 {
		t.Errorf("wrong num pending %d", status.NumPending)
	}
	if status.NumComplete != 0 {
		t.Errorf("wrong num complete %d", status.NumComplete)
	}
	if status.SizeTotal != 2048 {
		t.Errorf("wrong size total %d", status.SizeTotal)
	}
	if status.SizePending != 1024 {
		t.Errorf("wrong size pending %d", status.SizePending)
	}
	if status.SizeComplete != 0 {
		t.Errorf("wrong size complete %d", status.SizeComplete)
	}
}

func TestCafeRequestDB_UpdateGroupProgress(t *testing.T) {
	err := cafeRequestStore.UpdateGroupProgress("group2", 8, 16)
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select groupTransferred from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var groupTransferred int64
	err = stmt.QueryRow("abcdef").Scan(&groupTransferred)
	if err != nil {
		t.Error(err)
		return
	}
	if groupTransferred != 8 {
		t.Error("wrong group transferred")
	}
	stmt, err = cafeRequestStore.PrepareQuery("select groupSize from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var groupSize int64
	err = stmt.QueryRow("abcdef").Scan(&groupSize)
	if err != nil {
		t.Error(err)
		return
	}
	if groupSize != 16 {
		t.Error("wrong group size")
	}
}

func TestCafeRequestDB_UpdateGroupStatus(t *testing.T) {
	err := cafeRequestStore.UpdateGroupStatus("group2", pb.CafeRequest_COMPLETE)
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select status from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var status int
	err = stmt.QueryRow("abcdef").Scan(&status)
	if err != nil {
		t.Error(err)
		return
	}
	if status != 2 {
		t.Error("wrong status")
	}
}

func TestCafeRequestDB_DeleteCompleteSyncGroups(t *testing.T) {
	err := cafeRequestStore.DeleteCompleteSyncGroups()
	if err != nil {
		t.Error(err)
	}
}

func TestCafeRequestDB_SyncGroupCompleteAgain(t *testing.T) {
	err := cafeRequestStore.UpdateStatus("abcdefg", pb.CafeRequest_COMPLETE)
	if err != nil {
		t.Error(err)
		return
	}
	if !cafeRequestStore.SyncGroupComplete("sync_group2") {
		t.Error("sync group complete failed")
	}
}

func TestCafeRequestDB_DeleteCompleteSyncGroupsAgain(t *testing.T) {
	err := cafeRequestStore.DeleteCompleteSyncGroups()
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcdef").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
	err = stmt.QueryRow("abcdefg").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_AddAttempt(t *testing.T) {
	err := cafeRequestStore.AddAttempt("abcde")
	if err != nil {
		t.Error(err)
		return
	}
	req := cafeRequestStore.Get("abcde")
	if req.Attempts != 1 {
		t.Error("wrong attempts")
	}
}

func TestCafeRequestDB_Delete(t *testing.T) {
	err := cafeRequestStore.Delete("abcde")
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteByGroup(t *testing.T) {
	setupCafeRequestDB()
	err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "xyz",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE,
		Date:      ptypes.TimestampNow(),
		Size:      1024,
		Status:    pb.CafeRequest_NEW,
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = cafeRequestStore.DeleteByGroup("group")
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("xyz").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteBySyncGroup(t *testing.T) {
	setupCafeRequestDB()
	err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "xyz",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE,
		Date:      ptypes.TimestampNow(),
		Size:      1024,
		Status:    pb.CafeRequest_NEW,
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = cafeRequestStore.DeleteBySyncGroup("sync_group")
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("xyz").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteByCafe(t *testing.T) {
	setupCafeRequestDB()
	err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "xyz",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE,
		Date:      ptypes.TimestampNow(),
		Size:      1024,
		Status:    pb.CafeRequest_NEW,
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = cafeRequestStore.DeleteByCafe("peer")
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("zyx").Scan(&id)
	if err == nil {
		t.Error("delete by cafe failed")
	}
}
