package db

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
)

var peerdb repo.PeerStore

func init() {
	setupPeerDB()
}

func setupPeerDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	peerdb = NewPeerStore(conn, new(sync.Mutex))
}

func TestPeerDB_Add(t *testing.T) {
	err := peerdb.Add(&repo.Peer{
		Row:      "abc",
		Id:       ksuid.New().String(),
		ThreadId: ksuid.New().String(),
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := peerdb.PrepareQuery("select row from peers where row=?")
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

func TestPeerDB_Get(t *testing.T) {
	setupPeerDB()
	err := peerdb.Add(&repo.Peer{
		Row:      "abc",
		Id:       ksuid.New().String(),
		ThreadId: ksuid.New().String(),
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	p := peerdb.Get("abc")
	if p == nil {
		t.Error("could not get peer")
	}
}

func TestPeerDB_GetById(t *testing.T) {
	setupPeerDB()
	err := peerdb.Add(&repo.Peer{
		Row:      ksuid.New().String(),
		Id:       "abc",
		ThreadId: ksuid.New().String(),
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	p := peerdb.GetById("abc")
	if p == nil {
		t.Error("could not get peer")
	}
}

func TestPeerDB_List(t *testing.T) {
	setupPeerDB()
	err := peerdb.Add(&repo.Peer{
		Row:      "abc",
		Id:       ksuid.New().String(),
		ThreadId: "foo",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	err = peerdb.Add(&repo.Peer{
		Row:      "def",
		Id:       ksuid.New().String(),
		ThreadId: "boo",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	all := peerdb.List(-1, "")
	if len(all) != 2 {
		t.Error("returned incorrect number of peers")
		return
	}
	filtered := peerdb.List(-1, "threadId='boo'")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of peers")
		return
	}
}

func TestPeerDB_Count(t *testing.T) {
	setupPeerDB()
	err := peerdb.Add(&repo.Peer{
		Row:      "abc",
		Id:       "bar",
		ThreadId: "1",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	err = peerdb.Add(&repo.Peer{
		Row:      "def",
		Id:       "bar",
		ThreadId: "2",
		PubKey:   []byte(ksuid.New().String()),
	})
	if err != nil {
		t.Error(err)
	}
	cnt := peerdb.Count("", false)
	if cnt != 2 {
		t.Error("returned incorrect count of peers")
		return
	}
	distinct := peerdb.Count("", true)
	if distinct != 1 {
		t.Error("returned incorrect count of peers")
		return
	}
}

func TestPeerDB_Delete(t *testing.T) {
	err := peerdb.Delete("bar", "1")
	if err != nil {
		t.Error(err)
	}
	stmt, err := peerdb.PrepareQuery("select row from peers where row=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abc").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestPeerDB_DeleteByThreadId(t *testing.T) {
	err := peerdb.DeleteByThreadId("2")
	if err != nil {
		t.Error(err)
	}
	stmt, err := peerdb.PrepareQuery("select row from peers where row=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("def").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
