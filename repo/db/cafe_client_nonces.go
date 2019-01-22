package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type CafeClientNonceDB struct {
	modelStore
}

func NewCafeClientNonceStore(db *sql.DB, lock *sync.Mutex) repo.CafeClientNonceStore {
	return &CafeClientNonceDB{modelStore{db, lock}}
}

func (c *CafeClientNonceDB) Add(nonce *repo.CafeClientNonce) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_client_nonces(value, address, date) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		nonce.Value,
		nonce.Address,
		nonce.Date.UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeClientNonceDB) Get(value string) *repo.CafeClientNonce {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from cafe_client_nonces where value='" + value + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeClientNonceDB) Delete(value string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_client_nonces where value=?", value)
	return err
}

func (c *CafeClientNonceDB) handleQuery(stm string) []repo.CafeClientNonce {
	var ret []repo.CafeClientNonce
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var value, address string
		var dateInt int64
		if err := rows.Scan(&value, &address, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		ret = append(ret, repo.CafeClientNonce{
			Value:   value,
			Address: address,
			Date:    time.Unix(0, dateInt),
		})
	}
	return ret
}
