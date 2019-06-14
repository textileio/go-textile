package migrations

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mutecomm/go-sqlcipher"
)

type Minor002 struct{}

func (Minor002) Up(repoPath string, pinCode string, testnet bool) error {
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

	// add notifications table and indexes
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := `
    create table notifications (id text primary key not null, date integer not null, actorId text not null, targetId text not null, type integer not null, read integer not null, body text not null);
    create index notification_targetId on notifications (targetId);
    create index notification_actorId on notifications (actorId);
    create index notification_read on notifications (read);
    `
	stmt, err := tx.Prepare(query)
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
	f3, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f3.Close()
	if _, err = f3.Write([]byte("3")); err != nil {
		return err
	}
	return nil
}

func (Minor002) Down(repoPath string, pinCode string, testnet bool) error {
	return nil
}

func (Minor002) Major() bool {
	return false
}
