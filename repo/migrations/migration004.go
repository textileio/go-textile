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

	// add column categoryId to notifications
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt1, err := tx.Prepare("alter table notifications add column categoryId text not null default '';")
	if err != nil {
		return err
	}
	defer stmt1.Close()
	_, err = stmt1.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	// add index on categoryId
	stmt2, err := tx.Prepare("create index notification_categoryId on notifications (categoryId);")
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
