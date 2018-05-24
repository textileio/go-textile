package db

import (
	"github.com/textileio/textile-go/core"
	"os"
	"path"
	"testing"
	"time"
)

var testDB *SQLiteDatastore

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
	testDB.config.Configure(time.Now())
}

func teardown() {
	os.RemoveAll(path.Join("./", "datastore"))
}

func TestConfigDB_Create(t *testing.T) {
	if _, err := os.Stat(path.Join("./", "datastore", "mainnet.db")); os.IsNotExist(err) {
		t.Error("Failed to create database file")
	}
}

func TestConfigDB_GetCreationDate(t *testing.T) {
	_, err := testDB.config.GetCreationDate()
	if err != nil {
		t.Error(err)
	}
}

func TestConfigDB_GetVersion(t *testing.T) {
	sv, err := testDB.config.GetVersion()
	if err != nil {
		t.Error(err)
	}
	if sv != core.Version {
		t.Error("version mismatch")
	}
}

func TestConfigDB_IsEncrypted(t *testing.T) {
	encrypted := testDB.Config().IsEncrypted()
	if encrypted {
		t.Error("IsEncrypted returned incorrectly")
	}
}
