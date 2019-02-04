package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type CafeTokenDB struct {
	modelStore
}

func NewCafeTokenStore(db *sql.DB, lock *sync.Mutex) repo.CafeTokenStore {
	return &CafeTokenDB{modelStore{db, lock}}
}

func (c *CafeTokenDB) Add(token *repo.CafeToken) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_tokens(id, token, date) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		token.Id,
		token.Token,
		token.Date.UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeTokenDB) Get(id string) *repo.CafeToken {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from cafe_tokens where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeTokenDB) List() []repo.CafeToken {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_tokens order by id desc;"
	return c.handleQuery(stm)
}

func (c *CafeTokenDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_tokens where id=?", id)
	return err
}

func (c *CafeTokenDB) handleQuery(stm string) []repo.CafeToken {
	var ret []repo.CafeToken
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id string
		var token []byte
		var dateInt int64
		if err := rows.Scan(&id, &token, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		ret = append(ret, repo.CafeToken{
			Id:    id,
			Token: token,
			Date:  time.Unix(0, dateInt),
		})
	}
	return ret
}
