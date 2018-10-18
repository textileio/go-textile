package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
)

var accountpeerdb repo.AccountPeerStore

func init() {
	setupAccountPeerDB()
}

func setupAccountPeerDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	accountpeerdb = NewAccountPeerStore(conn, new(sync.Mutex))
}

func TestAccountPeerDB_Add(t *testing.T) {
	err := accountpeerdb.Add(&repo.AccountPeer{
		Id:   "abcde",
		Name: "boom",
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := accountpeerdb.PrepareQuery("select id from account_peers where id=?")
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

func TestAccountPeerDB_Get(t *testing.T) {
	block := accountpeerdb.Get("abcde")
	if block == nil {
		t.Error("could not get peer")
	}
}

func TestAccountPeerDB_List(t *testing.T) {
	setupAccountPeerDB()
	err := accountpeerdb.Add(&repo.AccountPeer{
		Id:   "abcde",
		Name: "boom",
	})
	if err != nil {
		t.Error(err)
	}
	err = accountpeerdb.Add(&repo.AccountPeer{
		Id:   "abcdef",
		Name: "booom",
	})
	if err != nil {
		t.Error(err)
	}
	all := accountpeerdb.List("")
	if len(all) != 2 {
		t.Error("returned incorrect number of peers")
		return
	}
	filtered := accountpeerdb.List("name='boom'")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of peers")
	}
}

func TestAccountPeerDB_Count(t *testing.T) {
	setupAccountPeerDB()
	err := accountpeerdb.Add(&repo.AccountPeer{
		Id:   "abcde",
		Name: "hello",
	})
	if err != nil {
		t.Error(err)
	}
	cnt := accountpeerdb.Count("")
	if cnt != 1 {
		t.Error("returned incorrect count of peers")
	}
}

func TestAccountPeerDB_Delete(t *testing.T) {
	err := accountpeerdb.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := accountpeerdb.PrepareQuery("select id from account_peers where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
