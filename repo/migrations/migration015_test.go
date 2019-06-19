package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt014(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table cafe_requests (id text primary key not null, peerId text not null, targetId text not null, cafeId text not null, cafe blob not null, groupId text not null, syncGroupId text not null, type integer not null, date integer not null, size integer not null, status integer not null, attempts integer not null);
    create index cafe_request_cafeId on cafe_requests (cafeId);
    create index cafe_request_groupId on cafe_requests (groupId);
    create index cafe_request_syncGroupId on cafe_requests (syncGroupId);
    create index cafe_request_date on cafe_requests (date);
    create index cafe_request_status on cafe_requests (status);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into cafe_requests(id, peerId, targetId, cafeId, cafe, groupId, syncGroupId, type, date, size, status, attempts) values(?,?,?,?,?,?,?,?,?,?,?,?)", "id", "peerId", "targetId", "cafeId", []byte("cafe"), "groupId", "syncGroupId", 0, 0, 8, 0, 0)
	if err != nil {
		return err
	}
	return nil
}

func Test015(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt014(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor015
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test new tables
	_, err = db.Exec("insert into cafe_requests(id, peerId, targetId, cafeId, cafe, groupId, syncGroupId, type, date, size, status, attempts, groupSize, groupTransferred) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)", "id2", "peerId", "targetId", "cafeId", []byte("cafe"), "groupId", "syncGroupId", 0, 0, 8, 0, 0, 0, 0)
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
	if string(version) != "16" {
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
