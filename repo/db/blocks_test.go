package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
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

func TestBlockDB_Add(t *testing.T) {
	err := bdb.Add(&repo.Block{
		Id:       "abcde",
		ThreadId: "thread_id",
		AuthorId: "author_id",
		Type:     repo.FilesBlock,
		Date:     time.Now(),
		Parents:  []string{"Qm123"},
		Target:   "Qm456",
		Body:     "body",
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
	}
}

func TestBlockDB_GetByTarget(t *testing.T) {
	block := bdb.GetByTarget("Qm456")
	if block == nil {
		t.Error("could not get block")
	}
}

func TestBlockDB_List(t *testing.T) {
	setupBlockDB()
	err := bdb.Add(&repo.Block{
		Id:       "abcde",
		ThreadId: "thread_id",
		AuthorId: "author_id",
		Type:     repo.FilesBlock,
		Date:     time.Now(),
		Parents:  []string{"Qm123"},
		Target:   "Qm456",
		Body:     "body",
	})
	if err != nil {
		t.Error(err)
	}
	err = bdb.Add(&repo.Block{
		Id:       "fghijk",
		ThreadId: "thread_id",
		AuthorId: "author_id",
		Type:     repo.FilesBlock,
		Date:     time.Now().Add(time.Minute),
		Parents:  []string{"Qm456"},
		Target:   "Qm789",
		Body:     "body",
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
	filtered := bdb.List("", -1, "threadId='thread_id'")
	if len(filtered) != 2 {
		t.Error("returned incorrect number of blocks")
	}
}

func TestBlockDB_Count(t *testing.T) {
	setupBlockDB()
	err := bdb.Add(&repo.Block{
		Id:       "abcde",
		ThreadId: "thread_id",
		AuthorId: "author_id",
		Type:     repo.FilesBlock,
		Date:     time.Now(),
		Parents:  []string{"Qm123"},
		Target:   "Qm456",
		Body:     "body",
	})
	if err != nil {
		t.Error(err)
	}
	err = bdb.Add(&repo.Block{
		Id:       "abcde2",
		ThreadId: "thread_id",
		AuthorId: "author_id",
		Date:     time.Now(),
		Type:     repo.FilesBlock,
		Parents:  []string{"Qm123"},
		Target:   "Qm456",
		Body:     "body",
	})
	if err != nil {
		t.Error(err)
	}
	cnt := bdb.Count("")
	if cnt != 2 {
		t.Error("returned incorrect count of blocks")
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
		t.Error("delete failed")
	}
}

func TestBlockDB_DeleteByThread(t *testing.T) {
	err := bdb.DeleteByThread("thread_id")
	if err != nil {
		t.Error(err)
	}
	stmt, err := bdb.PrepareQuery("select id from blocks where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde2").Scan(&id)
	if err == nil {
		t.Error("delete by thread id failed")
	}
}
