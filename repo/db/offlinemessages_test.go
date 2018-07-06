package db

import (
	"bytes"
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"testing"
)

var odb repo.OfflineMessageStore

func init() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	odb = NewOfflineMessageStore(conn, new(sync.Mutex))
}

func TestOfflineMessagesPut(t *testing.T) {
	err := odb.Put("abc")
	if err != nil {
		t.Error(err)
	}

	stmt, _ := odb.PrepareQuery("select url, date from offlinemessages where url=?")
	defer stmt.Close()

	var url string
	var date int
	err = stmt.QueryRow("abc").Scan(&url, &date)
	if err != nil {
		t.Error(err)
	}
	if url != "abc" || date <= 0 {
		t.Error("fffline messages put failed")
	}
}

func TestOfflineMessagesPutDuplicate(t *testing.T) {
	err := odb.Put("123")
	if err != nil {
		t.Error(err)
	}
	err = odb.Put("123")
	if err == nil {
		t.Error("put offline messages duplicate returned no error")
	}
}

func TestOfflineMessagesHas(t *testing.T) {
	err := odb.Put("abcc")
	if err != nil {
		t.Error(err)
	}
	has := odb.Has("abcc")
	if !has {
		t.Error("failed to find offline message url in db")
	}
	has = odb.Has("xyz")
	if has {
		t.Error("offline messages has returned incorrect")
	}
}

func TestOfflineMessagesSetMessage(t *testing.T) {
	err := odb.Put("abccc")
	if err != nil {
		t.Error(err)
	}
	err = odb.SetMessage("abccc", []byte("helloworld"))
	if err != nil {
		t.Error(err)
	}
	messages, err := odb.GetMessages()
	if err != nil {
		t.Error(err)
	}
	m, ok := messages["abccc"]
	if !ok || !bytes.Equal(m, []byte("helloworld")) {
		t.Error("returned incorrect value")
	}

	err = odb.DeleteMessage("abccc")
	if err != nil {
		t.Error(err)
	}
	messages, err = odb.GetMessages()
	if err != nil {
		t.Error(err)
	}
	m, ok = messages["abccc"]
	if ok {
		t.Error("failed to delete")
	}
}
