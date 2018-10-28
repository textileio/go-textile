package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type CafeClientThreadDB struct {
	modelStore
}

func NewCafeClientThreadStore(db *sql.DB, lock *sync.Mutex) repo.CafeClientThreadStore {
	return &CafeClientThreadDB{modelStore{db, lock}}
}

func (c *CafeClientThreadDB) AddOrUpdate(thrd *repo.CafeClientThread) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into cafe_client_threads(id, clientId, skCipher, headCipher, nameCipher) values(?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		thrd.Id,
		thrd.ClientId,
		thrd.SkCipher,
		thrd.HeadCipher,
		thrd.NameCipher,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeClientThreadDB) ListByClient(clientId string) []repo.CafeClientThread {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_client_threads where clientId='" + clientId + "';"
	return c.handleQuery(stm)
}

func (c *CafeClientThreadDB) Delete(id string, clientId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_client_threads where id=? and clientId=?", id, clientId)
	return err
}

func (c *CafeClientThreadDB) DeleteByClient(clientId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_client_threads where clientId=?", clientId)
	return err
}

func (c *CafeClientThreadDB) handleQuery(stm string) []repo.CafeClientThread {
	var ret []repo.CafeClientThread
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, clientId string
		var skCipher, headCipher, nameCipher []byte
		if err := rows.Scan(&id, &clientId, &skCipher, &headCipher, &nameCipher); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		thrd := repo.CafeClientThread{
			Id:         id,
			ClientId:   clientId,
			SkCipher:   skCipher,
			HeadCipher: headCipher,
			NameCipher: nameCipher,
		}
		ret = append(ret, thrd)
	}
	return ret
}
