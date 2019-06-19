package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor015 struct{}

func (Minor015) Up(repoPath string, pinCode string, testnet bool) error {
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
		_, err = db.Exec("pragma key='" + pinCode + "';")
		if err != nil {
			return err
		}
	}

	query := `
    alter table cafe_requests add column groupSize integer not null default 0;
    alter table cafe_requests add column groupTransferred integer not null default 0;
    `
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	// update version
	f16, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f16.Close()
	if _, err = f16.Write([]byte("16")); err != nil {
		return err
	}
	return nil
}

func (Minor015) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor015) Major() bool {
	return false
}
