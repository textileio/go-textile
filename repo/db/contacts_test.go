package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/textileio/textile-go/repo"
)

var contactStore repo.ContactStore

var testContact *repo.Contact

func init() {
	setupContactDB()
}

func setupContactDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	contactStore = NewContactStore(conn, new(sync.Mutex))
}

func TestContactDB_Add(t *testing.T) {
	if err := contactStore.Add(&repo.Contact{
		Id:      "abcde",
		Address: "address",
	}); err != nil {
		t.Error(err)
	}
	stmt, err := contactStore.PrepareQuery("select id from contacts where id=?")
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

func TestContactDB_Get(t *testing.T) {
	testContact = contactStore.Get("abcde")
	if testContact == nil {
		t.Error("could not get contact")
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
	testContact.Username = "joe"
	testContact.Avatar = "Qm123"
	testContact.Inboxes = []repo.Cafe{cafe}
	if err := contactStore.AddOrUpdate(testContact); err != nil {
		t.Error(err)
	}
	stmt, err := contactStore.PrepareQuery("select username, updated from contacts where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var username string
	var updated int
	if err := stmt.QueryRow("abcde").Scan(&username, &updated); err != nil {
		t.Error(err)
	}
	if username != "joe" {
		t.Errorf(`expected "joe" got %s`, username)
	}
	old := int(testContact.Updated.UnixNano())
	if updated <= old {
		t.Errorf("updated was not updated (old: %d, new: %d)", old, updated)
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
	if err := contactStore.Add(&repo.Contact{
		Id:       "abcde",
		Address:  "address1",
		Username: "joe",
		Avatar:   "Qm123",
		Inboxes:  []repo.Cafe{cafe},
	}); err != nil {
		t.Error(err)
	}
	if err := contactStore.Add(&repo.Contact{
		Id:       "fghij",
		Address:  "address2",
		Username: "joe",
		Avatar:   "Qm123",
		Inboxes:  []repo.Cafe{cafe, cafe},
	}); err != nil {
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
	if err := contactStore.Delete("abcde"); err != nil {
		t.Error(err)
	}
	stmt, err := contactStore.PrepareQuery("select id from contacts where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err = stmt.QueryRow("abcde").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}
