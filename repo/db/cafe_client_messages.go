package db

import (
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type CafeClientMessagesDB struct {
	modelStore
}

func NewCafeClientMessageStore(db *sql.DB, lock *sync.Mutex) repo.CafeClientMessageStore {
	return &CafeClientMessagesDB{modelStore{db, lock}}
}

func (c *CafeClientMessagesDB) AddOrUpdate(message *repo.CafeClientMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into cafe_client_messages(id, peerId, clientId, date) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		message.Id,
		message.PeerId,
		message.ClientId,
		int(message.Date.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeClientMessagesDB) ListByClient(clientId string, limit int) []repo.CafeClientMessage {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_client_messages where clientId='" + clientId + "' order by date asc limit " + strconv.Itoa(limit) + ";"
	return c.handleQuery(stm)
}

func (c *CafeClientMessagesDB) CountByClient(clientId string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from cafe_client_messages where clientId='" + clientId + "';")
	var count int
	row.Scan(&count)
	return count
}

func (c *CafeClientMessagesDB) Delete(id string, clientId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_client_messages where id=? and clientId=?", id, clientId)
	return err
}

func (c *CafeClientMessagesDB) DeleteByClient(clientId string, limit int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	sel := "select id from cafe_client_messages where clientId='" + clientId + "' order by date asc limit " + strconv.Itoa(limit)
	query := "delete from cafe_client_messages where id in (" + sel + ");"
	_, err := c.db.Exec(query)
	return err
}

func (c *CafeClientMessagesDB) handleQuery(stm string) []repo.CafeClientMessage {
	var ret []repo.CafeClientMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId, clientId string
		var dateInt int
		if err := rows.Scan(&id, &peerId, &clientId, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		ret = append(ret, repo.CafeClientMessage{
			Id:       id,
			PeerId:   peerId,
			ClientId: clientId,
			Date:     time.Unix(int64(dateInt), 0),
		})
	}
	return ret
}
