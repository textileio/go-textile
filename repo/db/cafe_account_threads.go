package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type CafeAccountThreadDB struct {
	modelStore
}

func NewCafeAccountThreadStore(db *sql.DB, lock *sync.Mutex) repo.CafeAccountThreadStore {
	return &CafeAccountThreadDB{modelStore{db, lock}}
}

func (c *CafeAccountThreadDB) Add(thrd *repo.CafeAccountThread) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into account_threads(id, accountId, skCipher, head) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		thrd.Id,
		thrd.AccountId,
		thrd.SkCipher,
		thrd.Head,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeAccountThreadDB) Get(id string, accountId string) *repo.CafeAccountThread {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from account_threads where id='" + id + "' and accountId='" + accountId + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeAccountThreadDB) ListByAccount(accountId string) []repo.CafeAccountThread {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from account_threads where accountId='" + accountId + "';"
	return c.handleQuery(stm)
}

func (c *CafeAccountThreadDB) UpdateHead(id string, accountId string, head string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update account_threads set head=? where id=? and accountId=?", head, id, accountId)
	return err
}

func (c *CafeAccountThreadDB) Delete(id string, accountId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from account_threads where id=? and accountId=?", id, accountId)
	return err
}

func (c *CafeAccountThreadDB) handleQuery(stm string) []repo.CafeAccountThread {
	var ret []repo.CafeAccountThread
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, accountId string
		var head *string
		var skCipher []byte
		if err := rows.Scan(&id, &accountId, &skCipher, &head); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		thrd := repo.CafeAccountThread{
			Id:        id,
			AccountId: accountId,
			SkCipher:  skCipher,
			Head:      head,
		}
		ret = append(ret, thrd)
	}
	return ret
}
