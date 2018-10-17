package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
)

type PeerDB struct {
	modelStore
}

func NewPeerStore(db *sql.DB, lock *sync.Mutex) repo.PeerStore {
	return &PeerDB{modelStore{db, lock}}
}

func (c *PeerDB) Add(peer *repo.Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into peers(row, id, pk, threadId) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		peer.Row,
		peer.Id,
		peer.PubKey,
		peer.ThreadId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *PeerDB) Get(row string) *repo.Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from peers where row='" + row + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *PeerDB) GetById(id string) *repo.Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from peers where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *PeerDB) List(limit int, query string) []repo.Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm, q string
	if query != "" {
		q = "where " + query + " "
	}
	stm = "select * from peers " + q + "limit " + strconv.Itoa(limit) + ";"
	return c.handleQuery(stm)
}

func (c *PeerDB) Count(query string, distinct bool) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm, q string
	if query != "" {
		q = " where " + query
	}
	if distinct {
		stm = "select Count(distinct id) from peers" + q + ";"
	} else {
		stm = "select Count(*) from peers" + q + ";"
	}
	row := c.db.QueryRow(stm)
	var count int
	row.Scan(&count)
	return count
}

func (c *PeerDB) Delete(id string, threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from peers where id=? and threadId=?", id, threadId)
	return err
}

func (c *PeerDB) DeleteByThreadId(threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from peers where threadId=?", threadId)
	return err
}

func (c *PeerDB) handleQuery(stm string) []repo.Peer {
	var ret []repo.Peer
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var row, id, threadId string
		var pk []byte
		if err := rows.Scan(&row, &id, &pk, &threadId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		block := repo.Peer{
			Row:      row,
			Id:       id,
			PubKey:   pk,
			ThreadId: threadId,
		}
		ret = append(ret, block)
	}
	return ret
}
