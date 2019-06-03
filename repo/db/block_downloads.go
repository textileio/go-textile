package db

import (
	"database/sql"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type BlockDownloadDB struct {
	modelStore
}

func NewBlockDownloadStore(db *sql.DB, lock *sync.Mutex) repo.BlockDownloadStore {
	return &BlockDownloadDB{modelStore{db, lock}}
}

func (c *BlockDownloadDB) Add(dl *pb.BlockDownload) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO block_downloads(id, threadId, parents, target, date, attempts) VALUES (?,?,?,?,?,?)`)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		dl.Id,
		dl.Thread,
		strings.Join(dl.Parents, ","),
		dl.Target,
		time.Now().UnixNano(),
		dl.Attempts,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()

	return nil
}

func (c *BlockDownloadDB) List(offset string, limit int) []pb.BlockDownload {
	c.lock.Lock()
	defer c.lock.Unlock()

	stm := "SELECT * FROM block_downloads"
	if offset != "" {
		stm += " WHERE date>(SELECT date FROM block_downloads WHERE id='" + offset + "')"
	}
	stm += " ORDER BY date ASC LIMIT " + strconv.Itoa(limit) + ";"

	return c.handleQuery(stm)
}

func (c *BlockDownloadDB) AddAttempt(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("UPDATE block_downloads SET attempts=attempts+1 WHERE id=?", id)
	return err
}

func (c *BlockDownloadDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("DELETE FROM block_downloads WHERE id=?", id)
	return err
}

func (c *BlockDownloadDB) handleQuery(stm string) []pb.BlockDownload {
	var list []pb.BlockDownload

	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, thread, parents, target string
		var attempts int
		var dateInt int64

		if err := rows.Scan(&id, &thread, &parents, &target, &dateInt, &attempts); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		list = append(list, pb.BlockDownload{
			Id:       id,
			Thread:   thread,
			Parents:  util.SplitString("parents", ","),
			Target:   target,
			Date:     util.ProtoTs(dateInt),
			Attempts: int32(attempts),
		})
	}

	return list
}
