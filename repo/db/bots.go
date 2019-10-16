package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type BotDB struct {
	modelStore
}

func NewBotstore(db *sql.DB, lock *sync.Mutex) repo.Botstore {
	return &BotDB{modelStore{db, lock}}
}

// AddOrUpdate Bot KV store adds namespace all bot requests by their key
func (c *BotDB) AddOrUpdate(key string, value []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into bots_store(id, value, created, updated) values(?,?,coalesce((select created from bots_store where id=?),?),?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		key,
		value,
		key,
		time.Now().UnixNano(),
		time.Now().UnixNano(),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

// Get Bot KV store gets namespace all bot requests by their key
func (c *BotDB) Get(key string) *pb.BotKV {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery(key)
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

// Delete Bot KV store deletes namespace all bot requests by their key
func (c *BotDB) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from bots_store where id=?", key)
	return err
}

func (c *BotDB) handleQuery(key string) []*pb.BotKV {
	list := make([]*pb.BotKV, 0)
	rows, err := c.db.Query("select * from bots_store where id=?", key)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id string
		var value []byte
		var createdInt, updatedInt int64
		if err := rows.Scan(&id, &value, &createdInt, &updatedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		row := c.handleRow(id, value, createdInt, updatedInt)
		if row != nil {
			list = append(list, row)
		}
	}
	return list
}

func (c *BotDB) handleRow(key string, value []byte, createdInt int64, updatedInt int64) *pb.BotKV {
	return &pb.BotKV{
		Key:     key,
		Value:   value,
		Created: util.ProtoTs(createdInt),
		Updated: util.ProtoTs(updatedInt),
	}
}
