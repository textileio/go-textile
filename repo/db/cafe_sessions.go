package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
	"time"
)

type CafeSessionDB struct {
	modelStore
}

func NewCafeSessionStore(db *sql.DB, lock *sync.Mutex) repo.CafeSessionStore {
	return &CafeSessionDB{modelStore{db, lock}}
}

func (c *CafeSessionDB) AddOrUpdate(session *repo.CafeSession) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into sessions(cafeId, access, refresh, expiry) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		session.CafeId,
		session.Access,
		session.Refresh,
		int(session.Expiry.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeSessionDB) Get(cafeId string) *repo.CafeSession {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from sessions where cafeId='" + cafeId + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeSessionDB) List() []repo.CafeSession {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from sessions order by expiry desc;"
	return c.handleQuery(stm)
}

func (c *CafeSessionDB) Delete(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from sessions where cafeId=?", cafeId)
	return err
}

func (c *CafeSessionDB) handleQuery(stm string) []repo.CafeSession {
	var ret []repo.CafeSession
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var cafeId, access, refresh string
		var expiryInt int
		if err := rows.Scan(&cafeId, &access, &refresh, &expiryInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		session := repo.CafeSession{
			CafeId:  cafeId,
			Access:  access,
			Refresh: refresh,
			Expiry:  time.Unix(int64(expiryInt), 0),
		}
		ret = append(ret, session)
	}
	return ret
}
