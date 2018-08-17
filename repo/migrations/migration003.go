package migrations

import (
	"database/sql"
	_ "github.com/mutecomm/go-sqlcipher"
	"os"
	"path"
)

type Migration003 struct{}

func (Migration003) Up(repoPath string, dbPassword string, testnet bool) error {
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

	// add column for username and category to notifications
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt1, err := tx.Prepare("alter table notifications add column actorUn text not null default '';")
	if err != nil {
		return err
	}
	defer stmt1.Close()
	_, err = stmt1.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	stmt2, err := tx.Prepare("alter table notifications add column category text not null default '';")
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
	f4, err := os.Create(path.Join(repoPath, "repover"))
	if err != nil {
		return err
	}
	defer f4.Close()
	if _, err = f4.Write([]byte("4")); err != nil {
		return err
	}
	return nil
}

func (Migration003) Down(repoPath string, dbPassword string, testnet bool) error {
	return nil
}
