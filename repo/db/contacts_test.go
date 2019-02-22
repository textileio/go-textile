package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
)

var contactStore repo.ContactStore

var testContact *pb.Contact
var testCafe = &pb.Cafe{
	Peer:     "peer",
	Address:  "address",
	Api:      "v0",
	Protocol: "/textile/cafe/1.0.0",
	Node:     "v1.0.0",
	Url:      "https://mycafe.com",
}

func init() {
	setupContactDB()
}

func setupContactDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	contactStore = NewContactStore(conn, new(sync.Mutex))
}

func TestContactDB_Add(t *testing.T) {
	if err := contactStore.Add(&pb.Contact{
		Id:      "abcde",
		Address: "address",
	}); err != nil {
		t.Error(err)
		return
	}
	stmt, err := contactStore.PrepareQuery("select id from contacts where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("abcde").Scan(&id); err != nil {
		t.Error(err)
		return
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

func TestContactDB_GetBest(t *testing.T) {
	testContact = contactStore.GetBest("abcde")
	if testContact == nil {
		t.Error("could not get best contact")
	}
}

func TestContactDB_AddOrUpdate(t *testing.T) {
	testContact.Username = "joe"
	testContact.Avatar = "Qm123"
	testContact.Inboxes = []*pb.Cafe{testCafe}
	if err := contactStore.AddOrUpdate(testContact); err != nil {
		t.Error(err)
		return
	}
	stmt, err := contactStore.PrepareQuery("select username, updated from contacts where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var username string
	var updated int64
	if err := stmt.QueryRow("abcde").Scan(&username, &updated); err != nil {
		t.Error(err)
		return
	}
	if username != "joe" {
		t.Errorf(`expected "joe" got %s`, username)
		return
	}
	old := util.ProtoNanos(testContact.Updated)
	if updated <= old {
		t.Errorf("updated was not updated (old: %d, new: %d)", old, updated)
	}
}

func TestContactDB_List(t *testing.T) {
	setupContactDB()
	if err := contactStore.Add(&pb.Contact{
		Id:       "abcde",
		Address:  "address1",
		Username: "joe",
		Avatar:   "Qm123",
		Inboxes:  []*pb.Cafe{testCafe},
	}); err != nil {
		t.Error(err)
		return
	}
	if err := contactStore.Add(&pb.Contact{
		Id:       "fghij",
		Address:  "address2",
		Username: "joe",
		Avatar:   "Qm123",
		Inboxes:  []*pb.Cafe{testCafe, testCafe},
	}); err != nil {
		t.Error(err)
		return
	}
	list := contactStore.List("")
	if len(list.Items) != 2 {
		t.Error("returned incorrect number of contacts")
	}
}

func TestContactDB_Count(t *testing.T) {
	if contactStore.Count("") != 2 {
		t.Error("returned incorrect count of contacts")
	}
}

func TestContactDB_UpdateUsername(t *testing.T) {
	if err := contactStore.UpdateUsername(testContact.Id, "mike"); err != nil {
		t.Error(err)
		return
	}
	updated := contactStore.Get(testContact.Id)
	if updated.Username != "mike" {
		t.Error("update username failed")
		return
	}
	if util.ProtoNanos(updated.Updated) <= util.ProtoNanos(testContact.Updated) {
		t.Error("update was not updated")
	}
	testContact = updated
}

func TestContactDB_UpdateAvatar(t *testing.T) {
	if err := contactStore.UpdateAvatar(testContact.Id, "avatar2"); err != nil {
		t.Error(err)
		return
	}
	updated := contactStore.Get(testContact.Id)
	if updated.Avatar != "avatar2" {
		t.Error("update avatar failed")
		return
	}
	if util.ProtoNanos(updated.Updated) <= util.ProtoNanos(testContact.Updated) {
		t.Error("update was not updated")
	}
	testContact = updated
}

func TestContactDB_UpdateInboxes(t *testing.T) {
	testCafe.Peer = "newone"
	if err := contactStore.UpdateInboxes(testContact.Id, []*pb.Cafe{testCafe}); err != nil {
		t.Error(err)
		return
	}
	updated := contactStore.Get(testContact.Id)
	if updated.Inboxes[0].Peer != "newone" {
		t.Error("update inboxes failed")
		return
	}
	if util.ProtoNanos(updated.Updated) <= util.ProtoNanos(testContact.Updated) {
		t.Error("update was not updated")
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
