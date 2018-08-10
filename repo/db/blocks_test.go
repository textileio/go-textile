package db

import (
	"crypto/rand"
	"database/sql"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"sync"
	"testing"
	"time"
)

var bdb repo.BlockStore

var threadId string

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
	_, pk, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
	}
	err = bdb.Add(&repo.Block{
		Id:                 "abcde",
		Date:               time.Now(),
		Parents:            []string{"Qm123"},
		ThreadId:           libp2pc.ConfigEncodeKey(pkb),
		AuthorPk:           "author_pk",
		Type:               repo.PhotoBlock,
		DataId:             "Qm456",
		DataKeyCipher:      key,
		DataCaptionCipher:  []byte("xxx"),
		DataUsernameCipher: []byte("un"),
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

func TestBlockDB_GetByDataId(t *testing.T) {
	block := bdb.GetByDataId("Qm456")
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
	_, pk, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	pkb, err := pk.Bytes()
	if err != nil {
		t.Error(err)
	}
	err = bdb.Add(&repo.Block{
		Id:                 "abcde",
		Date:               time.Now(),
		Parents:            []string{"Qm123"},
		ThreadId:           libp2pc.ConfigEncodeKey(pkb),
		AuthorPk:           "author_pk",
		Type:               repo.PhotoBlock,
		DataId:             "Qm456",
		DataKeyCipher:      key,
		DataCaptionCipher:  []byte("xxx"),
		DataUsernameCipher: []byte("un"),
	})
	if err != nil {
		t.Error(err)
	}
	_, pk2, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	pkb2, err := pk2.Bytes()
	if err != nil {
		t.Error(err)
	}
	threadId = libp2pc.ConfigEncodeKey(pkb2)
	err = bdb.Add(&repo.Block{
		Id:                 "fghijk",
		Date:               time.Now().Add(time.Minute),
		Parents:            []string{"Qm456"},
		ThreadId:           threadId,
		AuthorPk:           "author_pk",
		Type:               repo.PhotoBlock,
		DataId:             "Qm789",
		DataKeyCipher:      key,
		DataCaptionCipher:  []byte("xxx"),
		DataUsernameCipher: []byte("un"),
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
	filtered := bdb.List("", -1, "threadId='"+threadId+"'")
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

func TestBlockDB_DeleteByThreadId(t *testing.T) {
	err := bdb.DeleteByThreadId(threadId)
	if err != nil {
		t.Error(err)
	}
	stmt, err := bdb.PrepareQuery("select id from blocks where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("fghijk").Scan(&id)
	if err == nil {
		t.Error("Delete failed")
	}
}
