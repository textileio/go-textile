package migrations

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	native "github.com/ipfs/go-ipfs-config"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
)

type thread struct {
	id   string
	name string
	sk   []byte
}

type threadRow struct {
	Name  string   `json:"name"`
	Sk    string   `json:"sk"`
	Peers []string `json:"peers"`
}

type photoRow struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type Major005 struct{}

func (Major005) Up(repoPath string, pinCode string, testnet bool) error {
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

	// Get PeerId from IPFS config
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

	// Get username
	username := config.Identity.PeerID[len(config.Identity.PeerID)-7:]
	row := db.QueryRow("select value from profile where key='username';")
	_ = row.Scan(&username)
	jsonData := map[string]string{
		"peerid":   config.Identity.PeerID,
		"username": username,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(repoPath, "migration005_peerid.ndjson"), jsonBytes, 0644); err != nil {
		return err
	}

	// collect thread secrets
	var threads []*thread
	var defaults []*thread
	rows, err := db.Query("select id, name, sk from threads;")
	if err != nil {
		return err
	}
	for rows.Next() {
		var id, name string
		var sk []byte
		if err := rows.Scan(&id, &name, &sk); err != nil {
			return err
		}
		threads = append(threads, &thread{id: id, name: name, sk: sk})
	}

	// collect thread peers
	threadPeers := make(map[string][]string)
	for _, thread := range threads {
		rows, err := db.Query("select id from peers where threadId='" + thread.id + "'")
		if err != nil {
			return err
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return err
			}
			threadPeers[thread.id] = append(threadPeers[thread.id], id)
		}
	}

	// write to file
	tfile, err := os.Create(path.Join(repoPath, "migration005_threads.ndjson"))
	if err != nil {
		return err
	}
	defer tfile.Close()
	for _, thrd := range threads {
		if thrd.name == "default" {
			defaults = append(defaults, thrd)
		}
		sk64 := base64.StdEncoding.EncodeToString(thrd.sk)
		peers := threadPeers[thrd.id]
		if len(peers) == 0 {
			peers = make([]string, 0)
		}
		row, err := json.Marshal(&threadRow{
			Name:  thrd.name,
			Sk:    sk64,
			Peers: peers,
		})
		if err != nil {
			return err
		}
		if _, err = tfile.Write(append(row[:], []byte("\n")[:]...)); err != nil {
			return err
		}
	}

	// collect default thread photo blocks
	var photos []*photoRow
	for _, thread := range defaults {
		sk, err := ipfs.UnmarshalPrivateKey(thread.sk)
		if err != nil {
			return err
		}
		rows, err := db.Query("select dataId, dataKeyCipher from blocks where threadId='" + thread.id + "' and type=4;")
		if err != nil {
			return err
		}
		for rows.Next() {
			var id string
			var keyCipher []byte
			if err := rows.Scan(&id, &keyCipher); err != nil {
				return err
			}
			key, err := crypto.Decrypt(sk, keyCipher)
			if err != nil {
				return err
			}
			photos = append(photos, &photoRow{Id: id, Key: string(key)})
		}
	}

	// write to file
	bfile, err := os.Create(path.Join(repoPath, "migration005_default_photos.ndjson"))
	if err != nil {
		return err
	}
	defer bfile.Close()
	for _, photo := range photos {
		row, err := json.Marshal(photo)
		if err != nil {
			return err
		}
		if _, err = bfile.Write(append(row[:], []byte("\n")[:]...)); err != nil {
			return err
		}
	}

	// blast repo sans logs
	return blastRepo(repoPath)
}

func (Major005) Down(repoPath string, pinCode string, testnet bool) error {
	return ErrorCannotMigrateDown
}

func (Major005) Major() bool {
	return true
}
