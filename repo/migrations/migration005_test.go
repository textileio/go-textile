package migrations

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/textileio/go-textile/crypto"
)

func initAt004(db *sql.DB, pin string) error {
	// Only need the Identity stub...
	configStr := `
	{
		"Identity": {
			"PeerID": "QmQA7swSsZKoayPHaTPgzZ1u3SCQjLvLyKcN6RRMmTbLau",
			"PrivKey": "CAESYH1jZmeyepc6aWdAeOkLbkVDYt5FFHIvQramNAGglovRHSxkSGg54g2KJJ/9oqFXJuw2WL009Gap3XnFUxnvKGodLGRIaDniDYokn/2ioVcm7DZYvTT0ZqndecVTGe8oag=="
		}
	}
	`
	if err := ioutil.WriteFile("./config", []byte(configStr), 0644); err != nil {
		return err
	}
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table threads (id text primary key not null, name text not null, sk blob not null, head text not null);
    create table peers (row text primary key not null, id text not null, pk blob not null, threadId text not null);
    create table blocks (id text primary key not null, date integer not null, parents text not null, threadId text not null, authorPk text not null, type integer not null, dataId text, dataKeyCipher blob, dataCaptionCipher blob, dataUsernameCipher blob, dataMetadataCipher blob);
    create table profile (key text primary key not null, value blob);
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
	_, err = db.Exec("insert into profile(key, value) values(?,?)", "username", []byte("username"))
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into peers(row, id, pk, threadId) values(?,?,?,?)", "abc", "Qm123", []byte("foo"), "1")
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into peers(row, id, pk, threadId) values(?,?,?,?)", "def", "Qm456", []byte("bar"), "1")
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
	_ = os.Mkdir("./datastore", os.ModePerm)
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
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	pfile, err := ioutil.ReadFile("./migration005_peerid.ndjson")
	if err != nil {
		t.Error(err)
		return
	}
	var profileInfo map[string]string
	if err := json.Unmarshal(pfile, &profileInfo); err != nil {
		t.Error(err)
		return
	}
	if !strings.HasPrefix(profileInfo["peerid"], "Qm") {
		t.Error(fmt.Errorf("invalid/no peer id saved"))
		return
	}
	if profileInfo["username"] != "username" {
		t.Error(fmt.Errorf("no username saved"))
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
		t.Error(fmt.Errorf("saved wrong number of threads"))
		return
	}
	if len(threads[0].Peers) != 2 {
		t.Error(fmt.Errorf("saved wrong number of thread peers"))
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
		t.Error(fmt.Errorf("saved wrong number of photos"))
		return
	}

	_ = os.RemoveAll("./migration005_threads.ndjson")
	_ = os.RemoveAll("./migration005_peerid.ndjson")
	_ = os.RemoveAll("./migration005_default_photos.ndjson")
	_ = os.RemoveAll("./datastore")
	_ = os.RemoveAll("./repover")
	_ = os.RemoveAll("./config")
}
