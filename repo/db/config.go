package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/strkey"
)

type ConfigDB struct {
	db   *sql.DB
	lock *sync.Mutex
	path string
}

func NewConfigStore(db *sql.DB, lock *sync.Mutex, path string) repo.ConfigStore {
	return &ConfigDB{db, lock, path}
}

func (c *ConfigDB) Init(pin string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return initDatabaseTables(c.db, pin)
}

func (c *ConfigDB) Configure(accnt *keypair.Full, created time.Time) error {
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
	if _, err = stmt.Exec("seed", accnt.Seed()); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err = stmt.Exec("created", created.Format(time.RFC3339)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *ConfigDB) GetAccount() (*keypair.Full, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var seed string
	if err := stmt.QueryRow("seed").Scan(&seed); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if _, err = strkey.Decode(strkey.VersionByteSeed, seed); err != nil {
		return nil, err
	}
	kp, err := keypair.Parse(seed)
	if err != nil {
		return nil, err
	}
	full, ok := kp.(*keypair.Full)
	if !ok {
		return nil, fmt.Errorf("invalid seed")
	}
	return full, nil
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
	if err := stmt.QueryRow("created").Scan(&created); err != nil {
		return t, err
	}
	return time.Parse(time.RFC3339, string(created))
}

func (c *ConfigDB) IsEncrypted() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	pwdCheck := "select count(*) from sqlite_master;"
	if _, err := c.db.Exec(pwdCheck); err != nil {
		return true // wrong password
	}
	return false
}

func (c *ConfigDB) GetLastDaily() (time.Time, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var t time.Time
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		return t, err
	}
	defer stmt.Close()
	var last []byte
	if err := stmt.QueryRow("daily").Scan(&last); err != nil {
		if err == sql.ErrNoRows {
			return t, nil
		}
		return t, err
	}
	return time.Parse(time.RFC3339, string(last))
}

func (c *ConfigDB) SetLastDaily() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert or replace into config(key, value) values(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec("daily", time.Now().Format(time.RFC3339)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
