package db

import (
	"crypto/rand"
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	ps "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"sync"
	"testing"
	"time"
)

var pndb repo.PointerStore
var pointer repo.Pointer

func init() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	pndb = NewPointerStore(conn, new(sync.Mutex))
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	h, _ := multihash.Encode(randBytes, multihash.SHA2_256)
	id, _ := peer.IDFromBytes(h)
	maAddr, _ := ma.NewMultiaddr("/ipfs/QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr/")
	k, _ := cid.Decode("QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr")
	cancelId, _ := peer.IDB58Decode("QmbwSMS35CaYKdrYBvvR9aHU9FzeWhjJ7E3jLKeR2DWrs3")
	pointer = repo.Pointer{
		Cid: k,
		Value: ps.PeerInfo{
			ID:    id,
			Addrs: []ma.Multiaddr{maAddr},
		},
		Purpose:  repo.MESSAGE,
		Date:     time.Now(),
		CancelId: &cancelId,
	}
}

func TestPointersPut(t *testing.T) {
	err := pndb.Put(pointer)
	if err != nil {
		t.Error(err)
	}

	stmt, _ := pndb.PrepareQuery("select id, key, address, cancelId, purpose, date from pointers where id=?")
	defer stmt.Close()

	var pointerId string
	var key string
	var address string
	var purpose int
	var date int
	var cancelId string
	err = stmt.QueryRow(pointer.Value.ID.Pretty()).Scan(&pointerId, &key, &address, &cancelId, &purpose, &date)
	if err != nil {
		t.Error(err)
	}
	if pointerId != pointer.Value.ID.Pretty() || date <= 0 || key != pointer.Cid.String() || purpose != 1 || cancelId != pointer.CancelId.Pretty() {
		t.Error("pointer returned incorrect values")
	}
	err = pndb.Put(pointer)
	if err == nil {
		t.Error("allowed duplicate pointer")
	}
}

func TestDeletePointer(t *testing.T) {
	pndb.Put(pointer)
	err := pndb.Delete(pointer.Value.ID)
	if err != nil {
		t.Error("pointer delete failed")
	}
	stmt, _ := pndb.PrepareQuery("select id from pointers where id=?")
	defer stmt.Close()

	var pointerId string
	err = stmt.QueryRow(pointer.Value.ID.Pretty()).Scan(&pointerId)
	if err == nil {
		t.Error("pointer delete failed")
	}
}

func TestDeleteAllPointers(t *testing.T) {
	p := pointer
	p.Purpose = repo.MODERATOR
	pndb.Put(p)
	err := pndb.DeleteAll(repo.MODERATOR)
	if err != nil {
		t.Error("pointer delete failed")
	}
	stmt, _ := pndb.PrepareQuery("select id from pointers where purpose=?")
	defer stmt.Close()

	var pointerId string
	err = stmt.QueryRow(repo.MODERATOR).Scan(&pointerId)
	if err == nil {
		t.Error("pointer delete all failed")
	}
}

func TestGetAllPointers(t *testing.T) {
	pndb.Put(pointer)
	pointers, err := pndb.GetAll()
	if err != nil {
		t.Error("get all pointers returned error")
	}
	for _, p := range pointers {
		if p.Purpose != pointer.Purpose {
			t.Error("get all pointers returned incorrect data")
		}
		if p.Value.ID != pointer.Value.ID {
			t.Error("get all pointers returned incorrect data")
		}
		if !p.Cid.Equals(pointer.Cid) {
			t.Error("get all pointers returned incorrect data")
		}
		if p.CancelId.Pretty() != pointer.CancelId.Pretty() {
			t.Error("get all pointers returned incorrect data")
		}
	}
}

func TestPointersDB_GetByPurpose(t *testing.T) {
	pndb.Put(pointer)
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	h, _ := multihash.Encode(randBytes, multihash.SHA2_256)
	id, _ := peer.IDFromBytes(h)
	maAddr, _ := ma.NewMultiaddr("/ipfs/QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr/")
	k, _ := cid.Decode("QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr")
	m := repo.Pointer{
		Cid: k,
		Value: ps.PeerInfo{
			ID:    id,
			Addrs: []ma.Multiaddr{maAddr},
		},
		Purpose:  repo.MODERATOR,
		Date:     time.Now(),
		CancelId: nil,
	}
	err := pndb.Put(m)
	pointers, err := pndb.GetByPurpose(repo.MODERATOR)
	if err != nil {
		t.Error("get pointers returned error")
	}
	if len(pointers) != 1 {
		t.Error("returned incorrect number of pointers")
	}
	for _, p := range pointers {
		if p.Purpose != m.Purpose {
			t.Error("get pointers returned incorrect data")
		}
		if p.Value.ID != m.Value.ID {
			t.Error("get pointers returned incorrect data")
		}
		if !p.Cid.Equals(m.Cid) {
			t.Error("get pointers returned incorrect data")
		}
		if p.CancelId != nil {
			t.Error("get pointers returned incorrect data")
		}
	}
}

func TestPointersDB_Get(t *testing.T) {
	pndb.Put(pointer)
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	h, _ := multihash.Encode(randBytes, multihash.SHA2_256)
	id, _ := peer.IDFromBytes(h)
	maAddr, _ := ma.NewMultiaddr("/ipfs/QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr/")
	k, _ := cid.Decode("QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr")
	m := repo.Pointer{
		Cid: k,
		Value: ps.PeerInfo{
			ID:    id,
			Addrs: []ma.Multiaddr{maAddr},
		},
		Purpose:  repo.MODERATOR,
		Date:     time.Now(),
		CancelId: nil,
	}
	err := pndb.Put(m)
	if err != nil {
		t.Errorf("put pointer returned error: %s", err)
		return
	}
	p := pndb.Get(id)
	if p == nil {
		t.Error("get pointer returned nil")
		return
	}

	if p.Purpose != m.Purpose {
		t.Error("get pointers returned incorrect data")
	}
	if p.Value.ID != m.Value.ID {
		t.Error("get pointers returned incorrect data")
	}
	if !p.Cid.Equals(m.Cid) {
		t.Error("get pointers returned incorrect data")
	}
	if p.CancelId != nil {
		t.Error("get pointers returned incorrect data")
	}
}
