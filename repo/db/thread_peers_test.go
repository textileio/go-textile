package db

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
)

var threadpeerdb repo.ThreadPeerStore

func init() {
	setupThreadPeerDB()
}

func setupThreadPeerDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	threadpeerdb = NewThreadPeerStore(conn, new(sync.Mutex))
}

func TestThreadPeerDB_Add(t *testing.T) {
	err := threadpeerdb.Add(&repo.ThreadPeer{
		Row:      "abc",
		Id:       ksuid.New().String(),
		ThreadId: ksuid.New().String(),
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadpeerdb.PrepareQuery("select row from thread_peers where row=?")
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

func TestThreadPeerDB_Get(t *testing.T) {
	setupThreadPeerDB()
	err := threadpeerdb.Add(&repo.ThreadPeer{
		Row:      "abc",
		Id:       ksuid.New().String(),
		ThreadId: ksuid.New().String(),
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	p := threadpeerdb.Get("abc")
	if p == nil {
		t.Error("could not get peer")
	}
}

func TestThreadPeerDB_GetById(t *testing.T) {
	setupThreadPeerDB()
	err := threadpeerdb.Add(&repo.ThreadPeer{
		Row:      ksuid.New().String(),
		Id:       "abc",
		ThreadId: ksuid.New().String(),
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	p := threadpeerdb.GetById("abc")
	if p == nil {
		t.Error("could not get peer")
	}
}

func TestThreadPeerDB_List(t *testing.T) {
	setupThreadPeerDB()
	err := threadpeerdb.Add(&repo.ThreadPeer{
		Row:      "abc",
		Id:       ksuid.New().String(),
		ThreadId: "foo",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	err = threadpeerdb.Add(&repo.ThreadPeer{
		Row:      "def",
		Id:       ksuid.New().String(),
		ThreadId: "boo",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	all := threadpeerdb.List(-1, "")
	if len(all) != 2 {
		t.Error("returned incorrect number of peers")
		return
	}
	filtered := threadpeerdb.List(-1, "threadId='boo'")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of peers")
		return
	}
}

func TestThreadPeerDB_Count(t *testing.T) {
	setupThreadPeerDB()
	err := threadpeerdb.Add(&repo.ThreadPeer{
		Row:      "abc",
		Id:       "bar",
		ThreadId: "1",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	err = threadpeerdb.Add(&repo.ThreadPeer{
		Row:      "def",
		Id:       "bar",
		ThreadId: "2",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	cnt := threadpeerdb.Count("", false)
	if cnt != 2 {
		t.Error("returned incorrect count of peers")
		return
	}
	distinct := threadpeerdb.Count("", true)
	if distinct != 1 {
		t.Error("returned incorrect count of peers")
		return
	}
}

func TestThreadPeerDB_Delete(t *testing.T) {
	err := threadpeerdb.Delete("bar", "1")
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadpeerdb.PrepareQuery("select row from thread_peers where row=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abc").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestThreadPeerDB_DeleteByThreadId(t *testing.T) {
	err := threadpeerdb.DeleteByThreadId("2")
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadpeerdb.PrepareQuery("select row from thread_peers where row=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("def").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
