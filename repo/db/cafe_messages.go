package db

import (
	"database/sql"
	"strconv"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeMessageDB struct {
	modelStore
}

func NewCafeMessageStore(db *sql.DB, lock *sync.Mutex) repo.CafeMessageStore {
	return &CafeMessageDB{modelStore{db, lock}}
}

func (c *CafeMessageDB) Add(req *pb.CafeMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_messages(id, peerId, date, attempts) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		req.Id,
		req.Peer,
		util.ProtoNanos(req.Date),
		req.Attempts,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *CafeMessageDB) List(offset string, limit int) []pb.CafeMessage {
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

func (c *CafeMessageDB) AddAttempt(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update cafe_messages set attempts=attempts+1 where id=?", id)
	return err
}

func (c *CafeMessageDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_messages where id=?", id)
	return err
}

func (c *CafeMessageDB) handleQuery(stm string) []pb.CafeMessage {
	var list []pb.CafeMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId string
		var dateInt int64
		var attempts int
		if err := rows.Scan(&id, &peerId, &dateInt, &attempts); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		list = append(list, pb.CafeMessage{
			Id:       id,
			Peer:     peerId,
			Date:     util.ProtoTs(dateInt),
			Attempts: int32(attempts),
		})
	}
	return list
}
