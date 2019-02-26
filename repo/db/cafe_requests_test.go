package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
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
	err := cafeRequestStore.Add(&pb.CafeRequest{
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
		Type: pb.CafeRequest_STORE,
		Date: ptypes.TimestampNow(),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
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
	err := cafeRequestStore.Add(&pb.CafeRequest{
		Id:     "abcde",
		Peer:   "peer",
		Target: "zxy",
		Cafe:   cafe,
		Type:   pb.CafeRequest_STORE_THREAD,
		Date:   ptypes.TimestampNow(),
	})
	if err != nil {
		t.Error(err)
	}
	err = cafeRequestStore.Add(&pb.CafeRequest{
		Id:     "abcdef",
		Peer:   "peer",
		Target: "zxy",
		Cafe:   cafe,
		Type:   pb.CafeRequest_STORE,
		Date:   util.ProtoTs(time.Now().Add(time.Minute).UnixNano()),
	})
	if err != nil {
		t.Error(err)
	}
	all := cafeRequestStore.List("", -1)
	if len(all) != 2 {
		t.Error("returned incorrect number of requests")
		return
	}
	limited := cafeRequestStore.List("", 1)
	if len(limited) != 1 {
		t.Error("returned incorrect number of requests")
		return
	}
	offset := cafeRequestStore.List(limited[0].Id, -1)
	if len(offset) != 1 {
		t.Error("returned incorrect number of requests")
		return
	}
}

func TestCafeRequestDB_Delete(t *testing.T) {
	err := cafeRequestStore.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("abcde").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}

func TestCafeRequestDB_DeleteByCafe(t *testing.T) {
	setupCafeRequestDB()
	err := cafeRequestStore.Add(&pb.CafeRequest{
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
		Type: pb.CafeRequest_STORE,
		Date: ptypes.TimestampNow(),
	})
	if err != nil {
		t.Error(err)
	}
	err = cafeRequestStore.DeleteByCafe("boom")
	if err != nil {
		t.Error(err)
	}
	stmt, err := cafeRequestStore.PrepareQuery("select id from cafe_requests where id=?")
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("zyx").Scan(&id); err == nil {
		t.Error("delete by cafe failed")
	}
}
