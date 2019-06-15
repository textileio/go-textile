package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeClientDB struct {
	modelStore
}

func NewCafeClientStore(db *sql.DB, lock *sync.Mutex) repo.CafeClientStore {
	return &CafeClientDB{modelStore{db, lock}}
}

func (c *CafeClientDB) Add(client *pb.CafeClient) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_clients(id, address, created, lastSeen, tokenId) values(?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		client.Id,
		client.Address,
		util.ProtoNanos(client.Created),
		util.ProtoNanos(client.Seen),
		client.Token,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *CafeClientDB) Get(id string) *pb.CafeClient {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from cafe_clients where id='" + id + "';")
	if len(res) == 0 {
		return nil
	}
	return &res[0]
}

func (c *CafeClientDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from cafe_clients;")
	var count int
	_ = row.Scan(&count)
	return count
}

func (c *CafeClientDB) List() []pb.CafeClient {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_clients order by lastSeen desc;"
	return c.handleQuery(stm)
}

func (c *CafeClientDB) ListByAddress(address string) []pb.CafeClient {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_clients where address='" + address + "' order by lastSeen desc;"
	return c.handleQuery(stm)
}

func (c *CafeClientDB) UpdateLastSeen(id string, date time.Time) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update cafe_clients set lastSeen=? where id=?", int64(date.UnixNano()), id)
	return err
}

func (c *CafeClientDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_clients where id=?", id)
	return err
}

func (c *CafeClientDB) handleQuery(stm string) []pb.CafeClient {
	var list []pb.CafeClient
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, address, tokenId string
		var createdInt, lastSeenInt int64
		if err := rows.Scan(&id, &address, &createdInt, &lastSeenInt, &tokenId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		list = append(list, pb.CafeClient{
			Id:      id,
			Address: address,
			Created: util.ProtoTs(createdInt),
			Seen:    util.ProtoTs(lastSeenInt),
			Token:   tokenId,
		})
	}
	return list
}
