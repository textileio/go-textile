package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor008 struct{}

func (Minor008) Up(repoPath string, pinCode string, testnet bool) error {
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

	// add column for members and sharing
	if _, err := db.Exec("alter table threads add column members text not null default '';"); err != nil {
		return err
	}
	if _, err := db.Exec("alter table threads add column sharing integer not null default 0;"); err != nil {
		return err
	}

	// update existing threads to have sharing == 2 (shared), where type == 3 (open)
	if _, err := db.Exec("update threads set sharing=2 where type=3;"); err != nil {
		return err
	}

	// update version
	f9, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f9.Close()
	if _, err = f9.Write([]byte("9")); err != nil {
		return err
	}
	return nil
}

func (Minor008) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor008) Major() bool {
	return false
}
