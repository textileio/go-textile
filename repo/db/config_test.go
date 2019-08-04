package db

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/textileio/go-textile/wallet"
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
	_ = os.RemoveAll(path.Join("./", "datastore"))
	_ = os.MkdirAll(path.Join("./", "datastore"), os.ModePerm)
	testDB, _ = Create("", "letmein")
	_ = testDB.config.Init("letmein")

	w, err := wallet.WalletFromEntropy(128)
	if err != nil {
		panic(err)
	}
	a0, err := w.AccountAt(0, "letmeout")
	if err != nil {
		panic(err)
	}
	configAddress = a0.Address()
	if err := testDB.config.Configure(a0, time.Now()); err != nil {
		panic(err)
	}
}

func teardown() {
	_ = os.RemoveAll(path.Join("./", "datastore"))
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

func TestConfigDB_GetLastDaily(t *testing.T) {
	date, err := testDB.config.GetLastDaily()
	if err != nil {
		t.Error(err)
	}
	if date.Unix() > 0 {
		t.Error("last daily date should initially be less than 0")
	}
}

func TestConfigDB_SetLastDaily(t *testing.T) {
	if err := testDB.config.SetLastDaily(); err != nil {
		t.Error(err)
	}
}
func TestConfigDB_GetLastDailyAgain(t *testing.T) {
	date, err := testDB.config.GetLastDaily()
	if err != nil {
		t.Error(err)
	}
	if date.Unix() <= 0 {
		t.Error("last daily date should now be greater than 0")
	}
}
