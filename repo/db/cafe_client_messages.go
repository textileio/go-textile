package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type CafeClientMessagesDB struct {
	modelStore
}

func NewCafeClientMessageStore(db *sql.DB, lock *sync.Mutex) repo.CafeClientMessageStore {
	return &CafeClientMessagesDB{modelStore{db, lock}}
}

func (c *CafeClientMessagesDB) AddOrUpdate(message *repo.CafeMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into cafe_client_messages(id, clientId, date, read) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		message.Id,
		message.ClientId,
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

func (c *CafeClientMessagesDB) ListByClient(clientId string, offset string, limit int) []repo.CafeMessage {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from cafe_client_messages where clientId='" + clientId + "' and date<(select date from cafe_client_messages where id='" + offset + "') order by date desc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from cafe_client_messages where clientId='" + clientId + "' order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *CafeClientMessagesDB) Read(id string, clientId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update cafe_client_messages set read=1 where id=? and clientId=?", id, clientId)
	return err
}

func (c *CafeClientMessagesDB) Delete(id string, clientId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_client_messages where id=? and clientId=?", id, clientId)
	return err
}

func (c *CafeClientMessagesDB) DeleteByClient(clientId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_client_messages where clientId=?", clientId)
	return err
}

func (c *CafeClientMessagesDB) handleQuery(stm string) []repo.CafeMessage {
	var ret []repo.CafeMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, clientId string
		var dateInt, readInt int
		if err := rows.Scan(&id, &clientId, &dateInt, &readInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		read := false
		if readInt == 1 {
			read = true
		}
		message := repo.CafeMessage{
			Id:       id,
			ClientId: clientId,
			Date:     time.Unix(int64(dateInt), 0),
			Read:     read,
		}
		ret = append(ret, message)
	}
	return ret
}
