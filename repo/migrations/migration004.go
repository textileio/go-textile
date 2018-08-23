package migrations

import (
	"database/sql"
	_ "github.com/mutecomm/go-sqlcipher"
	"os"
	"path"
)

type Migration004 struct{}

func (Migration004) Up(repoPath string, dbPassword string, testnet bool) error {
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

	// delete notifications table
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt1, err := tx.Prepare("drop table notifications;")
	if err != nil {
		return err
	}
	defer stmt1.Close()
	_, err = stmt1.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	// re-add notifications table and indexes
	query := `
    create table notifications (id text primary key not null, date integer not null, actorId text not null, actorUsername text not null, subject text not null, subjectId text not null, blockId text, dataId text, type integer not null, body text not null, read integer not null);
    create index notification_actorId on notifications (actorId);    
    create index notification_subjectId on notifications (subjectId);
    create index notification_blockId on notifications (blockId);
    create index notification_read on notifications (read);
    `
	stmt2, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// update version
	f5, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f5.Close()
	if _, err = f5.Write([]byte("5")); err != nil {
		return err
	}
	return nil
}

func (Migration004) Down(repoPath string, dbPassword string, testnet bool) error {
	return nil
}
