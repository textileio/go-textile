package db

import (
	"database/sql"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"sync"
	"testing"
	"time"
)

var bdb repo.BlockStore

func init() {
	setupBlockDB()
}

func setupBlockDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	bdb = NewBlockStore(conn, new(sync.Mutex))
}

func TestBlockDB_Put(t *testing.T) {
	key, err := crypto.GenerateAESKey()
	if err != nil {
		t.Error(err)
	}
	_, pk, err := libp2pc.GenerateKeyPair(libp2pc.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
	}
	err = bdb.Add(&repo.Block{
		Id:           "abcde",
		Target:       "Qm456",
		Parents:      []string{"Qm123"},
		TargetKey:    key,
		ThreadPubKey: libp2pc.ConfigEncodeKey(pkb),
		Type:         repo.PhotoBlock,
		Date:         time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := bdb.PrepareQuery("select id from blocks where id=?")
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

func TestBlockDB_Get(t *testing.T) {
	block := bdb.Get("abcde")
	if block == nil {
		t.Error("could not get block")
		return
	}
}

func TestBlockDB_List(t *testing.T) {
	setupBlockDB()
	key, err := crypto.GenerateAESKey()
	if err != nil {
		t.Error(err)
	}
	_, pk, err := libp2pc.GenerateKeyPair(libp2pc.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
	}
	err = bdb.Add(&repo.Block{
		Id:           "abcde",
		Target:       "Qm456",
		Parents:      []string{"Qm123"},
		TargetKey:    key,
		ThreadPubKey: libp2pc.ConfigEncodeKey(pkb),
		Type:         repo.PhotoBlock,
		Date:         time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
	_, pk2, err := libp2pc.GenerateKeyPair(libp2pc.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	pkb2, err := pk2.Bytes()
	if err != nil {
		t.Error(err)
	}
	err = bdb.Add(&repo.Block{
		Id:           "fghijk",
		Target:       "Qm789",
		Parents:      []string{"Qm456"},
		TargetKey:    key,
		ThreadPubKey: libp2pc.ConfigEncodeKey(pkb2),
		Type:         repo.CommentBlock,
		Date:         time.Now().Add(time.Minute),
	})
	if err != nil {
		t.Error(err)
	}
	all := bdb.List("", -1, "")
	if len(all) != 2 {
		t.Error("returned incorrect number of blocks")
		return
	}
	limited := bdb.List("", 1, "")
	if len(limited) != 1 {
		t.Error("returned incorrect number of blocks")
		return
	}
	offset := bdb.List(limited[0].Id, -1, "")
	if len(offset) != 1 {
		t.Error("returned incorrect number of blocks")
		return
	}
	filtered := bdb.List("", -1, "pk='"+libp2pc.ConfigEncodeKey(pkb2)+"'")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of blocks")
		return
	}
}

func TestBlockDB_Delete(t *testing.T) {
	err := bdb.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := bdb.PrepareQuery("select id from blocks where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("Delete failed")
	}
}
