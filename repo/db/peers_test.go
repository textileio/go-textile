package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

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
	_ = initDatabaseTables(conn, "")
	peerStore = NewPeerStore(conn, new(sync.Mutex))
}

func TestPeerDB_Add(t *testing.T) {
	err := peerStore.Add(&pb.Peer{
		Id:      "abcde",
		Address: "address",
	})
	if err != nil {
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
	err = stmt.QueryRow("abcde").Scan(&id)
	if err != nil {
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

func TestPeerDB_GetBestUser(t *testing.T) {
	best := peerStore.GetBestUser("abcde")
	if best == nil {
		t.Error("could not get best user")
	}
}

func TestPeerDB_AddOrUpdate(t *testing.T) {
	testPeer.Name = "joe"
	testPeer.Avatar = "Qm123"
	testPeer.Inboxes = []*pb.Cafe{testCafe}
	err := peerStore.AddOrUpdate(testPeer)
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := peerStore.PrepareQuery("select username, updated from peers where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var name string
	var updated int64
	err = stmt.QueryRow("abcde").Scan(&name, &updated)
	if err != nil {
		t.Error(err)
		return
	}
	if name != "joe" {
		t.Errorf(`expected "joe" got %s`, name)
		return
	}
	old := util.ProtoNanos(testPeer.Updated)
	if updated <= old {
		t.Errorf("updated was not updated (old: %d, new: %d)", old, updated)
	}
}

func TestPeerDB_GetBestUserAgain(t *testing.T) {
	setupPeerDB()
	now := time.Now().UnixNano()
	err := peerStore.Add(&pb.Peer{
		Id:      "abcde",
		Address: "address",
		Updated: util.ProtoTs(now),
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = peerStore.Add(&pb.Peer{
		Id:      "abcdef",
		Address: "address",
		Name:    "name",
		Avatar:  "avatar",
		Updated: util.ProtoTs(now + 1e9),
	})
	if err != nil {
		t.Error(err)
		return
	}
	best := peerStore.GetBestUser("abcde")
	if best.Address != "address" {
		t.Error("wrong address")
	}
	if best.Name != "name" {
		t.Error("wrong name")
	}
	if best.Avatar != "avatar" {
		t.Error("wrong avatar")
	}
	err = peerStore.Add(&pb.Peer{
		Id:      "abcdefg",
		Address: "address",
		Name:    "new",
		Avatar:  "new",
		Updated: util.ProtoTs(now + 2e9),
	})
	if err != nil {
		t.Error(err)
		return
	}
	best = peerStore.GetBestUser("abcde")
	if best.Address != "address" {
		t.Error("wrong address")
	}
	if best.Name != "new" {
		t.Error("wrong name")
	}
	if best.Avatar != "new" {
		t.Error("wrong avatar")
	}
}

func TestPeerDB_List(t *testing.T) {
	setupPeerDB()
	err := peerStore.Add(&pb.Peer{
		Id:      "abcde",
		Address: "address1",
		Name:    "joe",
		Avatar:  "Qm123",
		Inboxes: []*pb.Cafe{testCafe},
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = peerStore.Add(&pb.Peer{
		Id:      "fghij",
		Address: "address2",
		Name:    "joe",
		Avatar:  "Qm123",
		Inboxes: []*pb.Cafe{testCafe, testCafe},
	})
	if err != nil {
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
	err := peerStore.UpdateName(testPeer.Id, "mike")
	if err != nil {
		t.Error(err)
		return
	}
	updated := peerStore.Get(testPeer.Id)
	if updated.Name != "mike" {
		t.Error("update name failed")
		return
	}
	if util.ProtoNanos(updated.Updated) <= util.ProtoNanos(testPeer.Updated) {
		t.Error("update was not updated")
	}
	testPeer = updated
}

func TestPeerDB_UpdateAvatar(t *testing.T) {
	err := peerStore.UpdateAvatar(testPeer.Id, "avatar2")
	if err != nil {
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
	err := peerStore.UpdateInboxes(testPeer.Id, []*pb.Cafe{testCafe})
	if err != nil {
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
	err := peerStore.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := peerStore.PrepareQuery("select id from peers where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
