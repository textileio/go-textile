package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type CafeInboxDB struct {
	modelStore
}

func NewCafeInboxStore(db *sql.DB, lock *sync.Mutex) repo.CafeInboxStore {
	return &CafeInboxDB{modelStore{db, lock}}
}

func (c *CafeInboxDB) Add(inbox *repo.CafeInbox) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_inboxes(peerId, cafeId) values(?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		inbox.PeerId,
		inbox.CafeId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeInboxDB) ListByPeer(peerId string) []repo.CafeInbox {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_inboxes where peerId='" + peerId + "';"
	return c.handleQuery(stm)
}

func (c *CafeInboxDB) DeleteByPeer(peerId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_inboxes where peerId=?", peerId)
	return err
}

func (c *CafeInboxDB) DeleteByCafe(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_inboxes where cafeId=?", cafeId)
	return err
}

func (c *CafeInboxDB) handleQuery(stm string) []repo.CafeInbox {
	var ret []repo.CafeInbox
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var peerId, cafeId string
		if err := rows.Scan(&peerId, &cafeId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		block := repo.CafeInbox{
			PeerId: peerId,
			CafeId: cafeId,
		}
		ret = append(ret, block)
	}
	return ret
}
