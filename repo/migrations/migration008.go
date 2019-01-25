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

	// add missing cafe_dev_token table
	// TODO: should we also edit the cafe_client table to add a field linking to cafe_dev_token?
	query := `
		create table cafe_dev_tokens (id text primary key not null, token blob not null, created integer not null);
		`
	if _, err := db.Exec(query); err != nil {
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

// Down is for a migration downgrade (not implemented)
func (Minor008) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

// Major is for a major version migration change (not implemented)
func (Minor008) Major() bool {
	return false
}
