package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt010(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table contacts (id text primary key not null, address text not null, username text not null, avatar text not null, inboxes blob not null, created integer not null, updated integer not null);
    create index contact_address on contacts (address);
    create index contact_username on contacts (username);
    create index contact_updated on contacts (updated);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into contacts(id, address, username, avatar, inboxes, created, updated) values(?,?,?,?,?,?,?)", "id", "address", "username", "avatar", []byte("foo"), 0, 0)
	if err != nil {
		return err
	}
	return nil
}

func Test011(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt010(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor011
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test new tables
	_, err = db.Exec("insert into peers(id, address, username, avatar, inboxes, created, updated) values(?,?,?,?,?,?,?)", "id2", "address", "username", "avatar", []byte("foo"), 0, 0)
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
	if string(version) != "12" {
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
