package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/textileio/textile-go/repo"
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
	err := cafeRequestStore.Add(&repo.CafeRequest{
		Id:       "abcde",
		PeerId:   "peer",
		TargetId: "zxy",
		Cafe: repo.Cafe{
			Peer:     "peer",
			Address:  "address",
			API:      "v0",
			Protocol: "/textile/cafe/1.0.0",
			Node:     "v1.0.0",
			URL:      "https://mycafe.com",
		},
		Type: repo.CafeStoreRequest,
		Date: time.Now(),
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
	cafe := repo.Cafe{
		Peer:     "peer",
		Address:  "address",
		API:      "v0",
		Protocol: "/textile/cafe/1.0.0",
		Node:     "v1.0.0",
		URL:      "https://mycafe.com",
	}
	err := cafeRequestStore.Add(&repo.CafeRequest{
		Id:       "abcde",
		PeerId:   "peer",
		TargetId: "zxy",
		Cafe:     cafe,
		Type:     repo.CafeStoreThreadRequest,
		Date:     time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	err = cafeRequestStore.Add(&repo.CafeRequest{
		Id:       "abcdef",
		PeerId:   "peer",
		TargetId: "zxy",
		Cafe:     cafe,
		Type:     repo.CafeStoreRequest,
		Date:     time.Now().Add(time.Minute),
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
	err := cafeRequestStore.Add(&repo.CafeRequest{
		Id:       "xyz",
		PeerId:   "peer",
		TargetId: "zxy",
		Cafe: repo.Cafe{
			Peer:     "peer",
			Address:  "address",
			API:      "v0",
			Protocol: "/textile/cafe/1.0.0",
			Node:     "v1.0.0",
			URL:      "https://mycafe.com",
		},
		Type: repo.CafeStoreRequest,
		Date: time.Now(),
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
