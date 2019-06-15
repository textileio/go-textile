package db

import (
	"database/sql"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeClientNonceDB struct {
	modelStore
}

func NewCafeClientNonceStore(db *sql.DB, lock *sync.Mutex) repo.CafeClientNonceStore {
	return &CafeClientNonceDB{modelStore{db, lock}}
}

func (c *CafeClientNonceDB) Add(nonce *pb.CafeClientNonce) error {
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
		util.ProtoNanos(nonce.Date),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *CafeClientNonceDB) Get(value string) *pb.CafeClientNonce {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from cafe_client_nonces where value='" + value + "';")
	if len(res) == 0 {
		return nil
	}
	return &res[0]
}

func (c *CafeClientNonceDB) Delete(value string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_client_nonces where value=?", value)
	return err
}

func (c *CafeClientNonceDB) handleQuery(stm string) []pb.CafeClientNonce {
	var list []pb.CafeClientNonce
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
		list = append(list, pb.CafeClientNonce{
			Value:   value,
			Address: address,
			Date:    util.ProtoTs(dateInt),
		})
	}
	return list
}
