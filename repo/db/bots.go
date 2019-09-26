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

func NewBotStore(db *sql.DB, lock *sync.Mutex) repo.BotStore {
	return &BotDB{modelStore{db, lock}}
}

func (c *BotDB) AddOrUpdate(id string, key string, value []byte, botVersion int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into botstore(id, key, value, version, created, updated) values(?,?,?,?,coalesce((select created from botstore where id=? and key=?),?),?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		id,
		key,
		value,
		botVersion,
		id,
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

func (c *BotDB) Get(id string, key string) *pb.BotKV {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery(id, key)
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

func (c *BotDB) Delete(id string, key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from botstore where id=? and key=?", id, key)
	return err
}

func (c *BotDB) handleQuery(id string, key string) []*pb.BotKV {
	list := make([]*pb.BotKV, 0)
	rows, err := c.db.Query("select * from botstore where id=? and key=?", id, key)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, key string
		var value []byte
		var version int32
		var createdInt, updatedInt int64
		if err := rows.Scan(&id, &key, &value, &version, &createdInt, &updatedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		row := c.handleRow(id, key, value, version, createdInt, updatedInt)
		if row != nil {
			list = append(list, row)
		}
	}
	return list
}

func (c *BotDB) handleRow(id string, key string, value []byte, version int32, createdInt int64, updatedInt int64) *pb.BotKV {
	return &pb.BotKV{
		Id:                id,
		Key:               key,
		Value:             value,
		BotReleaseVersion: version,
		Created:           util.ProtoTs(createdInt),
		Updated:           util.ProtoTs(updatedInt),
	}
}
