package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"time"
)

type ConfigDB struct {
	db   *sql.DB
	lock *sync.Mutex
	path string
}

func NewConfigStore(db *sql.DB, lock *sync.Mutex, path string) repo.ConfigStore {
	return &ConfigDB{db, lock, path}
}

func (c *ConfigDB) Init(password string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return initDatabaseTables(c.db, password)
}

func (c *ConfigDB) Configure(created time.Time, version string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into config(key, value) values(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec("created", created.Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("version", version)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ConfigDB) GetCreationDate() (time.Time, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var t time.Time
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		return t, err
	}
	defer stmt.Close()
	var created []byte
	err = stmt.QueryRow("created").Scan(&created)
	if err != nil {
		return t, err
	}
	return time.Parse(time.RFC3339, string(created))
}

func (c *ConfigDB) GetVersion() (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var sv string
	err = stmt.QueryRow("version").Scan(&sv)
	if err != nil {
		return "", err
	}
	return sv, nil
}

func (c *ConfigDB) IsEncrypted() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	pwdCheck := "select count(*) from sqlite_master;"
	_, err := c.db.Exec(pwdCheck) // Fails if wrong password is entered
	if err != nil {
		return true
	}
	return false
}
