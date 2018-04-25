package db

import (
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/photos"
)

var phdb repo.PhotoStore

func init() {
	setupPhotoDB()
}

func setupPhotoDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	phdb = NewPhotoStore(conn, new(sync.Mutex))
}

func TestPhotoDB_Put(t *testing.T) {
	err := phdb.Put(&repo.PhotoSet{
		Cid:      "Qmabc123",
		LastCid:  "",
		AlbumID:  "Qm",
		MetaData: photos.Metadata{},
		IsLocal:  false,
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := phdb.PrepareQuery("select cid from photos where cid=?")
	defer stmt.Close()
	var cid string
	err = stmt.QueryRow("Qmabc123").Scan(&cid)
	if err != nil {
		t.Error(err)
	}
	if cid != "Qmabc123" {
		t.Errorf(`expected "Qmabc123" got %s`, cid)
	}
}

func TestPhotoDB_GetPhoto(t *testing.T) {
	setupPhotoDB()
	err := phdb.Put(&repo.PhotoSet{
		Cid:     "Qmabc",
		LastCid: "",
		AlbumID: "Qm",
		MetaData: photos.Metadata{
			Added: time.Now(),
		},
		IsLocal: true,
	})
	if err != nil {
		t.Error(err)
	}
	p := phdb.GetPhoto("Qmabc")
	if p == nil {
		t.Error("could not get photo")
		return
	}
}

func TestPhotoDB_GetPhotos(t *testing.T) {
	setupPhotoDB()
	err := phdb.Put(&repo.PhotoSet{
		Cid:     "Qm123",
		LastCid: "",
		AlbumID: "Qm",
		MetaData: photos.Metadata{
			Added: time.Now(),
		},
		IsLocal: true,
	})
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 1)
	err = phdb.Put(&repo.PhotoSet{
		Cid:     "Qm456",
		LastCid: "Qm123",
		AlbumID: "Qm",
		MetaData: photos.Metadata{
			Added: time.Now(),
		},
		IsLocal: false,
	})
	if err != nil {
		t.Error(err)
	}
	ps := phdb.GetPhotos("", -1, "")
	if len(ps) != 2 {
		t.Error("returned incorrect number of photos")
		return
	}

	limited := phdb.GetPhotos("", 1, "")
	if len(limited) != 1 {
		t.Error("returned incorrect number of photos")
		return
	}

	offset := phdb.GetPhotos(limited[0].Cid, -1, "")
	if len(offset) != 1 {
		t.Error("returned incorrect number of photos")
		return
	}

	filtered := phdb.GetPhotos("", -1, "local=1")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of photos")
		return
	}
}

func TestPhotoDB_DeletePhoto(t *testing.T) {
	setupPhotoDB()
	err := phdb.Put(&repo.PhotoSet{
		Cid:      "Qm789",
		LastCid:  "",
		AlbumID:  "Qm",
		MetaData: photos.Metadata{},
		IsLocal:  true,
	})
	if err != nil {
		t.Error(err)
	}
	ps := phdb.GetPhotos("", -1, "")
	if len(ps) == 0 {
		t.Error("returned incorrect number of photos")
		return
	}
	err = phdb.DeletePhoto(ps[0].Cid)
	if err != nil {
		t.Error(err)
	}
	stmt, err := phdb.PrepareQuery("select cid from photos where cid=?")
	defer stmt.Close()
	var cid string
	err = stmt.QueryRow(ps[0].Cid).Scan(&cid)
	if err == nil {
		t.Error("Delete failed")
	}
}
