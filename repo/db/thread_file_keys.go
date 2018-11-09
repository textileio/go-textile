package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type ThreadFileKeyDB struct {
	modelStore
}

func NewThreadFileKeyStore(db *sql.DB, lock *sync.Mutex) repo.ThreadFileKeyStore {
	return &ThreadFileKeyDB{modelStore{db, lock}}
}

func (c *ThreadFileKeyDB) Add(key *repo.ThreadFileKey) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into thread_file_keys(hash, key) values(?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		key.Hash,
		key.Key,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ThreadFileKeyDB) Get(hash string) *repo.ThreadFileKey {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from thread_file_keys where hash='" + hash + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadFileKeyDB) Delete(hash string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_file_keys where hash=?", hash)
	return err
}

func (c *ThreadFileKeyDB) handleQuery(stm string) []repo.ThreadFileKey {
	var res []repo.ThreadFileKey
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var hash, key string
		if err := rows.Scan(&hash, &key); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		res = append(res, repo.ThreadFileKey{
			Hash: hash,
			Key:  key,
		})
	}
	return res
}
