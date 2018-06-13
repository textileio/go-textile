package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type ProfileDB struct {
	db   *sql.DB
	lock *sync.Mutex
}

func NewProfileStore(db *sql.DB, lock *sync.Mutex) repo.ProfileStore {
	return &ProfileDB{db, lock}
}

func (c *ProfileDB) SignIn(id string, secret []byte, username string, accessToken string, refreshToken string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert or replace into profile(key, value) values(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec("id", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("secret", secret)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("username", username)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("access", accessToken)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("refresh", refreshToken)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ProfileDB) SignOut() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("delete from profile where key=?")
	defer stmt.Close()
	_, err = stmt.Exec("id")
	if err != nil {
		return err
	}
	_, err = stmt.Exec("secret")
	if err != nil {
		return err
	}
	_, err = stmt.Exec("username")
	if err != nil {
		return err
	}
	_, err = stmt.Exec("access")
	if err != nil {
		return err
	}
	_, err = stmt.Exec("refresh")
	if err != nil {
		return err
	}
	return nil
}

func (c *ProfileDB) GetID() (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	defer stmt.Close()
	var id string
	err = stmt.QueryRow("id").Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (c *ProfileDB) GetSecret() ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	defer stmt.Close()
	var secret []byte
	err = stmt.QueryRow("secret").Scan(&secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (c *ProfileDB) GetUsername() (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	defer stmt.Close()
	var un string
	err = stmt.QueryRow("username").Scan(&un)
	if err != nil {
		return "", err
	}
	return un, nil
}

func (c *ProfileDB) GetTokens() (accessToken string, refreshToken string, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	defer stmt.Close()
	err = stmt.QueryRow("access").Scan(&accessToken)
	if err != nil {
		return "", "", err
	}
	err = stmt.QueryRow("refresh").Scan(&refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
