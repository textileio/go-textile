package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type CafeDevTokenDB struct {
	modelStore
}

func NewCafeDevTokenStore(db *sql.DB, lock *sync.Mutex) repo.CafeDevTokenStore {
	return &CafeDevTokenDB{modelStore{db, lock}}
}

func (c *CafeDevTokenDB) Add(token *repo.CafeDevToken) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_dev_tokens(id, token, created) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		token.Id,
		token.Token,
		token.Created.UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeDevTokenDB) Get(id string) *repo.CafeDevToken {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from cafe_dev_tokens where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *CafeDevTokenDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from cafe_dev_tokens;")
	var count int
	row.Scan(&count)
	return count
}

func (c *CafeDevTokenDB) List() []repo.CafeDevToken {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_dev_tokens order by id desc;"
	return c.handleQuery(stm)
}

// func (c *CafeDevTokenDB) ListByClient(clientId string) []repo.CafeDevToken {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	stm := "select * from cafe_dev_tokens where client='" + clientId + "' order by id desc;"
// 	return c.handleQuery(stm)
// }

func (c *CafeDevTokenDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_dev_tokens where id=?", id)
	return err
}

func (c *CafeDevTokenDB) handleQuery(stm string) []repo.CafeDevToken {
	var ret []repo.CafeDevToken
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id string
		var token []byte
		var createdInt int64
		if err := rows.Scan(&id, &token, &createdInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		ret = append(ret, repo.CafeDevToken{
			Id:      id,
			Token:   token,
			Created: time.Unix(0, createdInt),
		})
	}
	return ret
}
