package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func initAt015(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func Test016(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt015(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor016
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test new tables
	_, err = db.Exec("insert into botstore(id, key, value, version, created, updated) values(?,?,?,?,?,?)", "botid", "valuekey", []byte("some value"), 1, time.Now().UnixNano(), time.Now().UnixNano())
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
	if string(version) != "17" {
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
