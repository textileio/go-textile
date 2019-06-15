package migrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/textileio/go-textile/keypair"
)

const testPeerID = "QmQA7swSsZKoayPHaTPgzZ1u3SCQjLvLyKcN6RRMmTbLau"

func initAt005(db *sql.DB, pin string) error {
	// Only need the Identity stub...
	configStr := fmt.Sprintf(`
	{
		"Identity": {
			"PeerID": "%s",
			"PrivKey": "CAESYH1jZmeyepc6aWdAeOkLbkVDYt5FFHIvQramNAGglovRHSxkSGg54g2KJJ/9oqFXJuw2WL009Gap3XnFUxnvKGodLGRIaDniDYokn/2ioVcm7DZYvTT0ZqndecVTGe8oag=="
		}
	}
	`, testPeerID)
	if err := ioutil.WriteFile("./config", []byte(configStr), 0644); err != nil {
		return err
	}
	var sqlStmt string
	if pin != "" {
		sqlStmt = "PRAGMA key = '" + pin + "';"
	}
	sqlStmt += `
    create table config (key text primary key not null, value blob);
    create table profile (key text primary key not null, value blob);
    create table cafe_sessions (cafeId text primary key not null, access text not null, refresh text not null, expiry integer not null, cafe blob not null);
    create table contacts (id text primary key not null, address text not null, username text not null, avatar text not null, inboxes blob not null, created integer not null, updated integer not null);
    create index contact_address on contacts (address);
    create index contact_username on contacts (username);
    create index contact_updated on contacts (updated);
	`
	if _, err := db.Exec(sqlStmt); err != nil {
		return err
	}
	accnt := keypair.Random()
	_, err := db.Exec("insert into config(key, value) values(?,?)", "seed", accnt.Seed())
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into profile(key, value) values(?,?)", "username", []byte("username"))
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into profile(key, value) values(?,?)", "avatar", []byte("/ipfs/Qm123"))
	if err != nil {
		return err
	}
	session := cafeSession{
		Id:      "12D3KooWJq2e9xWrccbyfY3MnXBZyDDKiKV8vZ1Rt4HtofzUAMe6",
		Access:  "eyJhbGciOiJFZDI1NTE5IiwidHlwIjoiSldUIn0.eyJzY29wZXMiOiJhY2Nlc3MiLCJhdWQiOiIvdGV4dGlsZS9jYWZlLzEuMC4wIiwiZXhwIjoxNTUwMDM0MTU3LCJqdGkiOiIxRnBuWXFFOFZpV2lDVXFKM08wdEN6VXoySk4iLCJpYXQiOjE1NDc2MTQ5NTcsImlzcyI6IjEyRDNLb29XSnEyZTl4V3JjY2J5ZlkzTW5YQlp5RERLaUtWOHZaMVJ0NEh0b2Z6VUFNZTYiLCJzdWIiOiIxMkQzS29vV0FiRHpmV2k5OTRGTHpLRHBwb0tKa2l2WmlUUWh3aGJTcXc4NHQyNVp5M0I3In0.496zRbI-MdFRy98lH_w-QLUhFrmoKtCGxYKSUri3HCFQj6Oac7fFqpFwv6AM3o2RlzVazq18KqFWR-2sDbt4Bw",
		Refresh: "eyJhbGciOiJFZDI1NTE5IiwidHlwIjoiSldUIn0.eyJzY29wZXMiOiJyZWZyZXNoIiwiYXVkIjoiL3RleHRpbGUvY2FmZS8xLjAuMCIsImV4cCI6MTU1MjQ1MzM1NywianRpIjoicjFGcG5ZcUU4VmlXaUNVcUozTzB0Q3pVejJKTiIsImlhdCI6MTU0NzYxNDk1NywiaXNzIjoiMTJEM0tvb1dKcTJlOXhXcmNjYnlmWTNNblhCWnlEREtpS1Y4dloxUnQ0SHRvZnpVQU1lNiIsInN1YiI6IjEyRDNLb29XQWJEemZXaTk5NEZMektEcHBvS0praXZaaVRRaHdoYlNxdzg0dDI1WnkzQjcifQ.NMcTVmFDQ4xnBGGgpHYP7BKCW0ZzZyfLAc3z5rE9N5B0vlPKDIarplFX3APW1236L11CCrrP7squRZDBLgfnCg",
		Expiry:  time.Now(),
		Cafe: cafe{
			Peer:     "12D3KooWJq2e9xWrccbyfY3MnXBZyDDKiKV8vZ1Rt4HtofzUAMe6",
			Address:  "P8CUU863id8h1HeYaKJyP3Vj5YvT6jYNepYTg7ubWrM6aQPu",
			API:      "v0",
			Protocol: "/textile/cafe/1.0.0",
			Node:     "1.0.0-rc25",
			URL:      "http://127.0.0.1:42601",
			Swarm:    []string{"/ip4/127.0.0.1/tcp/5789"},
		},
	}
	cafe, err := json.Marshal(session.Cafe)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into cafe_sessions(cafeId, access, refresh, expiry, cafe) values(?,?,?,?,?)",
		session.Id,
		session.Access,
		session.Refresh,
		session.Expiry.UnixNano(),
		cafe,
	)
	if err != nil {
		return err
	}
	return nil
}

func Test006(t *testing.T) {
	var dbPath string
	_ = os.Mkdir("./datastore", os.ModePerm)
	dbPath = path.Join("./", "datastore", "mainnet.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	if err := initAt005(db, ""); err != nil {
		t.Error(err)
		return
	}

	// go up
	var m Minor006
	if err := m.Up("./", "", false); err != nil {
		t.Error(err)
		return
	}

	// look for own contact
	var ret []contact
	rows, err := db.Query("select * from contacts where id='" + testPeerID + "';")
	if err != nil {
		t.Error(err)
		return
	}
	for rows.Next() {
		var id, address, username, avatar string
		var inboxes []byte
		var createdInt, updatedInt int64
		if err := rows.Scan(&id, &address, &username, &avatar, &inboxes, &createdInt, &updatedInt); err != nil {
			t.Error(err)
			return
		}

		ilist := make([]cafe, 0)
		if err := json.Unmarshal(inboxes, &ilist); err != nil {
			t.Error(err)
			return
		}

		ret = append(ret, contact{
			Id:       id,
			Address:  address,
			Username: username,
			Avatar:   avatar,
			Inboxes:  ilist,
			Created:  time.Unix(0, createdInt),
			Updated:  time.Unix(0, updatedInt),
		})
	}
	if len(ret) == 0 {
		t.Error("failed to find contact for self")
		return
	}

	// ensure that version file was updated
	version, err := ioutil.ReadFile("./repover")
	if err != nil {
		t.Error(err)
		return
	}
	if string(version) != "7" {
		t.Error("failed to write new repo version")
		return
	}

	if err := m.Down("./", "", false); err != nil {
		t.Error(err)
		return
	}
	_ = os.RemoveAll("./datastore")
	_ = os.RemoveAll("./repover")
	_ = os.RemoveAll("./config")
}
