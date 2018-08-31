package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"time"
)

type ProfileDB struct {
	db   *sql.DB
	lock *sync.Mutex
}

func NewProfileStore(db *sql.DB, lock *sync.Mutex) repo.ProfileStore {
	return &ProfileDB{db, lock}
}

func (c *ProfileDB) SignIn(username string, tokens *repo.CafeTokens) error {
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
	_, err = stmt.Exec("username", username)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("access", tokens.Access)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("refresh", tokens.Refresh)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("expiry", int(tokens.Expiry.Unix()))
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
	if err != nil {
		return err
	}
	defer stmt.Close()
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

func (c *ProfileDB) GetUsername() (*string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var username string
	if err := stmt.QueryRow("username").Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &username, nil
}

func (c *ProfileDB) SetAvatarId(id string) error {
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
	_, err = stmt.Exec("avatar_id", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ProfileDB) GetAvatarId() (*string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var avatarId string
	if err := stmt.QueryRow("avatar_id").Scan(&avatarId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &avatarId, nil
}

func (c *ProfileDB) GetTokens() (*repo.CafeTokens, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var accessToken, refreshToken string
	var expiryInt int
	if err := stmt.QueryRow("access").Scan(&accessToken); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := stmt.QueryRow("refresh").Scan(&refreshToken); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := stmt.QueryRow("expiry").Scan(&expiryInt); err != nil {
		if err == sql.ErrNoRows {
			expiryInt = 0
		} else {
			return nil, err
		}
	}
	return &repo.CafeTokens{
		Access:  accessToken,
		Refresh: refreshToken,
		Expiry:  time.Unix(int64(expiryInt), 0),
	}, nil
}

func (c *ProfileDB) UpdateTokens(tokens *repo.CafeTokens) error {
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
	_, err = stmt.Exec("access", tokens.Access)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("refresh", tokens.Refresh)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmt.Exec("expiry", int(tokens.Expiry.Unix()))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
