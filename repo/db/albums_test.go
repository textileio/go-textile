package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/textileio/textile-go/repo"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

var aldb repo.AlbumStore

func init() {
	setupAlbumDB()
}

func setupAlbumDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	aldb = NewAlbumStore(conn, new(sync.Mutex))
}

func TestAlbumDB_Put(t *testing.T) {
	priv, _, err := libp2p.GenerateKeyPair(libp2p.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	err = aldb.Put(&repo.PhotoAlbum{
		Id:       "Qmabc123",
		Key:      priv,
		Mnemonic: "",
		Name:     "boom",
	})
	if err != nil {
		t.Error(err)
	}
	stmt, err := aldb.PrepareQuery("select id from albums where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("Qmabc123").Scan(&id)
	if err != nil {
		t.Error(err)
	}
	if id != "Qmabc123" {
		t.Errorf(`expected "Qmabc123" got %s`, id)
	}
}

func TestAlbumDB_GetAlbum(t *testing.T) {
	setupAlbumDB()
	priv, _, err := libp2p.GenerateKeyPair(libp2p.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	err = aldb.Put(&repo.PhotoAlbum{
		Id:       "Qmabc",
		Key:      priv,
		Mnemonic: "",
		Name:     "boom",
	})
	if err != nil {
		t.Error(err)
	}
	p := aldb.GetAlbum("Qmabc")
	if p == nil {
		t.Error("could not get album")
	}
}

func TestAlbumDB_GetAlbumByName(t *testing.T) {
	setupAlbumDB()
	priv, _, err := libp2p.GenerateKeyPair(libp2p.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	err = aldb.Put(&repo.PhotoAlbum{
		Id:       "Qmabc",
		Key:      priv,
		Mnemonic: "",
		Name:     "boom",
	})
	if err != nil {
		t.Error(err)
	}
	err = aldb.Put(&repo.PhotoAlbum{
		Id:       "Qmabc2",
		Key:      priv,
		Mnemonic: "",
		Name:     "boom",
	})
	if err == nil {
		t.Error("unique constraint on name failed")
	}
	p := aldb.GetAlbumByName("boom")
	if p == nil {
		t.Error("could not get album")
	}
}

func TestAlbumDB_GetAlbums(t *testing.T) {
	setupAlbumDB()
	priv, _, err := libp2p.GenerateKeyPair(libp2p.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	err = aldb.Put(&repo.PhotoAlbum{
		Id:       "Qm123",
		Key:      priv,
		Mnemonic: "",
		Name:     "boom",
	})
	if err != nil {
		t.Error(err)
	}
	priv, _, err = libp2p.GenerateKeyPair(libp2p.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	err = aldb.Put(&repo.PhotoAlbum{
		Id:       "Qm456",
		Key:      priv,
		Mnemonic: "",
		Name:     "boom2",
	})
	if err != nil {
		t.Error(err)
	}
	as := aldb.GetAlbums("")
	if len(as) != 2 {
		t.Error("returned incorrect number of albums")
		return
	}

	filtered := aldb.GetAlbums("name='boom2'")
	if len(filtered) != 1 {
		t.Error("returned incorrect number of albums")
		return
	}
}

func TestAlbumDB_DeleteAlbum(t *testing.T) {
	setupAlbumDB()
	priv, _, err := libp2p.GenerateKeyPair(libp2p.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	err = aldb.Put(&repo.PhotoAlbum{
		Id:       "Qm789",
		Key:      priv,
		Mnemonic: "",
		Name:     "boom",
	})
	if err != nil {
		t.Error(err)
	}
	as := aldb.GetAlbums("")
	if len(as) == 0 {
		t.Error("returned incorrect number of albums")
		return
	}
	err = aldb.DeleteAlbum(as[0].Id)
	if err != nil {
		t.Error(err)
	}
	stmt, err := phdb.PrepareQuery("select id from albums where id=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow(as[0].Id).Scan(&id)
	if err == nil {
		t.Error("Delete failed")
	}
}
