package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"time"
)

type CafeNonceDB struct {
	modelStore
}

func NewCafeNonceStore(db *sql.DB, lock *sync.Mutex) repo.CafeNonceStore {
	return &CafeNonceDB{modelStore{db, lock}}
}

func (c *CafeNonceDB) Add(nonce *repo.CafeNonce) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_nonces(value, address, date) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		nonce.Value,
		nonce.Address,
		int(nonce.Date.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeNonceDB) Get(value string) *repo.CafeNonce {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from cafe_nonces where value='" + value + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeNonceDB) Delete(value string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_nonces where value=?", value)
	return err
}

func (c *CafeNonceDB) handleQuery(stm string) []repo.CafeNonce {
	var ret []repo.CafeNonce
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var value, address string
		var dateInt int
		if err := rows.Scan(&value, &address, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		nonce := repo.CafeNonce{
			Value:   value,
			Address: address,
			Date:    time.Unix(int64(dateInt), 0),
		}
		ret = append(ret, nonce)
	}
	return ret
}
