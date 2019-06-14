package migrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	native "github.com/ipfs/go-ipfs-config"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/strkey"
)

type contact struct {
	Id       string    `json:"id"`
	Address  string    `json:"address"`
	Username string    `json:"username,omitempty"`
	Avatar   string    `json:"avatar,omitempty"`
	Inboxes  []cafe    `json:"inboxes,omitempty"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type cafe struct {
	Peer     string   `json:"peer"`
	Address  string   `json:"address"`
	API      string   `json:"api"`
	Protocol string   `json:"protocol"`
	Node     string   `json:"node"`
	URL      string   `json:"url"`
	Swarm    []string `json:"swarm"`
}

type cafeSession struct {
	Id      string    `json:"id"`
	Access  string    `json:"access"`
	Refresh string    `json:"refresh"`
	Expiry  time.Time `json:"expiry"`
	Cafe    cafe      `json:"cafe"`
}

type Minor006 struct{}

func (Minor006) Up(repoPath string, pinCode string, testnet bool) error {
	var dbPath string
	if testnet {
		dbPath = path.Join(repoPath, "datastore", "testnet.db")
	} else {
		dbPath = path.Join(repoPath, "datastore", "mainnet.db")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	if pinCode != "" {
		if _, err := db.Exec("pragma key='" + pinCode + "';"); err != nil {
			return err
		}
	}

	// get peer id from IPFS config
	configPath := path.Join(repoPath, "config")
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	var config native.Config
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return err
	}

	// get address from config
	var seed string
	row := db.QueryRow("select value from config where key='seed';")
	if err := row.Scan(&seed); err != nil {
		return fmt.Errorf("error getting address: %s", err)
	}
	if _, err = strkey.Decode(strkey.VersionByteSeed, seed); err != nil {
		return err
	}
	kp, err := keypair.Parse(seed)
	if err != nil {
		return err
	}
	full, ok := kp.(*keypair.Full)
	if !ok {
		return fmt.Errorf("invalid seed")
	}

	// get username
	username := config.Identity.PeerID[len(config.Identity.PeerID)-7:]
	row1 := db.QueryRow("select value from profile where key='username';")
	_ = row1.Scan(&username)

	// get avatar
	var avatar string
	row2 := db.QueryRow("select value from profile where key='avatar';")
	if err := row2.Scan(&avatar); err == nil {
		avatar = strings.Replace(avatar, "/ipfs/", "", 1)
	}

	// get inboxes
	var sessions []cafeSession
	rows, err := db.Query(`select * from cafe_sessions order by expiry desc;`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var cafeId, access, refresh string
		var expiryInt int64
		var c []byte
		if err := rows.Scan(&cafeId, &access, &refresh, &expiryInt, &c); err != nil {
			return err
		}
		var rcafe cafe
		if err := json.Unmarshal(c, &rcafe); err != nil {
			return err
		}
		sessions = append(sessions, cafeSession{
			Id:      cafeId,
			Access:  access,
			Refresh: refresh,
			Expiry:  time.Unix(0, expiryInt),
			Cafe:    rcafe,
		})
	}
	var cafes []cafe
	for _, ses := range sessions {
		cafes = append(cafes, ses.Cafe)
	}
	inboxes, err := json.Marshal(cafes)
	if err != nil {
		return err
	}

	// add contact for self
	q2 := `insert into contacts(id, address, username, avatar, inboxes, created, updated) values(?,?,?,?,?,?,?)`
	if _, err = db.Exec(
		q2,
		config.Identity.PeerID,
		full.Address(),
		username,
		avatar,
		inboxes,
		time.Now().UnixNano(),
		time.Now().UnixNano(),
	); err != nil {
		if !conflictError(err) {
			return err
		}
	}

	// delete profile table
	_, _ = db.Exec("drop table profile;")

	// update version
	f7, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f7.Close()
	if _, err = f7.Write([]byte("7")); err != nil {
		return err
	}
	return nil
}

func (Minor006) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor006) Major() bool {
	return false
}
