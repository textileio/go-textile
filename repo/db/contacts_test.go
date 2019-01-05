package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/textileio/textile-go/repo"
)

var contactStore repo.ContactStore

func init() {
	setupContactDB()
}

func setupContactDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	contactStore = NewContactStore(conn, new(sync.Mutex))
}

func TestContactDB_Add(t *testing.T) {
	err := contactStore.Add(&repo.Contact{
		Id:      "abcde",
		Address: "address",
		Added:   time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := contactStore.PrepareQuery("select id from contacts where id=?")
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

func TestContactDB_AddOrUpdate(t *testing.T) {
	cafe := repo.Cafe{
		Peer:     "peer",
		Address:  "address",
		API:      "v0",
		Protocol: "/textile/cafe/1.0.0",
		Node:     "v1.0.0",
		URL:      "https://mycafe.com",
	}
	err := contactStore.AddOrUpdate(&repo.Contact{
		Id:       "abcde",
		Address:  "address",
		Username: "joe",
		Inboxes:  []repo.Cafe{cafe},
		Added:    time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := contactStore.PrepareQuery("select username from contacts where id=?")
	defer stmt.Close()
	var username string
	err = stmt.QueryRow("abcde").Scan(&username)
	if err != nil {
		t.Error(err)
	}
	if username != "joe" {
		t.Errorf(`expected "joe" got %s`, username)
	}
}

func TestContactDB_Get(t *testing.T) {
	block := contactStore.Get("abcde")
	if block == nil {
		t.Error("could not get contact")
	}
}

func TestContactDB_List(t *testing.T) {
	setupContactDB()
	cafe := repo.Cafe{
		Peer:     "peer",
		Address:  "address",
		API:      "v0",
		Protocol: "/textile/cafe/1.0.0",
		Node:     "v1.0.0",
		URL:      "https://mycafe.com",
	}
	err := contactStore.Add(&repo.Contact{
		Id:       "abcde",
		Address:  "address1",
		Username: "joe",
		Inboxes:  []repo.Cafe{cafe},
		Added:    time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	err = contactStore.Add(&repo.Contact{
		Id:       "fghij",
		Address:  "address2",
		Username: "joe",
		Inboxes:  []repo.Cafe{cafe, cafe},
		Added:    time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	list := contactStore.List()
	if len(list) != 2 {
		t.Error("returned incorrect number of contacts")
		return
	}
}

func TestContactDB_ListByAddress(t *testing.T) {
	list := contactStore.ListByAddress("address1")
	if len(list) != 1 {
		t.Error("returned incorrect number of contacts")
		return
	}
}

func TestContactDB_Count(t *testing.T) {
	if contactStore.Count() != 2 {
		t.Error("returned incorrect count of contacts")
	}
}

func TestContactDB_Delete(t *testing.T) {
	err := contactStore.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := contactStore.PrepareQuery("select id from contacts where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
