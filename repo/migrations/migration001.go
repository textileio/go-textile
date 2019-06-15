package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor001 struct{}

func (Minor001) Up(repoPath string, pinCode string, testnet bool) error {
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

	// add column for encrypted metadata to blocks
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("alter table blocks add column dataMetadataCipher blob;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()

	// update version
	f2, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f2.Close()
	if _, err = f2.Write([]byte("2")); err != nil {
		return err
	}
	return nil
}

func (Minor001) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor001) Major() bool {
	return false
}
