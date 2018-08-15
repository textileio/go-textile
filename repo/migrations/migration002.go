package migrations

import (
	"database/sql"
	_ "github.com/mutecomm/go-sqlcipher"
	"os"
	"path"
)

type Migration002 struct{}

func (Migration002) Up(repoPath string, dbPassword string, testnet bool) error {
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
	if dbPassword != "" {
		p := "pragma key='" + dbPassword + "';"
		if _, err := db.Exec(p); err != nil {
			return err
		}
	}

	// add notifications table and indexes
	query := `
    create table notifications (id text primary key not null, date integer not null, actorId text not null, targetId text not null, type integer not null, read integer not null);
    create index notification_targetId on notifications (targetId);
    create index notification_actorId on notifications (actorId);
    create index notification_read on notifications (read);
    `
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

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

func (Migration002) Down(repoPath string, dbPassword string, testnet bool) error {
	return nil
}
