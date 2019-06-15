package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt013(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table blocks (id text primary key not null, threadId text not null, authorId text not null, type integer not null, date integer not null, parents text not null, target text not null, body text not null);
    create index block_threadId on blocks (threadId);
    create index block_type on blocks (type);
    create index block_date on blocks (date);
    create index block_target on blocks (target);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into blocks(id, threadId, authorId, type, date, parents, target, body) values(?,?,?,?,?,?,?,?)", "id", "threadId", "authorId", 7, 0, "parents", "target", "body")
	if err != nil {
		return err
	}
	return nil
}

func Test014(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt013(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor014
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test target -> data
	rows, err := db.Query("select target, data from blocks where type=7;")
	if err != nil {
		t.Error(err)
		return
	}
	var count int
	for rows.Next() {
		var target, data string
		err := rows.Scan(&target, &data)
		if err != nil {
			t.Error(err)
			return
		}
		if target != "" {
			t.Errorf("expected empty target, got %s", target)
		}
		if data != "target" {
			t.Errorf("expected data to be target, got %s", data)
		}
		count++
	}
	if count != 1 {
		t.Errorf("expected two rows, got %d", count)
	}

	// test new tables
	_, err = db.Exec("insert into blocks(id, threadId, authorId, type, date, parents, target, body, data, status, attempts) values(?,?,?,?,?,?,?,?,?,?,?)", "id2", "threadId", "authorId", 7, 0, "parents", "target", "body", "data", 0, 0)
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
	if string(version) != "15" {
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
