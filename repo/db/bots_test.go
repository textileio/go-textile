package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

var botStore repo.Botstore

var testBot *pb.BotKV

func init() {
	setupBotDB()
}

func setupBotDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	_ = initDatabaseTables(conn, "")
	botStore = NewBotstore(conn, new(sync.Mutex))
}

func TestBotDB_Add(t *testing.T) {
	newValue := "24242322"
	err := botStore.AddOrUpdate("ABCDEFG", []byte(newValue))
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := botStore.PrepareQuery("select value from bots_store where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var value []byte
	err = stmt.QueryRow("ABCDEFG").Scan(&value)
	if err != nil {
		t.Error(err)
		return
	}
	if string(value) != newValue {
		t.Errorf(`expected "24242322" got %s`, string(newValue))
	}
}

func TestBotDB_AddOrUpdate(t *testing.T) {
	valueKey := "ABCDEFG"
	newValue := "{count:2}"
	err := botStore.AddOrUpdate(valueKey, []byte(newValue))
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := botStore.PrepareQuery("select value, updated from bots_store where id=?")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	var value []byte
	var updated int64
	err = stmt.QueryRow(valueKey).Scan(&value, &updated)
	if err != nil {
		t.Error(err)
		return
	}
	if string(value) != newValue {
		t.Errorf(`expected "{count:2}" got %s`, string(value))
		return
	}
}

func TestBotDB_Get(t *testing.T) {
	testBotValue := botStore.Get("ABCDEFG")
	if testBotValue == nil {
		t.Error("could not get bot")
	}
}

func TestBotDB_Delete(t *testing.T) {
	err := botStore.Delete("ABCDEFG")
	if err != nil {
		t.Error(err)
	}
	stmt, err := botStore.PrepareQuery("select id from bots_store where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("ABCDEFG").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
