package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type AccountPeerDB struct {
	modelStore
}

func NewAccountPeerStore(db *sql.DB, lock *sync.Mutex) repo.AccountPeerStore {
	return &AccountPeerDB{modelStore{db, lock}}
}

func (c *AccountPeerDB) Add(peer *repo.AccountPeer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into account_peers(id, name) values(?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		peer.Id,
		peer.Name,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *AccountPeerDB) Get(id string) *repo.AccountPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from account_peers where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *AccountPeerDB) List(query string) []repo.AccountPeer {
	c.lock.Lock()
	defer c.lock.Unlock()
	var q string
	if query != "" {
		q = " where " + query
	}
	return c.handleQuery("select * from account_peers" + q + ";")
}

func (c *AccountPeerDB) Count(query string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	var q string
	if query != "" {
		q = " where " + query
	}
	row := c.db.QueryRow("select Count(*) from account_peers" + q + ";")
	var count int
	row.Scan(&count)
	return count
}

func (c *AccountPeerDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from account_peers where id=?", id)
	return err
}

func (c *AccountPeerDB) handleQuery(stm string) []repo.AccountPeer {
	var ret []repo.AccountPeer
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		peer := repo.AccountPeer{
			Id:   id,
			Name: name,
		}
		ret = append(ret, peer)
	}
	return ret
}
