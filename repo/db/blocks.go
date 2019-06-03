package db

import (
	"database/sql"
	"strconv"
	"strings"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
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
	stmt, err := tx.Prepare(`
        INSERT INTO blocks(
    	    id, threadId, authorId, type, date, parents, target, body, data, status, attempts
        ) VALUES (?,?,?,?,?,?,?,?,?,?,?)
    `)
	if err != nil {
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
		block.Data,
		int32(block.Status),
		block.Attempts,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *BlockDB) Replace(block *pb.Block) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        REPLACE INTO blocks(
    	    id, threadId, authorId, type, date, parents, target, body, data, status, attempts
        ) VALUES (?,?,?,?,?,?,?,?,?,?,coalesce((SELECT attempts FROM blocks WHERE id=?),?))
    `)
	if err != nil {
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
		block.Data,
		int32(block.Status),
		block.Id,
		block.Attempts,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *BlockDB) Get(id string) *pb.Block {
	c.lock.Lock()
	defer c.lock.Unlock()

	res := c.handleQuery("SELECT * FROM blocks WHERE id='" + id + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
}

func (c *BlockDB) List(offset string, limit int, query string) *pb.BlockList {
	c.lock.Lock()
	defer c.lock.Unlock()

	limits := strconv.Itoa(limit)
	stm := "SELECT * FROM blocks"
	if offset != "" {
		if query != "" {
			query += " and "
		}
		stm += " WHERE " + query + "(date<(SELECT date FROM blocks WHERE id='" + offset + "'))"
	} else if query != "" {
		stm += " WHERE " + query
	}
	stm += " ORDER BY date DESC LIMIT " + limits + ";"

	return c.handleQuery(stm)
}

func (c *BlockDB) Count(query string) int {
	c.lock.Lock()
	defer c.lock.Unlock()

	stm := "SELECT COUNT(*) FROM blocks"
	if query != "" {
		stm += " WHERE " + query
	}
	stm += ";"

	row := c.db.QueryRow(stm)
	var count int
	_ = row.Scan(&count)

	return count
}

func (c *BlockDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("DELETE FROM blocks WHERE id=?", id)
	return err
}

func (c *BlockDB) DeleteByThread(threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("DELETE FROM blocks WHERE threadId=?", threadId)
	return err
}

func (c *BlockDB) AddAttempt(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("UPDATE blocks SET attempts=attempts+1 WHERE id=?", id)
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
		var id, threadId, authorId, parents, target, body, data string
		var typeInt, statusInt, attempts int
		var dateInt int64

		err = rows.Scan(&id, &threadId, &authorId, &typeInt, &dateInt, &parents, &target, &body, &data, &statusInt, &attempts)
		if err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		list.Items = append(list.Items, &pb.Block{
			Id:       id,
			Thread:   threadId,
			Author:   authorId,
			Type:     pb.Block_BlockType(typeInt),
			Date:     util.ProtoTs(dateInt),
			Parents:  util.SplitString(parents, ","),
			Target:   target,
			Body:     body,
			Data:     data,
			Status:   pb.Block_BlockStatus(statusInt),
			Attempts: int32(attempts),
		})
	}

	return list
}
