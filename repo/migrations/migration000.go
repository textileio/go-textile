package migrations

import (
	"database/sql"
	_ "github.com/mutecomm/go-sqlcipher"
	"os"
	"path"
)

type Minor000 struct{}

func (Minor000) Up(repoPath string, pinCode string, testnet bool) error {
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
		p := "pragma key='" + pinCode + "';"
		if _, err := db.Exec(p); err != nil {
			return err
		}
	}

	// add column for encrypted username to blocks
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("alter table blocks add column dataUsernameCipher blob;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// update version
	f1, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f1.Close()
	if _, err = f1.Write([]byte("1")); err != nil {
		return err
	}
	return nil
}

func (Minor000) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor000) Major() bool {
	return false
}
