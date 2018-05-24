package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"sync"
)

type ThreadDB struct {
	modelStore
}

func NewThreadStore(db *sql.DB, lock *sync.Mutex) repo.ThreadStore {
	return &ThreadDB{modelStore{db, lock}}
}

func (c *ThreadDB) Add(thread *wallet.Thread) error {
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
	skb, err := thread.PrivKey.Bytes()
	if err != nil {
		log.Errorf("error getting key bytes: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		thread.Id,
		skb,
		thread.Name,
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

func (c *ThreadDB) Get(id string) *wallet.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from threads where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadDB) GetByName(name string) *wallet.Thread {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from threads where name='" + name + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadDB) List(query string) []wallet.Thread {
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

func (c *ThreadDB) handleQuery(stm string) []wallet.Thread {
	var ret []wallet.Thread
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
		sk, err := libp2p.UnmarshalPrivateKey(skb)
		if err != nil {
			log.Errorf("error unmarshaling private key: %s", err)
			continue
		}
		album := wallet.Thread{
			Id:      id,
			Name:    name,
			PrivKey: sk,
			Head:    head,
		}
		ret = append(ret, album)
	}
	return ret
}
