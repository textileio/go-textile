package db

import (
	"database/sql"
	"encoding/json"
	"sync"
	"testing"

	"github.com/textileio/textile-go/repo"
)

var sdb repo.ConfigurationStore
var settings repo.SettingsData

func init() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	sdb = NewConfigurationStore(conn, new(sync.Mutex))
	v := "/textile-go:0.1/"
	settings = repo.SettingsData{
		Version: &v,
	}
}

func TestSettingsPut(t *testing.T) {
	err := sdb.Put(settings)
	if err != nil {
		t.Error(err)
	}
	set := repo.SettingsData{}
	stmt, err := sdb.PrepareQuery("select value from config where key=?")
	defer stmt.Close()
	var settingsBytes []byte
	err = stmt.QueryRow("settings").Scan(&settingsBytes)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(settingsBytes, &set)
	if err != nil {
		t.Error(err)
	}
	if *set.Version!= "/textile-go:0.1/" {
		t.Error("Settings put failed to put correct value")
	}
}

func TestInvalidSettingsGet(t *testing.T) {
	tx, err := sdb.BeginTransaction()
	if err != nil {
		t.Error(err)
	}
	stmt, _ := tx.Prepare("insert or replace into config(key, value) values(?,?)")
	defer stmt.Close()

	_, err = stmt.Exec("settings", string("Test fail"))
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}
	tx.Commit()
	_, err = sdb.Get()
	if err == nil {
		t.Error("settings get didn't error with invalid data")
	}
}

func TestSettingsGet(t *testing.T) {
	err := sdb.Put(settings)
	if err != nil {
		t.Error(err)
	}
	set, err := sdb.Get()
	if err != nil {
		t.Error(err)
	}
	if *set.Version != "/textile-go:0.1/" {
		t.Error("Settings put failed to put correct value")
	}
}

func TestSettingsUpdate(t *testing.T) {
	err := sdb.Put(settings)
	if err != nil {
		t.Error(err)
	}
	l := "/textile-go:0.2/"
	setUpdt := repo.SettingsData{
		Version: &l,
	}
	err = sdb.Update(setUpdt)
	if err != nil {
		t.Error(err)
	}
	set, err := sdb.Get()
	if err != nil {
		t.Error(err)
	}
	if *set.Version != "/textile-go:0.2/" {
		t.Error("Settings update failed to put correct value")
	}
}
