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

func TestProfileDB_Init(t *testing.T) {
	err := pdb.Init("boom", []byte("..."))
	if err != nil {
		t.Error(err)
	}
}

func TestProfileDB_SignIn(t *testing.T) {
	err := pdb.SignIn("woohoo!", "...", "...")
	if err != nil {
		t.Error(err)
	}
}

func TestProfileDB_GetId(t *testing.T) {
	id, err := pdb.GetId()
	if err != nil {
		t.Error(err)
		return
	}
	if id != "boom" {
		t.Error("got bad id")
	}
}

func TestProfileDB_GetSecret(t *testing.T) {
	secret, err := pdb.GetSecret()
	if err != nil {
		t.Error(err)
		return
	}
	if string(secret) != "..." {
		t.Error("got bad secret")
	}
}

func TestProfileDB_GetUsername(t *testing.T) {
	un, err := pdb.GetUsername()
	if err != nil {
		t.Error(err)
		return
	}
	if un != "woohoo!" {
		t.Error("got bad username")
	}
}

func TestProfileDB_GetTokens(t *testing.T) {
	at, rt, err := pdb.GetTokens()
	if err != nil {
		t.Error(err)
		return
	}
	if at != "..." {
		t.Error("got bad access token")
		return
	}
	if rt != "..." {
		t.Error("got bad refresh token")
		return
	}
}

func TestProfileDB_SignOut(t *testing.T) {
	err := pdb.SignOut()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = pdb.GetUsername()
	if err == nil {
		t.Error("signed out but username still present")
	}
}
