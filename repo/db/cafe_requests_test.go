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
		Id:     "abcde",
		Peer:   "peer",
		Target: "zxy",
		Cafe: &pb.Cafe{
			Peer:     "peer",
			Address:  "address",
			Api:      "v0",
			Protocol: "/textile/cafe/1.0.0",
			Node:     "v1.0.0",
			Url:      "https://mycafe.com",
		},
		Type:  pb.CafeRequest_STORE,
		Date:  ptypes.TimestampNow(),
		Size:  1024,
		Group: "group",
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
		Id:     "abcde",
		Peer:   "peer",
		Target: "zxy",
		Cafe: &pb.Cafe{
			Peer:     "peer",
			Address:  "address",
			Api:      "v0",
			Protocol: "/textile/cafe/1.0.0",
			Node:     "v1.0.0",
			Url:      "https://mycafe.com",
		},
		Type:   pb.CafeRequest_STORE_THREAD,
		Date:   ptypes.TimestampNow(),
		Group:  "group",
		Status: pb.CafeRequest_NEW,
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
	cafe := &pb.Cafe{
		Peer:     "peer",
		Address:  "address",
		Api:      "v0",
		Protocol: "/textile/cafe/1.0.0",
		Node:     "v1.0.0",
		Url:      "https://mycafe.com",
	}
	if err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:     "abcde",
		Peer:   "peer",
		Target: "zxy",
		Cafe:   cafe,
		Type:   pb.CafeRequest_STORE_THREAD,
		Date:   ptypes.TimestampNow(),
		Group:  "group1",
		Status: pb.CafeRequest_NEW,
	}); err != nil {
		t.Error(err)
	}
	if err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:     "abcdef",
		Peer:   "peer",
		Target: "zxy",
		Cafe:   cafe,
		Type:   pb.CafeRequest_STORE,
		Date:   util.ProtoTs(time.Now().Add(time.Minute).UnixNano()),
		Size:   1024,
		Group:  "group2",
		Status: pb.CafeRequest_NEW,
	}); err != nil {
		t.Error(err)
	}
	if err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:     "abcdefg",
		Peer:   "peer",
		Target: "zxy",
		Cafe:   cafe,
		Type:   pb.CafeRequest_STORE,
		Date:   util.ProtoTs(time.Now().Add(time.Minute * 2).UnixNano()),
		Size:   1024,
		Group:  "group2",
		Status: pb.CafeRequest_PENDING,
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

func TestCafeRequestDB_ListCompletedGroups(t *testing.T) {
	list := cafeRequestStore.ListCompletedGroups()
	if len(list) != 0 {
		t.Error("list completed groups failed")
	}
}

func TestCafeRequestDB_CountByGroup(t *testing.T) {
	cnt := cafeRequestStore.CountByGroup("group1")
	if cnt != 1 {
		t.Error("count by group failed")
	}
}

func TestCafeRequestDB_GroupStatus(t *testing.T) {
	status := cafeRequestStore.GroupStatus("group2")
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

func TestCafeRequestDB_UpdateStatus(t *testing.T) {
	if err := cafeRequestStore.UpdateStatus("abcdef", pb.CafeRequest_COMPLETE); err != nil {
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

func TestCafeRequestDB_ListCompletedGroupsAgain(t *testing.T) {
	list := cafeRequestStore.ListCompletedGroups()
	if len(list) != 0 {
		t.Error("list completed groups failed")
	}
	if err := cafeRequestStore.UpdateStatus("abcdefg", pb.CafeRequest_COMPLETE); err != nil {
		t.Error(err)
	}
	list = cafeRequestStore.ListCompletedGroups()
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
	if err := stmt.QueryRow("abcdefg").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteByCafe(t *testing.T) {
	setupCafeRequestDB()
	if err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:     "xyz",
		Peer:   "peer",
		Target: "zxy",
		Cafe: &pb.Cafe{
			Peer:     "peer",
			Address:  "address",
			Api:      "v0",
			Protocol: "/textile/cafe/1.0.0",
			Node:     "v1.0.0",
			Url:      "https://mycafe.com",
		},
		Type:   pb.CafeRequest_STORE,
		Date:   ptypes.TimestampNow(),
		Size:   1024,
		Group:  "group",
		Status: pb.CafeRequest_NEW,
	}); err != nil {
		t.Error(err)
	}
	if err := cafeRequestStore.DeleteByCafe("boom"); err != nil {
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
