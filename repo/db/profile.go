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
	tx.Commit()
	return nil
}

func (c *ProfileDB) SignOut() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("delete from profile where key=?")
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

func (c *ProfileDB) GetUsername() (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	defer stmt.Close()
	var username string
	err = stmt.QueryRow("username").Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
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

func (c *ProfileDB) GetAvatarId() (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stmt, err := c.db.Prepare("select value from profile where key=?")
	defer stmt.Close()
	var avatarId string
	err = stmt.QueryRow("avatar_id").Scan(&avatarId)
	if err != nil {
		return "", err
	}
	return avatarId, nil
}

func (c *ProfileDB) GetTokens() (tokens *repo.CafeTokens, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stmt *sql.Stmt
	stmt, err = c.db.Prepare("select value from profile where key=?")
	if err != nil {
		return
	}
	defer stmt.Close()
	defer func() {
		if recover() != nil {
			log.Warning("get tokens recovered")
		}
	}()
	var accessToken, refreshToken string
	err = stmt.QueryRow("access").Scan(&accessToken)
	if err != nil {
		return
	}
	err = stmt.QueryRow("refresh").Scan(&refreshToken)
	if err != nil {
		return
	}
	tokens = &repo.CafeTokens{Access: accessToken, Refresh: refreshToken}
	return
}
