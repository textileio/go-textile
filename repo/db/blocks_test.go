package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
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
	if err := blockStore.Add(&pb.Block{
		Id:      "abcde",
		Thread:  "thread_id",
		Author:  "author_id",
		Type:    pb.Block_FILES,
		Date:    ptypes.TimestampNow(),
		Parents: []string{"Qm123"},
		Target:  "Qm456",
		Body:    "body",
	}); err != nil {
		t.Error(err)
		return
	}
	stmt, err := blockStore.PrepareQuery("select id from blocks where id=?")
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

func TestBlockDB_Get(t *testing.T) {
	if blockStore.Get("abcde") == nil {
		t.Error("could not get block")
	}
}

func TestBlockDB_List(t *testing.T) {
	setupBlockDB()
	if err := blockStore.Add(&pb.Block{
		Id:      "abcde",
		Thread:  "thread_id",
		Author:  "author_id",
		Type:    pb.Block_FILES,
		Date:    ptypes.TimestampNow(),
		Parents: []string{"Qm123"},
		Target:  "Qm456",
		Body:    "body",
	}); err != nil {
		t.Error(err)
		return
	}

	if err := blockStore.Add(&pb.Block{
		Id:      "fghijk",
		Thread:  "thread_id",
		Author:  "author_id",
		Type:    pb.Block_FILES,
		Date:    util.ProtoTs(time.Now().Add(time.Minute).UnixNano()),
		Parents: []string{"Qm456"},
		Target:  "Qm789",
		Body:    "body",
	}); err != nil {
		t.Error(err)
		return
	}

	all := blockStore.List("", -1, "").Items
	if len(all) != 2 {
		t.Error("returned incorrect number of blocks")
		return
	}

	limited := blockStore.List("", 1, "").Items
	if len(limited) != 1 {
		t.Error("returned incorrect number of blocks")
		return
	}

	offset := blockStore.List(limited[0].Id, -1, "").Items
	if len(offset) != 1 {
		t.Error("returned incorrect number of blocks")
		return
	}

	filtered := blockStore.List("", -1, "threadId='thread_id'").Items
	if len(filtered) != 2 {
		t.Error("returned incorrect number of blocks")
	}
}

func TestBlockDB_Count(t *testing.T) {
	setupBlockDB()
	if err := blockStore.Add(&pb.Block{
		Id:      "abcde",
		Thread:  "thread_id",
		Author:  "author_id",
		Type:    pb.Block_FILES,
		Date:    ptypes.TimestampNow(),
		Parents: []string{"Qm123"},
		Target:  "Qm456",
		Body:    "body",
	}); err != nil {
		t.Error(err)
		return
	}

	if err := blockStore.Add(&pb.Block{
		Id:      "abcde2",
		Thread:  "thread_id",
		Author:  "author_id",
		Date:    ptypes.TimestampNow(),
		Type:    pb.Block_FILES,
		Parents: []string{"Qm123"},
		Target:  "Qm456",
		Body:    "body",
	}); err != nil {
		t.Error(err)
		return
	}

	if blockStore.Count("") != 2 {
		t.Error("returned incorrect count of blocks")
	}
}

func TestBlockDB_Delete(t *testing.T) {
	if err := blockStore.Delete("abcde"); err != nil {
		t.Error(err)
		return
	}

	stmt, err := blockStore.PrepareQuery("select id from blocks where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()

	var id string
	if err := stmt.QueryRow("abcde").Scan(&id); err == nil {
		t.Error("delete failed")
	}
}

func TestBlockDB_DeleteByThread(t *testing.T) {
	if err := blockStore.DeleteByThread("thread_id"); err != nil {
		t.Error(err)
		return
	}
	stmt, err := blockStore.PrepareQuery("select id from blocks where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()

	var id string
	if err := stmt.QueryRow("abcde2").Scan(&id); err == nil {
		t.Error("delete by thread id failed")
	}
}
