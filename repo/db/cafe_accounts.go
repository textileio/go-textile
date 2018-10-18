package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"time"
)

type CafeAccountDB struct {
	modelStore
}

func NewCafeAccountStore(db *sql.DB, lock *sync.Mutex) repo.CafeAccountStore {
	return &CafeAccountDB{modelStore{db, lock}}
}

func (c *CafeAccountDB) Add(account *repo.CafeAccount) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into accounts(id, address, created, lastSeen) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		account.Id,
		account.Address,
		int(account.Created.Unix()),
		int(account.LastSeen.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeAccountDB) Get(id string) *repo.CafeAccount {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from accounts where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeAccountDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from accounts;")
	var count int
	row.Scan(&count)
	return count
}

func (c *CafeAccountDB) List() []repo.CafeAccount {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from accounts order by lastSeen desc;"
	return c.handleQuery(stm)
}

func (c *CafeAccountDB) ListByAddress(address string) []repo.CafeAccount {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from accounts where address='" + address + "' order by lastSeen desc;"
	return c.handleQuery(stm)
}

func (c *CafeAccountDB) UpdateLastSeen(id string, date time.Time) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update accounts set lastSeen=? where id=?", int(date.Unix()), id)
	return err
}

func (c *CafeAccountDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from accounts where id=?", id)
	return err
}

func (c *CafeAccountDB) handleQuery(stm string) []repo.CafeAccount {
	var ret []repo.CafeAccount
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, address string
		var createdInt, lastSeenInt int
		if err := rows.Scan(&id, &address, &createdInt, &lastSeenInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		accnt := repo.CafeAccount{
			Id:       id,
			Address:  address,
			Created:  time.Unix(int64(createdInt), 0),
			LastSeen: time.Unix(int64(lastSeenInt), 0),
		}
		ret = append(ret, accnt)
	}
	return ret
}
