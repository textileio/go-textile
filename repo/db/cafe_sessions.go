package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strings"
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
	stm := `insert or replace into cafe_sessions(cafeId, access, refresh, expiry, httpAddr, swarmAddrs) values(?,?,?,?,?,?)`
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
		session.HttpAddr,
		strings.Join(session.SwarmAddrs, ","),
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
	ret := c.handleQuery("select * from cafe_sessions where cafeId='" + cafeId + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeSessionDB) List() []repo.CafeSession {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_sessions order by expiry desc;"
	return c.handleQuery(stm)
}

func (c *CafeSessionDB) Delete(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_sessions where cafeId=?", cafeId)
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
		var cafeId, access, refresh, httpAddr, swarmAddrs string
		var expiryInt int
		if err := rows.Scan(&cafeId, &access, &refresh, &expiryInt, &httpAddr, &swarmAddrs); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		slist := make([]string, 0)
		for _, p := range strings.Split(swarmAddrs, ",") {
			if p != "" {
				slist = append(slist, p)
			}
		}
		ret = append(ret, repo.CafeSession{
			CafeId:     cafeId,
			Access:     access,
			Refresh:    refresh,
			Expiry:     time.Unix(int64(expiryInt), 0),
			HttpAddr:   httpAddr,
			SwarmAddrs: slist,
		})
	}
	return ret
}
