package migrations

import (
	"database/sql"
	"encoding/base64"
	_ "github.com/mutecomm/go-sqlcipher"
	"os"
	"path"
)

type thread struct {
	name string
	sk   []byte
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

	// collect thread secrets
	var threads []*thread
	rows, err := db.Query("select name, sk from threads;")
	if err != nil {
		return err
	}
	for rows.Next() {
		var name string
		var sk []byte
		if err := rows.Scan(&name, &sk); err != nil {
			return err
		}
		threads = append(threads, &thread{name: name, sk: sk})
	}

	// write to file
	file, err := os.Create(path.Join(repoPath, "major00_threads"))
	if err != nil {
		return err
	}
	defer file.Close()
	for _, thread := range threads {
		sk64 := base64.StdEncoding.EncodeToString(thread.sk)
		if _, err = file.Write([]byte(thread.name + " " + sk64 + "\n")); err != nil {
			return err
		}
	}

	// blast repo sans logs
	if err := blastRepo(repoPath); err != nil {
		return err
	}

	// update version
	f6, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f6.Close()
	if _, err = f6.Write([]byte("6")); err != nil {
		return err
	}
	return nil
}

func (Major005) Down(repoPath string, pinCode string, testnet bool) error {
	return ErrorCannotMigrateDown
}

func (Major005) Major() bool {
	return true
}
