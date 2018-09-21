package db

import (
	"crypto/rand"
	"database/sql"
	"github.com/textileio/textile-go/photo"
	"github.com/textileio/textile-go/repo"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"sync"
	"testing"
	"time"
)

var pdb repo.ProfileStore

func init() {
	setupProfileDB()
}

func setupProfileDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	pdb = NewProfileStore(conn, new(sync.Mutex))
}

func TestProfileDB_GetCafeTokensPreLogin(t *testing.T) {
	tokens, err := pdb.GetCafeTokens()
	if err != nil {
		t.Error(err)
		return
	}
	if tokens != nil {
		t.Error("tokens should be nil")
	}
}

func TestProfileDB_CafeLogin(t *testing.T) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	profileKey, err = photo.EncodeKey(sk)
	if err != nil {
		t.Error(err)
	}
	exp := time.Now().Add(time.Hour)
	if err := pdb.CafeLogin(&repo.CafeTokens{Access: "access", Refresh: "refresh", Expiry: exp}); err != nil {
		t.Error(err)
	}
}

func TestProfileDB_SetUsername(t *testing.T) {
	if err := pdb.SetUsername("psyched_mike_79"); err != nil {
		t.Error(err)
		return
	}
}

func TestProfileDB_GetUsername(t *testing.T) {
	un, err := pdb.GetUsername()
	if err != nil {
		t.Error(err)
		return
	}
	if *un != "psyched_mike_79" {
		t.Error("got bad username")
	}
}

func TestProfileDB_SetAvatarId(t *testing.T) {
	if err := pdb.SetAvatarId("/ipfs/Qm..."); err != nil {
		t.Error(err)
		return
	}
}

func TestProfileDB_GetAvatarId(t *testing.T) {
	av, err := pdb.GetAvatarId()
	if err != nil {
		t.Error(err)
		return
	}
	if *av != "/ipfs/Qm..." {
		t.Error("got bad avatar id")
	}
}

func TestProfileDB_GetCafeTokens(t *testing.T) {
	tokens, err := pdb.GetCafeTokens()
	if err != nil {
		t.Error(err)
		return
	}
	if tokens.Access != "access" {
		t.Error("got bad access token")
		return
	}
	if tokens.Refresh != "refresh" {
		t.Error("got bad refresh token")
		return
	}
	if tokens.Expiry.Before(time.Now()) {
		t.Error("got bad expiry")
	}
}

func TestProfileDB_CafeLogout(t *testing.T) {
	if err := pdb.CafeLogout(); err != nil {
		t.Error(err)
		return
	}
	tokens, err := pdb.GetCafeTokens()
	if err != nil {
		t.Error(err)
	}
	if tokens != nil {
		t.Error("logged out but tokens still present")
	}
}
