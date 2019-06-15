package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt006(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table thread_invites (id text primary key not null, block blob not null, name text not null, inviter text not null, date integer not null);
    create index thread_invite_date on thread_invites (date);
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into thread_invites(id, block, name, inviter, date) values(?,?,?,?,?)", "test", []byte("block"), "name", "inviter", 0)
	if err != nil {
		return err
	}
	return nil
}

func Test007(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt006(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor007
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test new table
	_, err = db.Exec("insert into thread_invites(id, block, name, contact, date) values(?,?,?,?,?)", "test", []byte("block"), "name", []byte("contact"), 0)
	if err != nil {
		t.Error(err)
		return
	}

	// ensure that version file was updated
	version, err := ioutil.ReadFile("./repover")
	if err != nil {
		t.Error(err)
		return
	}
	if string(version) != "8" {
		t.Error("failed to write new repo version")
		return
	}

	if err := m.Down("./", "", false); err != nil {
		t.Error(err)
		return
	}
	_ = os.RemoveAll("./datastore")
	_ = os.RemoveAll("./repover")
}
