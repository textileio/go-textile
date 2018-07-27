package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BlockDB struct {
	modelStore
}

func NewBlockStore(db *sql.DB, lock *sync.Mutex) repo.BlockStore {
	return &BlockDB{modelStore{db, lock}}
}

func (c *BlockDB) Add(block *repo.Block) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into blocks(id, date, parents, threadId, authorPk, type, dataId, dataKeyCipher, dataCaptionCipher) values(?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		block.Id,
		int(block.Date.Unix()),
		strings.Join(block.Parents, ","),
		block.ThreadId,
		block.AuthorPk,
		int(block.Type),
		block.DataId,
		block.DataKeyCipher,
		block.DataCaptionCipher,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *BlockDB) Get(id string) *repo.Block {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from blocks where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *BlockDB) GetByDataId(dataId string) *repo.Block {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from blocks where dataId='" + dataId + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *BlockDB) List(offset string, limit int, query string) []repo.Block {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		q := ""
		if query != "" {
			q = query + " and "
		}
		stm = "select * from blocks where " + q + "date<(select date from blocks where id='" + offset + "') order by date desc limit " + strconv.Itoa(limit) + " ;"
	} else {
		q := ""
		if query != "" {
			q = "where " + query + " "
		}
		stm = "select * from blocks " + q + "order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *BlockDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from blocks where id=?", id)
	return err
}

func (c *BlockDB) DeleteByThreadId(threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from blocks where threadId=?", threadId)
	return err
}

func (c *BlockDB) handleQuery(stm string) []repo.Block {
	var ret []repo.Block
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, parents, threadId, authorPk, dataId string
		var dateInt, typeInt int
		var dataKeyCipher, dataCaptionCipher []byte
		if err := rows.Scan(&id, &dateInt, &parents, &threadId, &authorPk, &typeInt, &dataId, &dataKeyCipher, &dataCaptionCipher); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		block := repo.Block{
			Id:                id,
			Date:              time.Unix(int64(dateInt), 0),
			Parents:           strings.Split(parents, ","),
			ThreadId:          threadId,
			AuthorPk:          authorPk,
			Type:              repo.BlockType(typeInt),
			DataId:            dataId,
			DataKeyCipher:     dataKeyCipher,
			DataCaptionCipher: dataCaptionCipher,
		}
		ret = append(ret, block)
	}
	return ret
}
