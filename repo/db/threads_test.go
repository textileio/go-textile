package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
)

var tdb repo.ThreadStore

func init() {
	setupThreadDB()
}

func setupThreadDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	tdb = NewThreadStore(conn, new(sync.Mutex))
}

func TestThreadDB_Add(t *testing.T) {
	err := tdb.Add(&repo.Thread{
		Id:      "Qmabc123",
		Name:    "boom",
		PrivKey: make([]byte, 8),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := tdb.PrepareQuery("select id from threads where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("Qmabc123").Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != "Qmabc123" {
		t.Errorf(`expected "Qmabc123" got %s`, id)
	}
}

func TestThreadDB_Get(t *testing.T) {
	setupThreadDB()
	err := tdb.Add(&repo.Thread{
		Id:      "Qmabc",
		Name:    "boom",
		PrivKey: make([]byte, 8),
	})
	if err != nil {
		t.Error(err)
	}
	th := tdb.Get("Qmabc")
	if th == nil {
		t.Error("could not get thread")
	}
}

func TestThreadDB_List(t *testing.T) {
	setupThreadDB()
	err := tdb.Add(&repo.Thread{
		Id:      "Qm123",
		Name:    "boom",
		PrivKey: make([]byte, 8),
	})
	if err != nil {
		t.Error(err)
	}
	err = tdb.Add(&repo.Thread{
		Id:      "Qm456",
		Name:    "boom2",
		PrivKey: make([]byte, 8),
	})
	if err != nil {
		t.Error(err)
	}
	all := tdb.List()
	if len(all) != 2 {
		t.Error("returned incorrect number of threads")
		return
	}
}

func TestThreadDB_Count(t *testing.T) {
	setupThreadDB()
	err := tdb.Add(&repo.Thread{
		Id:      "Qm123count",
		Name:    "boom",
		PrivKey: make([]byte, 8),
	})
	if err != nil {
		t.Error(err)
	}
	cnt := tdb.Count()
	if cnt != 1 {
		t.Error("returned incorrect count of threads")
		return
	}
}

func TestThreadDB_UpdateHead(t *testing.T) {
	setupThreadDB()
	err := tdb.Add(&repo.Thread{
		Id:      "Qmabc",
		Name:    "boom",
		PrivKey: make([]byte, 8),
	})
	if err != nil {
		t.Error(err)
	}
	err = tdb.UpdateHead("Qmabc", "12345")
	if err != nil {
		t.Error(err)
	}
	th := tdb.Get("Qmabc")
	if th == nil {
		t.Error("could not get thread")
	}
	if th.Head != "12345" {
		t.Error("update head failed")
	}
}

func TestThreadDB_Delete(t *testing.T) {
	setupThreadDB()
	err := tdb.Add(&repo.Thread{
		Id:      "Qm789",
		Name:    "boom",
		PrivKey: make([]byte, 8),
	})
	if err != nil {
		t.Error(err)
	}
	all := tdb.List()
	if len(all) == 0 {
		t.Error("returned incorrect number of threads")
		return
	}
	err = tdb.Delete(all[0].Id)
	if err != nil {
		t.Error(err)
	}
	stmt, err := tdb.PrepareQuery("select id from threads where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow(all[0].Id).Scan(&id)
	if err == nil {
		t.Error("Delete failed")
	}
}
