package db

import (
	"crypto/rand"
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"sync"
	"testing"
	"time"
)

var pdb repo.ProfileStore

var profileKey string

func init() {
	setupProfileDB()
}

func setupProfileDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	pdb = NewProfileStore(conn, new(sync.Mutex))
}

func TestProfileDB_GetTokensPreLogin(t *testing.T) {
	tokens, err := pdb.GetTokens()
	if err != nil {
		t.Error(err)
		return
	}
	if tokens != nil {
		t.Error("tokens should be nil")
	}
}

func TestProfileDB_Login(t *testing.T) {
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	profileKey, err = util.EncodeKey(sk)
	if err != nil {
		t.Error(err)
	}
	if err := pdb.Login(sk, &repo.CafeTokens{Access: "access", Refresh: "refresh", Expiry: time.Now()}); err != nil {
		t.Error(err)
	}
}

func TestProfileDB_GetKey(t *testing.T) {
	key, err := pdb.GetKey()
	if err != nil {
		t.Error(err)
		return
	}
	if key == nil {
		t.Error("missing key")
		return
	}
	keystr, err := util.EncodeKey(key)
	if err != nil {
		t.Error(err)
		return
	}
	if keystr != profileKey {
		t.Error("got bad key")
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

func TestProfileDB_GetTokens(t *testing.T) {
	tokens, err := pdb.GetTokens()
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
}

func TestProfileDB_UpdateTokens(t *testing.T) {
	err := pdb.UpdateTokens(&repo.CafeTokens{Access: "access", Refresh: "refresh", Expiry: time.Now()})
	if err != nil {
		t.Error(err)
	}
}

func TestProfileDB_Logout(t *testing.T) {
	if err := pdb.Logout(); err != nil {
		t.Error(err)
		return
	}
	tokens, err := pdb.GetTokens()
	if err != nil {
		t.Error(err)
	}
	if tokens != nil {
		t.Error("logged out but tokens still present")
	}
}
