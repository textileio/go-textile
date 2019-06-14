package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor003 struct{}

func (Minor003) Up(repoPath string, pinCode string, testnet bool) error {
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

	// add column for username and category to notifications
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt1, err := tx.Prepare("alter table notifications add column actorUn text not null default '';")
	if err != nil {
		return err
	}
	defer stmt1.Close()
	_, err = stmt1.Exec()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	stmt2, err := tx.Prepare("alter table notifications add column category text not null default '';")
	if err != nil {
		return err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()

	// update version
	f4, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f4.Close()
	if _, err = f4.Write([]byte("4")); err != nil {
		return err
	}
	return nil
}

func (Minor003) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor003) Major() bool {
	return false
}
