package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt003(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
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

func Test004(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
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
	var m Minor004
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test new table
	_, err = db.Exec("insert into notifications(id, date, actorId, actorUsername, subject, subjectId, blockId, dataId, type, body, read) values(?,?,?,?,?,?,?,?,?,?,?)", "test", 0, "actorId", "james", "cats", "catId", "", "", 0, "hey!", 0)
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
	_ = os.RemoveAll("./datastore")
	_ = os.RemoveAll("./repover")
}
