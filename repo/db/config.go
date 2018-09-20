package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
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

func (c *ConfigDB) Init(pin string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return initDatabaseTables(c.db, pin)
}

func (c *ConfigDB) Configure(key libp2pc.PrivKey, created time.Time) error {
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
	keyb, err := libp2pc.MarshalPrivateKey(key)
	if err != nil {
		return err
	}
	id, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return err
	}
	_, err = stmt.Exec("id", id.Pretty())
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("key", keyb)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("created", created.Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ConfigDB) GetId() (*string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("id").Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &id, nil
}

func (c *ConfigDB) GetKey() (libp2pc.PrivKey, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var keyb []byte
	if err := stmt.QueryRow("key").Scan(&keyb); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return libp2pc.UnmarshalPrivateKey(keyb)
}

func (c *ConfigDB) GetCreationDate() (time.Time, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var t time.Time
	stmt, err := c.db.Prepare("select value from config where key=?")
	if err != nil {
		if err == sql.ErrNoRows {
			return t, nil
		}
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
