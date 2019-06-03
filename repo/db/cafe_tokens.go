package db

import (
	"database/sql"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeTokenDB struct {
	modelStore
}

func NewCafeTokenStore(db *sql.DB, lock *sync.Mutex) repo.CafeTokenStore {
	return &CafeTokenDB{modelStore{db, lock}}
}

func (c *CafeTokenDB) Add(token *pb.CafeToken) error {
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
		token.Value,
		util.ProtoNanos(token.Date),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *CafeTokenDB) Get(id string) *pb.CafeToken {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from cafe_tokens where id='" + id + "';")
	if len(res) == 0 {
		return nil
	}
	return &res[0]
}

func (c *CafeTokenDB) List() []pb.CafeToken {
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

func (c *CafeTokenDB) handleQuery(stm string) []pb.CafeToken {
	var list []pb.CafeToken
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
		list = append(list, pb.CafeToken{
			Id:    id,
			Value: token,
			Date:  util.ProtoTs(dateInt),
		})
	}
	return list
}
