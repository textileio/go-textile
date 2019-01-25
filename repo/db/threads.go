package db

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/textileio/textile-go/repo"
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
	stm := `insert into threads(id, key, sk, name, schema, initiator, type, state, head, members, sharing) values(?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		thread.Id,
		thread.Key,
		thread.PrivKey,
		thread.Name,
		thread.Schema,
		thread.Initiator,
		int(thread.Type),
		int(thread.State),
		thread.Head,
		strings.Join(thread.Members, ","),
		int(thread.Sharing),
	)
	if err != nil {
		tx.Rollback()
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

func (c *ThreadDB) GetByKey(key string) *repo.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from threads where key='" + key + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadDB) List() []repo.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.handleQuery("select * from threads;")
}

func (c *ThreadDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from threads;")
	var count int
	row.Scan(&count)
	return count
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

func (c *ThreadDB) handleQuery(stm string) []repo.Thread {
	var ret []repo.Thread
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, key, name, schema, initiator, head, members string
		var skb []byte
		var typeInt, stateInt, sharingInt int
		if err := rows.Scan(&id, &key, &skb, &name, &schema, &initiator, &typeInt, &stateInt, &head, &members, &sharingInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		mlist := make([]string, 0)
		for _, m := range strings.Split(members, ",") {
			if m != "" {
				mlist = append(mlist, m)
			}
		}
		ret = append(ret, repo.Thread{
			Id:        id,
			Key:       key,
			PrivKey:   skb,
			Name:      name,
			Schema:    schema,
			Initiator: initiator,
			Type:      repo.ThreadType(typeInt),
			Sharing:   repo.ThreadSharing(sharingInt),
			Members:   mlist,
			State:     repo.ThreadState(stateInt),
			Head:      head,
		})
	}
	return ret
}
