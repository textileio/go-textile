package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

var peerStore repo.PeerStore

var testPeer *pb.Peer
var testCafe = &pb.Cafe{
	Peer:     "peer",
	Address:  "address",
	Api:      "v0",
	Protocol: "/textile/cafe/1.0.0",
	Node:     "v1.0.0",
	Url:      "https://mycafe.com",
}

func init() {
	setupPeerDB()
}

func setupPeerDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	peerStore = NewPeerStore(conn, new(sync.Mutex))
}

func TestPeerDB_Add(t *testing.T) {
	if err := peerStore.Add(&pb.Peer{
		Id:      "abcde",
		Address: "address",
	}); err != nil {
		t.Error(err)
		return
	}
	stmt, err := peerStore.PrepareQuery("select id from peers where id=?")
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

func TestPeerDB_Get(t *testing.T) {
	testPeer = peerStore.Get("abcde")
	if testPeer == nil {
		t.Error("could not get peer")
	}
}

func TestPeerDB_GetBest(t *testing.T) {
	testPeer = peerStore.GetBest("abcde")
	if testPeer == nil {
		t.Error("could not get best peer")
	}
}

func TestPeerDB_AddOrUpdate(t *testing.T) {
	testPeer.Name = "joe"
	testPeer.Avatar = "Qm123"
	testPeer.Inboxes = []*pb.Cafe{testCafe}
	if err := peerStore.AddOrUpdate(testPeer); err != nil {
		t.Error(err)
		return
	}
	stmt, err := peerStore.PrepareQuery("select username, updated from peers where id=?")
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
	old := util.ProtoNanos(testPeer.Updated)
	if updated <= old {
		t.Errorf("updated was not updated (old: %d, new: %d)", old, updated)
	}
}

func TestPeerDB_List(t *testing.T) {
	setupPeerDB()
	if err := peerStore.Add(&pb.Peer{
		Id:      "abcde",
		Address: "address1",
		Name:    "joe",
		Avatar:  "Qm123",
		Inboxes: []*pb.Cafe{testCafe},
	}); err != nil {
		t.Error(err)
		return
	}
	if err := peerStore.Add(&pb.Peer{
		Id:      "fghij",
		Address: "address2",
		Name:    "joe",
		Avatar:  "Qm123",
		Inboxes: []*pb.Cafe{testCafe, testCafe},
	}); err != nil {
		t.Error(err)
		return
	}
	list := peerStore.List("")
	if len(list) != 2 {
		t.Error("returned incorrect number of peers")
	}
}

func TestPeerDB_Count(t *testing.T) {
	if peerStore.Count("") != 2 {
		t.Error("returned incorrect count of peers")
	}
}

func TestPeerDB_UpdateName(t *testing.T) {
	if err := peerStore.UpdateName(testPeer.Id, "mike"); err != nil {
		t.Error(err)
		return
	}
	updated := peerStore.Get(testPeer.Id)
	if updated.Name != "mike" {
		t.Error("update username failed")
		return
	}
	if util.ProtoNanos(updated.Updated) <= util.ProtoNanos(testPeer.Updated) {
		t.Error("update was not updated")
	}
	testPeer = updated
}

func TestPeerDB_UpdateAvatar(t *testing.T) {
	if err := peerStore.UpdateAvatar(testPeer.Id, "avatar2"); err != nil {
		t.Error(err)
		return
	}
	updated := peerStore.Get(testPeer.Id)
	if updated.Avatar != "avatar2" {
		t.Error("update avatar failed")
		return
	}
	if util.ProtoNanos(updated.Updated) <= util.ProtoNanos(testPeer.Updated) {
		t.Error("update was not updated")
	}
	testPeer = updated
}

func TestPeerDB_UpdateInboxes(t *testing.T) {
	testCafe.Peer = "newone"
	if err := peerStore.UpdateInboxes(testPeer.Id, []*pb.Cafe{testCafe}); err != nil {
		t.Error(err)
		return
	}
	updated := peerStore.Get(testPeer.Id)
	if updated.Inboxes[0].Peer != "newone" {
		t.Error("update inboxes failed")
		return
	}
	if util.ProtoNanos(updated.Updated) <= util.ProtoNanos(testPeer.Updated) {
		t.Error("update was not updated")
	}
}

func TestPeerDB_Delete(t *testing.T) {
	if err := peerStore.Delete("abcde"); err != nil {
		t.Error(err)
	}
	stmt, err := peerStore.PrepareQuery("select id from peers where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err = stmt.QueryRow("abcde").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}
