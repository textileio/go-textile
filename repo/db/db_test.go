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

func TestConfigDB_SignIn(t *testing.T) {
	err := testDB.config.SignIn("woohoo!", "...", "...")
	if err != nil {
		t.Error(err)
	}
}

func TestConfigDB_GetUsername(t *testing.T) {
	un, err := testDB.config.GetUsername()
	if err != nil {
		t.Error(err)
		return
	}
	if un != "woohoo!" {
		t.Error("got bad username")
	}
}

func TestConfigDB_GetTokens(t *testing.T) {
	at, rt, err := testDB.config.GetTokens()
	if err != nil {
		t.Error(err)
		return
	}
	if at != "..." {
		t.Error("got bad access token")
		return
	}
	if rt != "..." {
		t.Error("got bad refresh token")
		return
	}
}

func TestConfigDB_SignOut(t *testing.T) {
	err := testDB.config.SignOut()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = testDB.config.GetUsername()
	if err == nil {
		t.Error("signed out but username still present")
	}
}

func TestConfigDB_GetCreationDate(t *testing.T) {
	_, err := testDB.config.GetCreationDate()
	if err != nil {
		t.Error(err)
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

func TestConfigDB_IsEncrypted(t *testing.T) {
	encrypted := testDB.Config().IsEncrypted()
	if encrypted {
		t.Error("IsEncrypted returned incorrectly")
	}
}
