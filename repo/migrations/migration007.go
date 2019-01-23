package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor007 struct{}

func (Minor007) Up(repoPath string, pinCode string, testnet bool) error {
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

	// delete thread invites
	if _, err := db.Exec("drop table thread_invites;"); err != nil {
		return err
	}

	// add it back
	query := `
    create table thread_invites (id text primary key not null, block blob not null, name text not null, contact blob not null, date integer not null);
    create index thread_invite_date on thread_invites (date);
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f8, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f8.Close()
	if _, err = f8.Write([]byte("8")); err != nil {
		return err
	}
	return nil
}

func (Minor007) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor007) Major() bool {
	return false
}
