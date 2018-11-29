package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/textileio/textile-go/repo"
)

var blockStore repo.BlockStore

func init() {
	setupBlockDB()
}

func setupBlockDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	blockStore = NewBlockStore(conn, new(sync.Mutex))
}

func TestBlockDB_Add(t *testing.T) {
	err := blockStore.Add(&repo.Block{
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
	stmt, err := blockStore.PrepareQuery("select id from blocks where id=?")
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
	block := blockStore.Get("abcde")
	if block == nil {
		t.Error("could not get block")
	}
}

func TestBlockDB_List(t *testing.T) {
	setupBlockDB()
	err := blockStore.Add(&repo.Block{
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
	err = blockStore.Add(&repo.Block{
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
	all := blockStore.List("", -1, "")
	if len(all) != 2 {
		t.Error("returned incorrect number of blocks")
		return
	}
	limited := blockStore.List("", 1, "")
	if len(limited) != 1 {
		t.Error("returned incorrect number of blocks")
		return
	}
	offset := blockStore.List(limited[0].Id, -1, "")
	if len(offset) != 1 {
		t.Error("returned incorrect number of blocks")
		return
	}
	filtered := blockStore.List("", -1, "threadId='thread_id'")
	if len(filtered) != 2 {
		t.Error("returned incorrect number of blocks")
	}
}

func TestBlockDB_Count(t *testing.T) {
	setupBlockDB()
	err := blockStore.Add(&repo.Block{
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
	err = blockStore.Add(&repo.Block{
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
	cnt := blockStore.Count("")
	if cnt != 2 {
		t.Error("returned incorrect count of blocks")
	}
}

func TestBlockDB_Delete(t *testing.T) {
	err := blockStore.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := blockStore.PrepareQuery("select id from blocks where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}

func TestBlockDB_DeleteByThread(t *testing.T) {
	err := blockStore.DeleteByThread("thread_id")
	if err != nil {
		t.Error(err)
	}
	stmt, err := blockStore.PrepareQuery("select id from blocks where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde2").Scan(&id)
	if err == nil {
		t.Error("delete by thread id failed")
	}
}
