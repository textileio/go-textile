package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
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
	err := threadStore.Add(&pb.Thread{
		Id:        "Qmabc123",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_Open,
		Members:   []string{"P1,P2"},
		Sharing:   pb.Thread_Shared,
		State:     pb.Thread_Loaded,
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
	err := threadStore.Add(&pb.Thread{
		Id:        "Qmabc",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_Open,
		Members:   []string{},
		Sharing:   pb.Thread_Shared,
		State:     pb.Thread_Loaded,
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
	err := threadStore.Add(&pb.Thread{
		Id:        "Qm123",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_Private,
		Members:   []string{},
		Sharing:   pb.Thread_NotShared,
		State:     pb.Thread_Loaded,
	})
	if err != nil {
		t.Error(err)
	}
	err = threadStore.Add(&pb.Thread{
		Id:      "Qm456",
		Key:     ksuid.New().String(),
		Sk:      make([]byte, 8),
		Name:    "boom",
		Schema:  "Qm...",
		Type:    pb.Thread_Private,
		Members: []string{},
		Sharing: pb.Thread_NotShared,
		State:   pb.Thread_Loaded,
	})
	if err != nil {
		t.Error(err)
	}
	all := threadStore.List()
	if len(all.Items) != 2 {
		t.Error("returned incorrect number of threads")
		return
	}
}

func TestThreadDB_Count(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&pb.Thread{
		Id:        "Qm123count",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_Private,
		Members:   []string{},
		Sharing:   pb.Thread_NotShared,
		State:     pb.Thread_Loaded,
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
	err := threadStore.Add(&pb.Thread{
		Id:        "Qmabc",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_Private,
		Members:   []string{},
		Sharing:   pb.Thread_NotShared,
		State:     pb.Thread_Loaded,
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
	err := threadStore.Add(&pb.Thread{
		Id:        "Qm789",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_Private,
		Members:   []string{},
		Sharing:   pb.Thread_NotShared,
		State:     pb.Thread_Loaded,
	})
	if err != nil {
		t.Error(err)
	}
	all := threadStore.List()
	if len(all.Items) == 0 {
		t.Error("returned incorrect number of threads")
		return
	}
	err = threadStore.Delete(all.Items[0].Id)
	if err != nil {
		t.Error(err)
	}
	stmt, err := threadStore.PrepareQuery("select id from threads where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow(all.Items[0].Id).Scan(&id)
	if err == nil {
		t.Error("Delete failed")
	}
}
