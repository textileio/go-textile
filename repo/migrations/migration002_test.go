package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestMigration002(t *testing.T) {
	var dbPath string
	os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Migration002
	err = m.Up("./", "", false)
	if err != nil {
		t.Error(err)
		return
	}

	// test new table
	_, err = db.Exec("insert into notifications(id, date, actorId, targetId, type, read) values(?,?,?,?,?,?,?)", "test", 0, "actorId", "targetId", 0, 0, "hey!")
	if err != nil {
		t.Error(err)
	}

	// ensure that version file was updated
	version, err := ioutil.ReadFile("./repover")
	if err != nil {
		t.Error(err)
		return
	}
	if string(version) != "3" {
		t.Error("failed to write new repo version")
		return
	}

	if err := m.Down("./", "", false); err != nil {
		t.Error(err)
		return
	}
	os.RemoveAll("./datastore")
	os.RemoveAll("./repover")
}
