package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
)

var threadStore repo.ThreadStore

func init() {
	setupThreadDB()
}

func setupThreadDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	threadStore = NewThreadStore(conn, new(sync.Mutex))
}

func TestThreadDB_Add(t *testing.T) {
	err := threadStore.Add(&repo.Thread{
		Id:        "Qmabc123",
		Key:       ksuid.New().String(),
		PrivKey:   make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      repo.OpenThread,
		Members:   []string{"P1,P2"},
		Sharing:   repo.SharedThread,
		State:     repo.ThreadLoaded,
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadStore.PrepareQuery("select id from threads where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("Qmabc123").Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != "Qmabc123" {
		t.Errorf(`expected "Qmabc123" got %s`, id)
	}
}

func TestThreadDB_Get(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&repo.Thread{
		Id:        "Qmabc",
		Key:       ksuid.New().String(),
		PrivKey:   make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      repo.OpenThread,
		Members:   []string{},
		Sharing:   repo.SharedThread,
		State:     repo.ThreadLoaded,
	})
	if err != nil {
		t.Error(err)
	}
	th := threadStore.Get("Qmabc")
	if th == nil {
		t.Error("could not get thread")
	}
}

func TestThreadDB_List(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&repo.Thread{
		Id:        "Qm123",
		Key:       ksuid.New().String(),
		PrivKey:   make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      repo.PrivateThread,
		Members:   []string{},
		Sharing:   repo.NotSharedThread,
		State:     repo.ThreadLoaded,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadStore.Add(&repo.Thread{
		Id:      "Qm456",
		Key:     ksuid.New().String(),
		PrivKey: make([]byte, 8),
		Name:    "boom",
		Schema:  "Qm...",
		Type:    repo.PrivateThread,
		Members: []string{},
		Sharing: repo.NotSharedThread,
		State:   repo.ThreadLoaded,
	})
	if err != nil {
		t.Error(err)
	}
	all := threadStore.List()
	if len(all) != 2 {
		t.Error("returned incorrect number of threads")
		return
	}
}

func TestThreadDB_Count(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&repo.Thread{
		Id:        "Qm123count",
		Key:       ksuid.New().String(),
		PrivKey:   make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      repo.PrivateThread,
		Members:   []string{},
		Sharing:   repo.NotSharedThread,
		State:     repo.ThreadLoading,
	})
	if err != nil {
		t.Error(err)
	}
	cnt := threadStore.Count()
	if cnt != 1 {
		t.Error("returned incorrect count of threads")
		return
	}
}

func TestThreadDB_UpdateHead(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&repo.Thread{
		Id:        "Qmabc",
		Key:       ksuid.New().String(),
		PrivKey:   make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      repo.PrivateThread,
		Members:   []string{},
		Sharing:   repo.NotSharedThread,
		State:     repo.ThreadLoading,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadStore.UpdateHead("Qmabc", "12345")
	if err != nil {
		t.Error(err)
	}
	th := threadStore.Get("Qmabc")
	if th == nil {
		t.Error("could not get thread")
	}
	if th.Head != "12345" {
		t.Error("update head failed")
	}
}

func TestThreadDB_Delete(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&repo.Thread{
		Id:        "Qm789",
		Key:       ksuid.New().String(),
		PrivKey:   make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      repo.PrivateThread,
		Members:   []string{},
		Sharing:   repo.NotSharedThread,
		State:     repo.ThreadLoaded,
	})
	if err != nil {
		t.Error(err)
	}
	all := threadStore.List()
	if len(all) == 0 {
		t.Error("returned incorrect number of threads")
		return
	}
	err = threadStore.Delete(all[0].Id)
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadStore.PrepareQuery("select id from threads where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow(all[0].Id).Scan(&id)
	if err == nil {
		t.Error("Delete failed")
	}
}
