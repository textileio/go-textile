package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor017 struct{}

func (Minor017) Up(repoPath string, pinCode string, testnet bool) error {
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

	query := `
		drop table if exists botstore;
		create table bots_store (id text primary key not null, value blob, created integer not null, updated integer not null);
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f18, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f18.Close()
	if _, err = f18.Write([]byte("18")); err != nil {
		return err
	}
	return nil
}

func (Minor017) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor017) Major() bool {
	return false
}
