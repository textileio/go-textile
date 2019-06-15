package db

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type ThreadDB struct {
	modelStore
}

func NewThreadStore(db *sql.DB, lock *sync.Mutex) repo.ThreadStore {
	return &ThreadDB{modelStore{db, lock}}
}

func (c *ThreadDB) Add(thread *pb.Thread) error {
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
		thread.Sk,
		thread.Name,
		thread.Schema,
		thread.Initiator,
		int(thread.Type),
		int(thread.State),
		thread.Head,
		strings.Join(thread.Whitelist, ","),
		int(thread.Sharing),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *ThreadDB) Get(id string) *pb.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from threads where id='" + id + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
}

func (c *ThreadDB) GetByKey(key string) *pb.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from threads where key='" + key + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
}

func (c *ThreadDB) List() *pb.ThreadList {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.handleQuery("select * from threads;")
}

func (c *ThreadDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from threads;")
	var count int
	_ = row.Scan(&count)
	return count
}

func (c *ThreadDB) UpdateHead(id string, heads []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update threads set head=? where id=?", strings.Join(heads, ","), id)
	return err
}

func (c *ThreadDB) UpdateName(id string, name string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update threads set name=? where id=?", name, id)
	return err
}

func (c *ThreadDB) UpdateSchema(id string, hash string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update threads set schema=? where id=?", hash, id)
	return err
}

func (c *ThreadDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from threads where id=?", id)
	return err
}

func (c *ThreadDB) handleQuery(stm string) *pb.ThreadList {
	list := &pb.ThreadList{Items: make([]*pb.Thread, 0)}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, key, name, schema, initiator, head, whitelist string
		var skb []byte
		var typeInt, stateInt, sharingInt int
		err := rows.Scan(&id, &key, &skb, &name, &schema, &initiator, &typeInt, &stateInt, &head, &whitelist, &sharingInt)
		if err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		list.Items = append(list.Items, &pb.Thread{
			Id:        id,
			Key:       key,
			Sk:        skb,
			Name:      name,
			Schema:    schema,
			Initiator: initiator,
			Type:      pb.Thread_Type(typeInt),
			Sharing:   pb.Thread_Sharing(sharingInt),
			Whitelist: util.SplitString(whitelist, ","),
			State:     pb.Thread_State(stateInt),
			Head:      head,
		})
	}
	return list
}
