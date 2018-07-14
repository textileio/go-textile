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
	stm := `insert into blocks(id, target, parents, key, pk, ppk, type, date) values(?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		block.Id,
		block.Target,
		strings.Join(block.Parents, ","),
		block.TargetKey,
		block.ThreadPubKey,
		block.PeerPubKey,
		int(block.Type),
		int(block.Date.Unix()),
	)
	if err != nil {
		tx.Rollback()
		log.Errorf("error in db exec: %s", err)
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

func (c *BlockDB) GetByTarget(target string) *repo.Block {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from blocks where target='" + target + "';")
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

func (c *BlockDB) handleQuery(stm string) []repo.Block {
	var ret []repo.Block
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, target, parents, pk, ppk string
		var key []byte
		var typeInt, dateInt int
		if err := rows.Scan(&id, &target, &parents, &key, &pk, &ppk, &typeInt, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		block := repo.Block{
			Id:           id,
			Target:       target,
			Parents:      strings.Split(parents, ","),
			TargetKey:    key,
			ThreadPubKey: pk,
			PeerPubKey:   ppk,
			Type:         repo.BlockType(typeInt),
			Date:         time.Unix(int64(dateInt), 0),
		}
		ret = append(ret, block)
	}
	return ret
}
