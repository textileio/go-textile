package db

import (
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
	testDB.config.Configure("Mnemonic Passphrase", []byte("Private Key"), time.Now())
}

func teardown() {
	os.RemoveAll(path.Join("./", "datastore"))
}

func TestCreate(t *testing.T) {
	if _, err := os.Stat(path.Join("./", "datastore", "mainnet.db")); os.IsNotExist(err) {
		t.Error("Failed to create database file")
	}
}

func TestConfig(t *testing.T) {
	mn, err := testDB.config.GetMnemonic()
	if err != nil {
		t.Error(err)
	}
	if mn != "Mnemonic Passphrase" {
		t.Error("Config returned wrong mnemonic")
	}
	pk, err := testDB.config.GetIdentityKey()
	if err != nil {
		t.Error(err)
	}
	testKey := []byte("Private Key")
	for i := range pk {
		if pk[i] != testKey[i] {
			t.Error("Config returned wrong identity key")
		}
	}
}

func TestInterface(t *testing.T) {
	if testDB.Config() != testDB.config {
		t.Error("Config() return wrong value")
	}
	if testDB.Settings() != testDB.settings {
		t.Error("Settings() return wrong value")
	}
}

func TestEncryptedDb(t *testing.T) {
	encrypted := testDB.Config().IsEncrypted()
	if encrypted {
		t.Error("IsEncrypted returned incorrectly")
	}
}
