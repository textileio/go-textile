package db

import (
	"crypto/rand"
	"github.com/textileio/textile-go/ipfs"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"os"
	"path"
	"testing"
	"time"
)

var testDB *SQLiteDatastore
var profileId, profileKey string

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func setup() {
	os.MkdirAll(path.Join("./", "datastore"), os.ModePerm)
	testDB, _ = Create("", "LetMeIn")
	testDB.config.Init("LetMeIn")
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}
	profileKey, err = ipfs.EncodeKey(sk)
	if err != nil {
		panic(err)
	}
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		panic(err)
	}
	profileId = id.Pretty()
	testDB.config.Configure(sk, time.Now())
}

func teardown() {
	os.RemoveAll(path.Join("./", "datastore"))
}

func TestConfigDB_Create(t *testing.T) {
	if _, err := os.Stat(path.Join("./", "datastore", "mainnet.db")); os.IsNotExist(err) {
		t.Error("failed to create database file")
	}
}

func TestConfigDB_GetId(t *testing.T) {
	id, err := testDB.config.GetId()
	if err != nil {
		t.Error(err)
		return
	}
	if id == nil {
		t.Error("missing id")
		return
	}
	if *id != profileId {
		t.Error("got bad id")
	}
}

func TestConfigDB_GetKey(t *testing.T) {
	key, err := testDB.config.GetKey()
	if err != nil {
		t.Error(err)
		return
	}
	if key == nil {
		t.Error("missing key")
		return
	}
	keystr, err := ipfs.EncodeKey(key)
	if err != nil {
		t.Error(err)
		return
	}
	if keystr != profileKey {
		t.Error("got bad key")
	}
}

func TestConfigDB_GetCreationDate(t *testing.T) {
	_, err := testDB.config.GetCreationDate()
	if err != nil {
		t.Error(err)
	}
}

func TestConfigDB_IsEncrypted(t *testing.T) {
	encrypted := testDB.Config().IsEncrypted()
	if encrypted {
		t.Error("IsEncrypted returned incorrectly")
	}
}
