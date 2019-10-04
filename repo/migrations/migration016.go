package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor016 struct{}

func (Minor016) Up(repoPath string, pinCode string, testnet bool) error {
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
    create table botstore (id text primary key not null, key text not null, value blob, version integer not null, created integer not null, updated integer not null);
    create index botstore_key on botstore (key);
    create index botstore_updated on botstore (updated);
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f17, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f17.Close()
	if _, err = f17.Write([]byte("17")); err != nil {
		return err
	}
	return nil
}

func (Minor016) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor016) Major() bool {
	return false
}
