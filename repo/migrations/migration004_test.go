package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt003(db *sql.DB, password string) error {
	var sqlStmt string
	if password != "" {
		sqlStmt = "PRAGMA key = '" + password + "';"
	}
	sqlStmt += `
    create table notifications (id text primary key not null, date integer not null, actorId text not null, targetId text not null, type integer not null, read integer not null, body text not null, actorUn text not null, category text not null);
    create index notification_targetId on notifications (targetId);
    create index notification_actorId on notifications (actorId);
    create index notification_read on notifications (read);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into notifications(id, date, actorId, targetId, type, read, body, actorUn, category) values(?,?,?,?,?,?,?,?,?)", "test", 0, "actorId", "targetId", 0, 0, "hey!", "bob", "cats")
	if err != nil {
		return err
	}
	return nil
}

func TestMigration004(t *testing.T) {
	var dbPath string
	os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt003(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Migration004
	err = m.Up("./", "", false)
	if err != nil {
		t.Error(err)
		return
	}

	// test new fields
	_, err = db.Exec("update notifications set categoryId=? where id=?", "catId", "test")
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
	if string(version) != "5" {
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
