package migrations

import (
	"crypto/rand"
	"database/sql"
	"github.com/textileio/textile-go/crypto"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func initAt004(db *sql.DB, password string) error {
	var sqlStmt string
	if password != "" {
		sqlStmt = "PRAGMA key = '" + password + "';"
	}
	sqlStmt += `
    create table threads (id text primary key not null, name text not null, sk blob not null, head text not null);
    create table blocks (id text primary key not null, date integer not null, parents text not null, threadId text not null, authorPk text not null, type integer not null, dataId text, dataKeyCipher blob, dataCaptionCipher blob, dataUsernameCipher blob, dataMetadataCipher blob);
    create index block_dataId on blocks (dataId);
    create index block_threadId_type_date on blocks (threadId, type, date);
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return err
	}
	skb, err := sk.Bytes()
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into threads(id, name, sk, head) values(?,?,?,?)", "1", "default", skb, "")
	if err != nil {
		return err
	}
	keyc1, err := crypto.Encrypt(sk.GetPublic(), []byte("imakey"))
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into blocks(id, date, parents, threadId, authorPk, type, dataId, dataKeyCipher, dataCaptionCipher, dataUsernameCipher, dataMetadataCipher) values(?,?,?,?,?,?,?,?,?,?,?)", "1", 0, "", "1", "", 4, "Qmtester1", keyc1, []byte("x"), []byte("x"), []byte("x"))
	if err != nil {
		return err
	}
	keyc2, err := crypto.Encrypt(sk.GetPublic(), []byte("imakey2"))
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into blocks(id, date, parents, threadId, authorPk, type, dataId, dataKeyCipher, dataCaptionCipher, dataUsernameCipher, dataMetadataCipher) values(?,?,?,?,?,?,?,?,?,?,?)", "2", 0, "", "1", "", 4, "Qmtester2", keyc2, []byte("x"), []byte("x"), []byte("x"))
	if err != nil {
		return err
	}
	return nil
}

func Test005(t *testing.T) {
	var dbPath string
	os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt004(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Major005
	err = m.Up("./", "", false)
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
	if string(version) != "6" {
		t.Error("failed to write new repo version")
		return
	}

	os.RemoveAll("./migration005_threads.ndjson")
	os.RemoveAll("./migration005_default_photos.ndjson")
	os.RemoveAll("./datastore")
	os.RemoveAll("./repover")
}
