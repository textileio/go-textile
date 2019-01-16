package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type ThreadInviteDB struct {
	modelStore
}

func NewThreadInviteStore(db *sql.DB, lock *sync.Mutex) repo.ThreadInviteStore {
	return &ThreadInviteDB{modelStore{db, lock}}
}

func (c *ThreadInviteDB) Add(invite *repo.ThreadInvite) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into thread_invites(id, block, name, inviter, date) values(?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		invite.Id,
		invite.Block,
		invite.Name,
		invite.Inviter,
		invite.Date.UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ThreadInviteDB) Get(id string) *repo.ThreadInvite {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from thread_invites where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ThreadInviteDB) List() []repo.ThreadInvite {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.handleQuery("select * from thread_invites order by date desc;")
}

func (c *ThreadInviteDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_invites where id=?", id)
	return err
}

func (c *ThreadInviteDB) handleQuery(stm string) []repo.ThreadInvite {
	var ret []repo.ThreadInvite
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, name, inviter string
		var block []byte
		var dateInt int64
		if err := rows.Scan(&id, &block, &name, &inviter, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		ret = append(ret, repo.ThreadInvite{
			Id:      id,
			Block:   block,
			Name:    name,
			Inviter: inviter,
			Date:    time.Unix(0, dateInt),
		})
	}
	return ret
}
