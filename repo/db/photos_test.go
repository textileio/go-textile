package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/textileio/textile-go/repo"
)

var phdb repo.PhotoStore

func init() {
	setupDB()
}

func setupDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	phdb = NewPhotoStore(conn, new(sync.Mutex))
}

func TestPhotoDB_Put(t *testing.T) {
	err := phdb.Put("Qmabc123", time.Now())
	if err != nil {
		t.Error(err)
	}
	stmt, err := phdb.PrepareQuery("select cid, timestamp from photos where cid=?")
	defer stmt.Close()
	var cid string
	var timestamp int
	err = stmt.QueryRow("Qmabc123").Scan(&cid, &timestamp)
	if err != nil {
		t.Error(err)
	}
	if cid != "Qmabc123" {
		t.Errorf(`Expected "Qmabc123" got %s`, cid)
	}
	if timestamp <= 0 {
		t.Error("Returned incorrect timestamp")
	}
}

func TestPhotoDB_GetPhotos(t *testing.T) {
	setupDB()
	err := phdb.Put("Qm123", time.Now())
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 1)
	err = phdb.Put("Qm456", time.Now())
	if err != nil {
		t.Error(err)
	}
	photos := phdb.GetPhotos("", -1)
	if len(photos) != 2 {
		t.Error("Returned incorrect number of photos")
		return
	}

	limted := phdb.GetPhotos("", 1)
	if len(limted) != 1 {
		t.Error("Returned incorrect number of photos")
		return
	}

	offset := phdb.GetPhotos(limted[0].Cid, -1)
	if len(offset) != 1 {
		t.Error("Returned incorrect number of photos")
		return
	}
}

func TestPhotoDB_DeletePhoto(t *testing.T) {
	setupDB()
	err := phdb.Put("Qm789", time.Now())
	if err != nil {
		t.Error(err)
	}
	photos := phdb.GetPhotos("", -1)
	if len(photos) == 0 {
		t.Error("Returned incorrect number of photos")
		return
	}
	err = phdb.DeletePhoto(photos[0].Cid)
	if err != nil {
		t.Error(err)
	}
	stmt, err := phdb.PrepareQuery("select cid from photos where cid=?")
	defer stmt.Close()
	var cid string
	err = stmt.QueryRow(photos[0].Cid).Scan(&cid)
	if err == nil {
		t.Error("Delete failed")
	}
}
