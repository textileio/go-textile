package db

import (
	"database/sql"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

type ThreadPeerDB struct {
	modelStore
}

func NewThreadPeerStore(db *sql.DB, lock *sync.Mutex) repo.ThreadPeerStore {
	return &ThreadPeerDB{modelStore{db, lock}}
}

func (c *ThreadPeerDB) Add(peer *pb.ThreadPeer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into thread_peers(id, threadId, welcomed) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		peer.Id,
		peer.Thread,
		false,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *ThreadPeerDB) List() []pb.ThreadPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from thread_peers;"
	return c.handleQuery(stm)
}

func (c *ThreadPeerDB) ListById(id string) []pb.ThreadPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from thread_peers where id='" + id + "';"
	return c.handleQuery(stm)
}

func (c *ThreadPeerDB) ListByThread(threadId string) []pb.ThreadPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from thread_peers where threadId='" + threadId + "';"
	return c.handleQuery(stm)
}

func (c *ThreadPeerDB) ListUnwelcomedByThread(threadId string) []pb.ThreadPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from thread_peers where threadId='" + threadId + "' and welcomed=0;"
	return c.handleQuery(stm)
}

func (c *ThreadPeerDB) WelcomeByThread(threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update thread_peers set welcomed=1 where threadId=?", threadId)
	return err
}

func (c *ThreadPeerDB) Count(distinct bool) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if distinct {
		stm = "select Count(distinct id) from thread_peers;"
	} else {
		stm = "select Count(*) from thread_peers;"
	}
	row := c.db.QueryRow(stm)
	var count int
	_ = row.Scan(&count)
	return count
}

func (c *ThreadPeerDB) Delete(id string, threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_peers where id=? and threadId=?", id, threadId)
	return err
}

func (c *ThreadPeerDB) DeleteById(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_peers where id=?", id)
	return err
}

func (c *ThreadPeerDB) DeleteByThread(threadId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_peers where threadId=?", threadId)
	return err
}

func (c *ThreadPeerDB) handleQuery(stm string) []pb.ThreadPeer {
	var list []pb.ThreadPeer
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, threadId string
		var welcomedInt int
		if err := rows.Scan(&id, &threadId, &welcomedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		welcomed := false
		if welcomedInt == 1 {
			welcomed = true
		}
		list = append(list, pb.ThreadPeer{
			Id:       id,
			Thread:   threadId,
			Welcomed: welcomed,
		})
	}
	return list
}
