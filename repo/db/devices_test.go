package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
)

var devicedb repo.DeviceStore

func init() {
	setupDeviceDB()
}

func setupDeviceDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	devicedb = NewDeviceStore(conn, new(sync.Mutex))
}

func TestDeviceDB_Add(t *testing.T) {
	err := devicedb.Add(&repo.Device{
		Id:   "abcde",
		Name: "boom",
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := devicedb.PrepareQuery("select id from devices where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != "abcde" {
		t.Errorf(`expected "abcde" got %s`, id)
	}
}

func TestDeviceDB_Get(t *testing.T) {
	block := devicedb.Get("abcde")
	if block == nil {
		t.Error("could not get device")
	}
}

func TestDeviceDB_List(t *testing.T) {
	setupDeviceDB()
	err := devicedb.Add(&repo.Device{
		Id:   "abcde",
		Name: "boom",
	})
	if err != nil {
		t.Error(err)
	}
	err = devicedb.Add(&repo.Device{
		Id:   "abcdef",
		Name: "booom",
	})
	if err != nil {
		t.Error(err)
	}
	all := devicedb.List("")
	if len(all) != 2 {
		t.Error("returned incorrect number of devices")
		return
	}
	filtered := devicedb.List("name='boom'")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of devices")
	}
}

func TestDeviceDB_Count(t *testing.T) {
	setupDeviceDB()
	err := devicedb.Add(&repo.Device{
		Id:   "abcde",
		Name: "hello",
	})
	if err != nil {
		t.Error(err)
	}
	cnt := devicedb.Count("")
	if cnt != 1 {
		t.Error("returned incorrect count of devices")
	}
}

func TestDeviceDB_Delete(t *testing.T) {
	err := devicedb.Delete("abcde")
	if err != nil {
		t.Error(err)
	}
	stmt, err := devicedb.PrepareQuery("select id from devices where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("abcde").Scan(&id)
	if err == nil {
		t.Error("delete failed")
	}
}
