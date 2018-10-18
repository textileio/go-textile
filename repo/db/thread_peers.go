package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
)

type ThreadPeerDB struct {
	modelStore
}

func NewThreadPeerStore(db *sql.DB, lock *sync.Mutex) repo.ThreadPeerStore {
	return &ThreadPeerDB{modelStore{db, lock}}
}

func (c *ThreadPeerDB) Add(peer *repo.ThreadPeer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into thread_peers(row, id, pk, threadId) values(?,?,?,?)`
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

func (c *ThreadPeerDB) Get(row string) *repo.ThreadPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from thread_peers where row='" + row + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadPeerDB) GetById(id string) *repo.ThreadPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from thread_peers where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadPeerDB) List(limit int, query string) []repo.ThreadPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm, q string
	if query != "" {
		q = "where " + query + " "
	}
	stm = "select * from thread_peers " + q + "limit " + strconv.Itoa(limit) + ";"
	return c.handleQuery(stm)
}

func (c *ThreadPeerDB) Count(query string, distinct bool) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm, q string
	if query != "" {
		q = " where " + query
	}
	if distinct {
		stm = "select Count(distinct id) from thread_peers" + q + ";"
	} else {
		stm = "select Count(*) from thread_peers" + q + ";"
	}
	row := c.db.QueryRow(stm)
	var count int
	row.Scan(&count)
	return count
}

func (c *ThreadPeerDB) Delete(id string, threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_peers where id=? and threadId=?", id, threadId)
	return err
}

func (c *ThreadPeerDB) DeleteByThreadId(threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_peers where threadId=?", threadId)
	return err
}

func (c *ThreadPeerDB) handleQuery(stm string) []repo.ThreadPeer {
	var ret []repo.ThreadPeer
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
		block := repo.ThreadPeer{
			Row:      row,
			Id:       id,
			PubKey:   pk,
			ThreadId: threadId,
		}
		ret = append(ret, block)
	}
	return ret
}
