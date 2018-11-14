package migrations

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/textileio/textile-go/crypto"
)

func initAt004(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
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

	// check threads
	tfile, err := ioutil.ReadFile("./migration005_threads.ndjson")
	if err != nil {
		t.Error(err)
		return
	}
	var threads []*threadRow
	threadRows := strings.Split(string(tfile), "\n")
	for _, row := range threadRows {
		if len(row) == 0 {
			continue
		}
		thrd := new(threadRow)
		if err := json.Unmarshal([]byte(row), &thrd); err != nil {
			t.Error(err)
			return
		}
		threads = append(threads, thrd)
	}
	if len(threads) != 1 {
		t.Error(errors.New("saved wrong number of threads"))
		return
	}

	// check photos
	ffile, err := ioutil.ReadFile("./migration005_default_photos.ndjson")
	if err != nil {
		t.Error(err)
		return
	}
	var photos []*photoRow
	photoRows := strings.Split(string(ffile), "\n")
	for _, row := range photoRows {
		if len(row) == 0 {
			continue
		}
		photo := new(photoRow)
		if err := json.Unmarshal([]byte(row), &photo); err != nil {
			t.Error(err)
			return
		}
		photos = append(photos, photo)
	}
	if len(photos) != 2 {
		t.Error(errors.New("saved wrong number of photos"))
		return
	}

	os.RemoveAll("./migration005_threads.ndjson")
	os.RemoveAll("./migration005_default_photos.ndjson")
	os.RemoveAll("./datastore")
	os.RemoveAll("./repover")
}
