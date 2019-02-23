package db

import (
	"database/sql"
	"strconv"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
)

type BlockMessageDB struct {
	modelStore
}

func NewBlockMessageStore(db *sql.DB, lock *sync.Mutex) repo.BlockMessageStore {
	return &BlockMessageDB{modelStore{db, lock}}
}

func (c *BlockMessageDB) Add(msg *pb.BlockMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into block_messages(id, peerId, envelope, date) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	env, err := proto.Marshal(msg.Env)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		msg.Id,
		msg.Peer,
		env,
		util.ProtoNanos(msg.Date),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *BlockMessageDB) List(offset string, limit int) []pb.BlockMessage {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from block_messages where date>(select date from block_messages where id='" + offset + "') order by date asc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from block_messages order by date asc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *BlockMessageDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from block_messages where id=?", id)
	return err
}

func (c *BlockMessageDB) handleQuery(stm string) []pb.BlockMessage {
	var list []pb.BlockMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId string
		var dateInt int64
		var envelopeb []byte
		if err := rows.Scan(&id, &peerId, &envelopeb, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		env := new(pb.Envelope)
		if err := proto.Unmarshal(envelopeb, env); err != nil {
			log.Errorf("error unmarshaling envelope: %s", err)
			continue
		}

		list = append(list, pb.BlockMessage{
			Id:   id,
			Peer: peerId,
			Env:  env,
			Date: util.ProtoTs(dateInt),
		})
	}
	return list
}
