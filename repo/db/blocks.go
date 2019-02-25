package db

import (
	"database/sql"
	"strconv"
	"strings"
	"sync"

	"github.com/textileio/textile-go/util"

	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

type BlockDB struct {
	modelStore
}

func NewBlockStore(db *sql.DB, lock *sync.Mutex) repo.BlockStore {
	return &BlockDB{modelStore{db, lock}}
}

func (c *BlockDB) Add(block *pb.Block) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into blocks(id, threadId, authorId, type, date, parents, target, body) values(?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		block.Id,
		block.Thread,
		block.Author,
		int(block.Type),
		util.ProtoNanos(block.Date),
		strings.Join(block.Parents, ","),
		block.Target,
		block.Body,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *BlockDB) Get(id string) *pb.Block {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from blocks where id='" + id + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
}

func (c *BlockDB) List(offset string, limit int, query string) *pb.BlockList {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm, q string
	if offset != "" {
		if query != "" {
			q = query + " and "
		}
		stm = "select * from blocks where " + q + "(date<(select date from blocks where id='" + offset + "')) order by date desc limit " + strconv.Itoa(limit) + ";"
	} else {
		if query != "" {
			q = "where " + query + " "
		}
		stm = "select * from blocks " + q + "order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *BlockDB) Count(query string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	var q string
	if query != "" {
		q = " where " + query
	}
	row := c.db.QueryRow("select Count(*) from blocks" + q + ";")
	var count int
	row.Scan(&count)
	return count
}

func (c *BlockDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from blocks where id=?", id)
	return err
}

func (c *BlockDB) DeleteByThread(threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from blocks where threadId=?", threadId)
	return err
}

func (c *BlockDB) handleQuery(stm string) *pb.BlockList {
	list := &pb.BlockList{Items: make([]*pb.Block, 0)}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, threadId, authorId, parents, target, body string
		var typeInt int
		var dateInt int64
		if err := rows.Scan(&id, &threadId, &authorId, &typeInt, &dateInt, &parents, &target, &body); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		list.Items = append(list.Items, &pb.Block{
			Id:      id,
			Thread:  threadId,
			Author:  authorId,
			Type:    pb.Block_BlockType(typeInt),
			Date:    util.ProtoTs(dateInt),
			Parents: util.SplitString(parents, ","),
			Target:  target,
			Body:    body,
		})
	}
	return list
}
