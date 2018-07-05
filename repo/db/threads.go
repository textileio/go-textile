package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type ThreadDB struct {
	modelStore
}

func NewThreadStore(db *sql.DB, lock *sync.Mutex) repo.ThreadStore {
	return &ThreadDB{modelStore{db, lock}}
}

func (c *ThreadDB) Add(thread *repo.Thread) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into threads(id, name, sk, head) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		thread.Id,
		thread.Name,
		thread.PrivKey,
		thread.Head,
	)
	if err != nil {
		tx.Rollback()
		log.Errorf("error in db exec: %s", err)
		return err
	}
	tx.Commit()
	return nil
}

func (c *ThreadDB) Get(id string) *repo.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from threads where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadDB) GetByName(name string) *repo.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from threads where name='" + name + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadDB) List(query string) []repo.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := ""
	if query != "" {
		q = " where " + query
	}
	return c.handleQuery("select * from threads" + q + ";")
}

func (c *ThreadDB) UpdateHead(id string, head string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update threads set head=? where id=?", head, id)
	return err
}

func (c *ThreadDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from threads where id=?", id)
	return err
}

func (c *ThreadDB) DeleteByName(name string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from threads where name=?", name)
	return err
}

func (c *ThreadDB) handleQuery(stm string) []repo.Thread {
	var ret []repo.Thread
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, name, head string
		var skb []byte
		if err := rows.Scan(&id, &name, &skb, &head); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		thread := repo.Thread{
			Id:      id,
			Name:    name,
			PrivKey: skb,
			Head:    head,
		}
		ret = append(ret, thread)
	}
	return ret
}
