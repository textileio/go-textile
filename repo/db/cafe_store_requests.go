package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type CafeStoreRequestDB struct {
	modelStore
}

func NewCafeStoreRequestStore(db *sql.DB, lock *sync.Mutex) repo.CafeStoreRequestStore {
	return &CafeStoreRequestDB{modelStore{db, lock}}
}

func (c *CafeStoreRequestDB) Put(req *repo.CafeStoreRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into storereqs(id, cafeId, date) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		req.Id,
		req.CafeId,
		int(req.Date.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeStoreRequestDB) List(offset string, limit int) []repo.CafeStoreRequest {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from storereqs where date<(select date from storereqs where id='" + offset + "') order by date desc limit " + strconv.Itoa(limit) + " ;"
	} else {
		stm = "select * from storereqs order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *CafeStoreRequestDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from storereqs where id=?", id)
	return err
}

func (c *CafeStoreRequestDB) DeleteByCafe(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from storereqs where cafeId=?", cafeId)
	return err
}

func (c *CafeStoreRequestDB) handleQuery(stm string) []repo.CafeStoreRequest {
	var ret []repo.CafeStoreRequest
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, cafeId string
		var dateInt int
		if err := rows.Scan(&id, &cafeId, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		req := repo.CafeStoreRequest{
			Id:     id,
			CafeId: cafeId,
			Date:   time.Unix(int64(dateInt), 0),
		}
		ret = append(ret, req)
	}
	return ret
}
