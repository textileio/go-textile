package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt008(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table cafe_clients (id text primary key not null, address text not null, created integer not null, lastSeen integer not null);
`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into cafe_clients(id, address, created, lastSeen) values(?,?,?,?)", "test", "address", 0, 0)
	if err != nil {
		return err
	}
	return nil
}

func Test009(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
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

	// test new field
	_, err = db.Exec("update cafe_clients set tokenId=? where id=?", "token", "test")
	if err != nil {
		t.Error(err)
		return
	}

	// test new table
	_, err = db.Exec("insert into cafe_tokens(id, token, date) values(?,?,?)", "id", []byte("token"), 0)
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
	_ = os.RemoveAll("./datastore")
	_ = os.RemoveAll("./repover")
}
