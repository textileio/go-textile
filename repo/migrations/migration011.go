package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor011 struct{}

func (Minor011) Up(repoPath string, pinCode string, testnet bool) error {
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
    alter table contacts rename to peers;
    drop index contact_address;
    drop index contact_username;
    drop index contact_updated;
    create index peer_address on peers (address);
    create index peer_username on peers (username);
    create index peer_updated on peers (updated);
    `
	if _, err := db.Exec(query); err != nil {
		return err
	}

	// update version
	f12, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f12.Close()
	if _, err = f12.Write([]byte("12")); err != nil {
		return err
	}
	return nil
}

// Down is for a migration downgrade (not implemented)
func (Minor011) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

// Major is for a major version migration change (not implemented)
func (Minor011) Major() bool {
	return false
}
