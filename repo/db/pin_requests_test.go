package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
	"time"
)

var pinreqdb repo.PinRequestStore

func init() {
	setupPinRequestDB()
}

func setupPinRequestDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	pinreqdb = NewPinRequestStore(conn, new(sync.Mutex))
}

func TestPinRequestDB_Put(t *testing.T) {
	err := pinreqdb.Put(&repo.StoreRequest{
		Id:   "abcde",
		Date: time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := pinreqdb.PrepareQuery("select id from pinrequests where id=?")
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

func TestPinRequestDB_List(t *testing.T) {
	setupPinRequestDB()
	err := pinreqdb.Put(&repo.StoreRequest{
		Id:   "abcde",
		Date: time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	err = pinreqdb.Put(&repo.StoreRequest{
		Id:   "abcdef",
		Date: time.Now().Add(time.Minute),
	})
	if err != nil {
		t.Error(err)
	}
	all := pinreqdb.List("", -1)
	if len(all) != 2 {
		t.Error("returned incorrect number of pin requests")
		return
	}
	limited := pinreqdb.List("", 1)
	if len(limited) != 1 {
		t.Error("returned incorrect number of pin requests")
		return
	}
	offset := pinreqdb.List(limited[0].Id, -1)
	if len(offset) != 1 {
		t.Error("returned incorrect number of pin requests")
		return
	}
}

func TestPinRequestDB_Delete(t *testing.T) {
	err := pinreqdb.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := pinreqdb.PrepareQuery("select id from pinrequests where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
