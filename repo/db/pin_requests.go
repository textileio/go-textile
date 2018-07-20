package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type PinRequestDB struct {
	modelStore
}

func NewPinRequestStore(db *sql.DB, lock *sync.Mutex) repo.PinRequestStore {
	return &PinRequestDB{modelStore{db, lock}}
}

func (c *PinRequestDB) Put(pr *repo.PinRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into pinrequests(id, date) values(?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		pr.Id,
		int(pr.Date.Unix()),
	)
	if err != nil {
		tx.Rollback()
		log.Errorf("error in db exec: %s", err)
		return err
	}
	tx.Commit()
	return nil
}

func (c *PinRequestDB) List(offset string, limit int) []repo.PinRequest {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from pinrequests where date<(select date from pinrequests where id='" + offset + "') order by date desc limit " + strconv.Itoa(limit) + " ;"
	} else {
		stm = "select * from pinrequests order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *PinRequestDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from pinrequests where id=?", id)
	return err
}

func (c *PinRequestDB) handleQuery(stm string) []repo.PinRequest {
	var ret []repo.PinRequest
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id string
		var dateInt int
		if err := rows.Scan(&id, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		pr := repo.PinRequest{
			Id:   id,
			Date: time.Unix(int64(dateInt), 0),
		}
		ret = append(ret, pr)
	}
	return ret
}
