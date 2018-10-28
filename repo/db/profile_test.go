package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
)

var pdb repo.ProfileStore

func init() {
	setupProfileDB()
}

func setupProfileDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	pdb = NewProfileStore(conn, new(sync.Mutex))
}

func TestProfileDB_SetUsername(t *testing.T) {
	if err := pdb.SetUsername("psyched_mike_79"); err != nil {
		t.Error(err)
		return
	}
}

func TestProfileDB_GetUsername(t *testing.T) {
	un, err := pdb.GetUsername()
	if err != nil {
		t.Error(err)
		return
	}
	if *un != "psyched_mike_79" {
		t.Error("got bad username")
	}
}

func TestProfileDB_SetAvatar(t *testing.T) {
	if err := pdb.SetAvatar("/ipfs/Qm..."); err != nil {
		t.Error(err)
		return
	}
}

func TestProfileDB_GetAvatar(t *testing.T) {
	av, err := pdb.GetAvatar()
	if err != nil {
		t.Error(err)
		return
	}
	if *av != "/ipfs/Qm..." {
		t.Error("got bad avatar")
	}
}
