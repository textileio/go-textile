package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt008(db *sql.DB, pin string) error {
	// nothing to do here
	return nil
}

func Test009(t *testing.T) {
	var dbPath string
	os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt008(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor009
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test new table
	_, err = db.Exec("insert into cafe_dev_tokens(id, token, created) values(?,?,?)", "id", "token", 0)
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
	if string(version) != "10" {
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
