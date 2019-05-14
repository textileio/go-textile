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
	initDatabaseTables(conn, "")
	cafeRequestStore = NewCafeRequestStore(conn, new(sync.Mutex))
}

func TestCafeRequestDB_Add(t *testing.T) {
	if err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "abcde",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE,
		Date:      ptypes.TimestampNow(),
		Size:      1024,
		Status:    pb.CafeRequest_NEW,
		Attempts:  0,
	}); err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("abcde").Scan(&id); err != nil {
		t.Error(err)
	}
	if id != "abcde" {
		t.Errorf(`expected "abcde" got %s`, id)
	}
}

func TestCafeRequestDB_Get(t *testing.T) {
	setupCafeRequestDB()
	if err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "abcde",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE_THREAD,
		Date:      ptypes.TimestampNow(),
		Status:    pb.CafeRequest_NEW,
	}); err != nil {
		t.Error(err)
	}
	req := cafeRequestStore.Get("abcde")
	if req == nil {
		t.Error("get request failed")
		return
	}
}

func TestCafeRequestDB_List(t *testing.T) {
	setupCafeRequestDB()
	if err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:        "abcde",
		Peer:      "peer",
		Target:    "zxy",
		Cafe:      testCafe,
		Group:     "group1",
		SyncGroup: "sync_group",
		Type:      pb.CafeRequest_STORE_THREAD,
		Date:      ptypes.TimestampNow(),
		Status:    pb.CafeRequest_NEW,
	}); err != nil {
		t.Error(err)
	}
	if err := cafeRequestStore.Add(&pb.CafeRequest{
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
	}); err != nil {
		t.Error(err)
	}
	if err := cafeRequestStore.Add(&pb.CafeRequest{
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
	}); err != nil {
		t.Error(err)
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

func TestCafeRequestDB_ListCompletedSyncGroups(t *testing.T) {
	list := cafeRequestStore.ListCompletedSyncGroups()
	if len(list) != 0 {
		t.Error("list completed groups failed")
	}
}

func TestCafeRequestDB_UpdateStatus(t *testing.T) {
	if err := cafeRequestStore.UpdateStatus("abcdef", pb.CafeRequest_PENDING); err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select status from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var status int
	if err := stmt.QueryRow("abcdef").Scan(&status); err != nil {
		t.Error(err)
	}
	if status != 1 {
		t.Error("wrong status")
	}
}

func TestCafeRequestDB_SyncGroupStatus(t *testing.T) {
	status := cafeRequestStore.SyncGroupStatus("sync_group2")
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

func TestCafeRequestDB_UpdateGroupStatus(t *testing.T) {
	if err := cafeRequestStore.UpdateGroupStatus("group2", pb.CafeRequest_COMPLETE); err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select status from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var status int
	if err := stmt.QueryRow("abcdef").Scan(&status); err != nil {
		t.Error(err)
	}
	if status != 2 {
		t.Error("wrong status")
	}
}

func TestCafeRequestDB_ListCompletedSyncGroupsAgain(t *testing.T) {
	list := cafeRequestStore.ListCompletedSyncGroups()
	if len(list) != 0 {
		t.Error("list completed groups failed")
	}
	if err := cafeRequestStore.UpdateStatus("abcdefg", pb.CafeRequest_COMPLETE); err != nil {
		t.Error(err)
	}
	list = cafeRequestStore.ListCompletedSyncGroups()
	if len(list) != 1 {
		t.Error("list completed groups failed")
	}
}

func TestCafeRequestDB_Delete(t *testing.T) {
	if err := cafeRequestStore.Delete("abcde"); err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("abcde").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteByGroup(t *testing.T) {
	if err := cafeRequestStore.DeleteByGroup("group2"); err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("abcdef").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteBySyncGroup(t *testing.T) {
	if err := cafeRequestStore.DeleteBySyncGroup("sync_group2"); err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("abcdefg").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteByCafe(t *testing.T) {
	setupCafeRequestDB()
	if err := cafeRequestStore.Add(&pb.CafeRequest{
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
	}); err != nil {
		t.Error(err)
	}
	if err := cafeRequestStore.DeleteByCafe("peer"); err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("zyx").Scan(&id); err == nil {
		t.Error("delete by cafe failed")
	}
}
