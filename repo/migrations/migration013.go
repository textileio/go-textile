package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor013 struct{}

func (Minor013) Up(repoPath string, pinCode string, testnet bool) error {
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
    alter table invites add column parents text not null default '';
    drop table cafe_requests;
    create table cafe_requests (id text primary key not null, peerId text not null, targetId text not null, cafeId text not null, cafe blob not null, groupId text not null, syncGroupId text not null, type integer not null, date integer not null, size integer not null, status integer not null, attempts integer not null);
    create index cafe_request_cafeId on cafe_requests (cafeId);
    create index cafe_request_groupId on cafe_requests (groupId);
    create index cafe_request_syncGroupId on cafe_requests (syncGroupId);
    create index cafe_request_date on cafe_requests (date);
    create index cafe_request_status on cafe_requests (status);
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f14, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f14.Close()
	if _, err = f14.Write([]byte("14")); err != nil {
		return err
	}
	return nil
}

func (Minor013) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor013) Major() bool {
	return false
}
