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
	stm := `insert into peers(row, id, thread, pk) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		peer.Row,
		peer.Id,
		peer.ThreadId,
		peer.PubKey,
	)
	if err != nil {
		tx.Rollback()
		log.Errorf("error in db exec: %s", err)
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

func (c *PeerDB) List(offset string, limit int, query string) []repo.Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		q := ""
		if query != "" {
			q = query + " and "
		}
		stm = "select * from peers where " + q + "row<(select row from peers where row='" + offset + "') order by row desc limit " + strconv.Itoa(limit) + " ;"
	} else {
		q := ""
		if query != "" {
			q = "where " + query + " "
		}
		stm = "select * from peers " + q + "order by row desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *PeerDB) Delete(row string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from peers where row=?", row)
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
		var row, id, thread string
		var pk []byte
		if err := rows.Scan(&row, &id, &thread, &pk); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		block := repo.Peer{
			Row:      row,
			Id:       id,
			ThreadId: thread,
			PubKey:   pk,
		}
		ret = append(ret, block)
	}
	return ret
}
