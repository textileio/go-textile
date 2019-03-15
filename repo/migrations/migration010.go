package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor010 struct{}

func (Minor010) Up(repoPath string, pinCode string, testnet bool) error {
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
    drop table thread_messages;
    drop table thread_invites;
    create table block_messages (id text primary key not null, peerId text not null, envelope blob not null, date integer not null);
    create index block_message_date on block_messages (date);
    create table invites (id text primary key not null, block blob not null, name text not null, inviter blob not null, date integer not null);
    create index invite_date on invites (date);
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f11, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f11.Close()
	if _, err = f11.Write([]byte("11")); err != nil {
		return err
	}
	return nil
}

// Down is for a migration downgrade (not implemented)
func (Minor010) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

// Major is for a major version migration change (not implemented)
func (Minor010) Major() bool {
	return false
}
