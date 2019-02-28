package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

var threadPeerStore repo.ThreadPeerStore

func init() {
	setupThreadPeerDB()
}

func setupThreadPeerDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	threadPeerStore = NewThreadPeerStore(conn, new(sync.Mutex))
}

func TestThreadPeerDB_Add(t *testing.T) {
	err := threadPeerStore.Add(&pb.ThreadPeer{
		Id:       "abc",
		Thread:   ksuid.New().String(),
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadPeerStore.PrepareQuery("select id from thread_peers where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abc").Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != "abc" {
		t.Errorf(`expected id "abc" got %s`, id)
	}
}

func TestThreadPeerDB_ListById(t *testing.T) {
	setupThreadPeerDB()
	err := threadPeerStore.Add(&pb.ThreadPeer{
		Id:       ksuid.New().String(),
		Thread:   ksuid.New().String(),
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadPeerStore.Add(&pb.ThreadPeer{
		Id:       "boo",
		Thread:   ksuid.New().String(),
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	filtered := threadPeerStore.ListById("boo")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of peers")
		return
	}
}

func TestThreadPeerDB_ListByThread(t *testing.T) {
	setupThreadPeerDB()
	err := threadPeerStore.Add(&pb.ThreadPeer{
		Id:       ksuid.New().String(),
		Thread:   "foo",
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadPeerStore.Add(&pb.ThreadPeer{
		Id:       ksuid.New().String(),
		Thread:   "boo",
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	filtered := threadPeerStore.ListByThread("boo")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of peers")
		return
	}
}

func TestThreadPeerDB_Count(t *testing.T) {
	setupThreadPeerDB()
	err := threadPeerStore.Add(&pb.ThreadPeer{
		Id:       "bar",
		Thread:   "1",
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadPeerStore.Add(&pb.ThreadPeer{
		Id:       "bar",
		Thread:   "2",
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadPeerStore.Add(&pb.ThreadPeer{
		Id:       "bar2",
		Thread:   "2",
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	cnt := threadPeerStore.Count(false)
	if cnt != 3 {
		t.Error("returned incorrect count of peers")
		return
	}
	distinct := threadPeerStore.Count(true)
	if distinct != 2 {
		t.Error("returned incorrect count of peers")
		return
	}
}

func TestThreadPeerDB_Delete(t *testing.T) {
	err := threadPeerStore.Add(&pb.ThreadPeer{
		Id:       "car",
		Thread:   "3",
		Welcomed: false,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadPeerStore.Delete("car", "3")
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadPeerStore.PrepareQuery("select id from thread_peers where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("car").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestThreadPeerDB_DeleteById(t *testing.T) {
	err := threadPeerStore.DeleteById("bar")
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadPeerStore.PrepareQuery("select id from thread_peers where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("bar").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestThreadPeerDB_DeleteByThread(t *testing.T) {
	err := threadPeerStore.DeleteByThread("2")
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadPeerStore.PrepareQuery("select id from thread_peers where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("bar2").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
