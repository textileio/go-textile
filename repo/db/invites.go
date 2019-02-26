package db

import (
	"database/sql"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
)

type InviteDB struct {
	modelStore
}

func NewInviteStore(db *sql.DB, lock *sync.Mutex) repo.InviteStore {
	return &InviteDB{modelStore{db, lock}}
}

func (c *InviteDB) Add(invite *pb.Invite) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into invites(id, block, name, inviter, date) values(?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	inviter, err := proto.Marshal(invite.Inviter)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		invite.Id,
		invite.Block,
		invite.Name,
		inviter,
		util.ProtoNanos(invite.Date),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *InviteDB) Get(id string) *pb.Invite {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from invites where id='" + id + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
}

func (c *InviteDB) List() *pb.InviteList {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.handleQuery("select * from invites order by date desc;")
}

func (c *InviteDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from invites where id=?", id)
	return err
}

func (c *InviteDB) handleQuery(stm string) *pb.InviteList {
	list := &pb.InviteList{Items: make([]*pb.Invite, 0)}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, name string
		var block, inviterb []byte
		var dateInt int64
		if err := rows.Scan(&id, &block, &name, &inviterb, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		inviter := new(pb.Contact)
		if err := proto.Unmarshal(inviterb, inviter); err != nil {
			log.Errorf("error unmarshaling inviter: %s", err)
			continue
		}

		list.Items = append(list.Items, &pb.Invite{
			Id:      id,
			Block:   block,
			Name:    name,
			Inviter: inviter,
			Date:    util.ProtoTs(dateInt),
		})
	}
	return list
}
