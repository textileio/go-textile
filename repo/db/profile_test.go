package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/textileio/textile-go/repo"
)

var profileStore repo.ProfileStore

func init() {
	setupProfileDB()
}

func setupProfileDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	profileStore = NewProfileStore(conn, new(sync.Mutex))
}

func TestProfileDB_SetUsername(t *testing.T) {
	if err := profileStore.SetUsername("psyched_mike_79"); err != nil {
		t.Error(err)
		return
	}
}

func TestProfileDB_GetUsername(t *testing.T) {
	un, err := profileStore.GetUsername()
	if err != nil {
		t.Error(err)
		return
	}
	if *un != "psyched_mike_79" {
		t.Error("got bad username")
	}
}

func TestProfileDB_SetAvatar(t *testing.T) {
	if err := profileStore.SetAvatar("/ipfs/Qm..."); err != nil {
		t.Error(err)
		return
	}
}

func TestProfileDB_GetAvatar(t *testing.T) {
	av, err := profileStore.GetAvatar()
	if err != nil {
		t.Error(err)
		return
	}
	if *av != "/ipfs/Qm..." {
		t.Error("got bad avatar")
	}
}
