package db

import (
	"database/sql"
	"strconv"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeClientMessagesDB struct {
	modelStore
}

func NewCafeClientMessageStore(db *sql.DB, lock *sync.Mutex) repo.CafeClientMessageStore {
	return &CafeClientMessagesDB{modelStore{db, lock}}
}

func (c *CafeClientMessagesDB) AddOrUpdate(message *pb.CafeClientMessage) error {
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
		message.Peer,
		message.Client,
		util.ProtoNanos(message.Date),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *CafeClientMessagesDB) ListByClient(clientId string, limit int) []pb.CafeClientMessage {
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
	_ = row.Scan(&count)
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
	sel := "select id from cafe_client_messages where clientId='" + clientId + "'"
	if limit > 0 {
		sel += " order by date asc limit " + strconv.Itoa(limit)
	}
	query := "delete from cafe_client_messages where id in (" + sel + ");"
	_, err := c.db.Exec(query)
	return err
}

func (c *CafeClientMessagesDB) handleQuery(stm string) []pb.CafeClientMessage {
	var list []pb.CafeClientMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId, clientId string
		var dateInt int64
		if err := rows.Scan(&id, &peerId, &clientId, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		list = append(list, pb.CafeClientMessage{
			Id:     id,
			Peer:   peerId,
			Client: clientId,
			Date:   util.ProtoTs(dateInt),
		})
	}
	return list
}
