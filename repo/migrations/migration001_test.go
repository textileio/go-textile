package migrations

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt000(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table blocks (id text primary key not null, date integer not null, parents text not null, threadId text not null, authorPk text not null, type integer not null, dataId text, dataKeyCipher blob, dataCaptionCipher blob, dataUsernameCipher blob);
    create index block_dataId on blocks (dataId);
    create index block_threadId_type_date on blocks (threadId, type, date);
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into blocks(id, date, parents, threadId, authorPk, type, dataId, dataKeyCipher, dataCaptionCipher, dataUsernameCipher) values(?,?,?,?,?,?,?,?,?,?)", "test", 0, "parents", "threadId", "authorPk", 4, "dataId", []byte("dataKeyCipher"), []byte("dataCaptionCipher"), []byte("dataUsernameCipher"))
	if err != nil {
		return err
	}
	return nil
}

func Test001(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt000(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor001
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// test new field
	_, err = db.Exec("update blocks set dataMetadataCipher=? where id=?", []byte("dataMetadataCipher"), "test")
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
	if string(version) != "2" {
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
