package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type CafeMessagesDB struct {
	modelStore
}

func NewCafeMessageStore(db *sql.DB, lock *sync.Mutex) repo.CafeMessagesStore {
	return &CafeMessagesDB{modelStore{db, lock}}
}

func (c *CafeMessagesDB) AddOrUpdate(message *repo.CafeMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into cafe_messages(id, accountId, date, read) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		message.Id,
		message.AccountId,
		int(message.Date.Unix()),
		false,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeMessagesDB) ListByAccount(accountId string, offset string, limit int) []repo.CafeMessage {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from cafe_messages where accountId='" + accountId + "' and date<(select date from cafe_messages where id='" + offset + "') order by date desc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from cafe_messages where accountId='" + accountId + "' order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *CafeMessagesDB) Read(id string, accountId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update cafe_messages set read=1 where id=? and accountId=?", id, accountId)
	return err
}

func (c *CafeMessagesDB) Delete(id string, accountId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_messages where id=? and accountId=?", id, accountId)
	return err
}

func (c *CafeMessagesDB) DeleteByAccount(accountId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_messages where accountId=?", accountId)
	return err
}

func (c *CafeMessagesDB) handleQuery(stm string) []repo.CafeMessage {
	var ret []repo.CafeMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, accountId string
		var dateInt, readInt int
		if err := rows.Scan(&id, &accountId, &dateInt, &readInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		read := false
		if readInt == 1 {
			read = true
		}
		message := repo.CafeMessage{
			Id:        id,
			AccountId: accountId,
			Date:      time.Unix(int64(dateInt), 0),
			Read:      read,
		}
		ret = append(ret, message)
	}
	return ret
}
