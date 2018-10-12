package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"time"
)

type AccountDB struct {
	modelStore
}

func NewAccountStore(db *sql.DB, lock *sync.Mutex) repo.AccountStore {
	return &AccountDB{modelStore{db, lock}}
}

func (c *AccountDB) Add(account *repo.Account) error {
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

func (c *AccountDB) Get(id string) *repo.Account {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from accounts where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *AccountDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from accounts;")
	var count int
	row.Scan(&count)
	return count
}

func (c *AccountDB) ListByAddress(address string) []repo.Account {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from accounts where address='" + address + "' order by lastSeen desc;"
	return c.handleQuery(stm)
}

func (c *AccountDB) UpdateLastSeen(id string, date time.Time) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update accounts set lastSeen=? where id=?", int(date.Unix()), id)
	return err
}

func (c *AccountDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from accounts where id=?", id)
	return err
}

func (c *AccountDB) handleQuery(stm string) []repo.Account {
	var ret []repo.Account
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
		block := repo.Account{
			Id:       id,
			Address:  address,
			Created:  time.Unix(int64(createdInt), 0),
			LastSeen: time.Unix(int64(lastSeenInt), 0),
		}
		ret = append(ret, block)
	}
	return ret
}
