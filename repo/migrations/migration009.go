package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor009 struct{}

func (Minor009) Up(repoPath string, pinCode string, testnet bool) error {
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
    create table cafe_tokens (id text primary key not null, token blob not null, date integer not null);
    alter table cafe_clients add column tokenId text not null default '';
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f10, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f10.Close()
	if _, err = f10.Write([]byte("10")); err != nil {
		return err
	}
	return nil
}

// Down is for a migration downgrade (not implemented)
func (Minor009) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

// Major is for a major version migration change (not implemented)
func (Minor009) Major() bool {
	return false
}
