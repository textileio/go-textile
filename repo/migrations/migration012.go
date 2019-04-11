package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor012 struct{}

func (Minor012) Up(repoPath string, pinCode string, testnet bool) error {
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
    alter table cafe_requests add column size integer not null default 0;
    alter table cafe_requests add column groupId text not null default '';
    alter table cafe_requests add column status integer not null default 0;
    create index cafe_request_groupId on cafe_requests (groupId);
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f13, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f13.Close()
	if _, err = f13.Write([]byte("13")); err != nil {
		return err
	}
	return nil
}

func (Minor012) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor012) Major() bool {
	return false
}
