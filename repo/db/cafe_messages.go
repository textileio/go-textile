package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type CafeMessageDB struct {
	modelStore
}

func NewCafeMessageStore(db *sql.DB, lock *sync.Mutex) repo.CafeMessageStore {
	return &CafeMessageDB{modelStore{db, lock}}
}

func (c *CafeMessageDB) Add(req *repo.CafeMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_messages(id, peerId, date) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		req.Id,
		req.PeerId,
		int(req.Date.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeMessageDB) List(offset string, limit int) []repo.CafeMessage {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from cafe_messages where date>(select date from cafe_messages where id='" + offset + "') order by date asc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from cafe_messages order by date asc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *CafeMessageDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_messages where id=?", id)
	return err
}

func (c *CafeMessageDB) handleQuery(stm string) []repo.CafeMessage {
	var ret []repo.CafeMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId string
		var dateInt int
		if err := rows.Scan(&id, &peerId, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		ret = append(ret, repo.CafeMessage{
			Id:     id,
			PeerId: peerId,
			Date:   time.Unix(int64(dateInt), 0),
		})
	}
	return ret
}
