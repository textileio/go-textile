package db

import (
	"github.com/textileio/textile-go/wallet"
	"os"
	"path"
	"testing"
	"time"
)

var testDB *SQLiteDatastore
var configAddress string

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func setup() {
	os.RemoveAll(path.Join("./", "datastore"))
	os.MkdirAll(path.Join("./", "datastore"), os.ModePerm)
	testDB, _ = Create("", "letmein")
	testDB.config.Init("letmein")

	w, err := wallet.NewWallet(128)
	if err != nil {
		panic(err)
	}
	a0, err := w.AccountAt(0, "letmeout")
	if err != nil {
		panic(err)
	}
	configAddress = a0.Address()
	if err := testDB.config.Configure(a0, true, time.Now()); err != nil {
		panic(err)
	}
}

func teardown() {
	os.RemoveAll(path.Join("./", "datastore"))
}

func TestConfigDB_Create(t *testing.T) {
	if _, err := os.Stat(path.Join("./", "datastore", "mainnet.db")); os.IsNotExist(err) {
		t.Error("failed to create database file")
	}
}

func TestConfigDB_GetAccount(t *testing.T) {
	account, err := testDB.config.GetAccount()
	if err != nil {
		t.Error(err)
		return
	}
	if account == nil {
		t.Error("missing account")
		return
	}
	if account.Address() != configAddress {
		t.Error("got bad account")
	}
}

func TestConfigDB_GetMobile(t *testing.T) {
	mobile, err := testDB.config.GetMobile()
	if err != nil {
		t.Error(err)
	}
	if !mobile {
		t.Error("mobile should be true")
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
