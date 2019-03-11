package db

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type ContactDB struct {
	modelStore
}

func NewContactStore(db *sql.DB, lock *sync.Mutex) repo.ContactStore {
	return &ContactDB{modelStore{db, lock}}
}

func (c *ContactDB) Add(contact *pb.Contact) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into contacts(id, address, username, avatar, inboxes, created, updated) values(?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	inboxes, err := json.Marshal(contact.Inboxes)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		contact.Id,
		contact.Address,
		contact.Name,
		contact.Avatar,
		inboxes,
		time.Now().UnixNano(),
		time.Now().UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ContactDB) AddOrUpdate(contact *pb.Contact) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into contacts(id, address, username, avatar, inboxes, created, updated) values(?,?,?,?,?,coalesce((select created from contacts where id=?),?),?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	inboxes, err := json.Marshal(contact.Inboxes)
	if err != nil {
		return err
	}

	var created int64
	if contact.Created == nil {
		created = time.Now().UnixNano()
	} else {
		created = util.ProtoNanos(contact.Created)
	}

	_, err = stmt.Exec(
		contact.Id,
		contact.Address,
		contact.Name,
		contact.Avatar,
		inboxes,
		contact.Id,
		created,
		time.Now().UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ContactDB) Get(id string) *pb.Contact {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from contacts where id='" + id + "';")
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

func (c *ContactDB) GetBest(id string) *pb.Contact {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select *, (select address from contacts where id='" + id + "') as addr from contacts where address=addr order by updated desc limit 1;"
	row := c.db.QueryRow(stm)
	var _id, address, username, avatar, addr string
	var inboxes []byte
	var createdInt, updatedInt int64
	if err := row.Scan(&_id, &address, &username, &avatar, &inboxes, &createdInt, &updatedInt, &addr); err != nil {
		return nil
	}
	return c.handleRow(id, address, username, avatar, inboxes, createdInt, updatedInt)
}

func (c *ContactDB) List(query string) []*pb.Contact {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := "select * from contacts"
	if query != "" {
		q += " where " + query
	}
	q += " order by updated desc;"
	return c.handleQuery(q)
}

func (c *ContactDB) Find(address string, name string, exclude []string) []*pb.Contact {
	c.lock.Lock()
	defer c.lock.Unlock()
	if address == "" && name == "" {
		return nil
	}
	var q string
	if address != "" {
		q += "address='" + address + "'"
	}
	if name != "" {
		if len(q) > 0 {
			q += " and "
		}
		q += "username like '%" + name + "%'"
	}
	if len(exclude) > 0 {
		q += " and id not in ("
		for i, e := range exclude {
			q += "'" + e + "'"
			if i != len(exclude)-1 {
				q += ","
			}
		}
		q += ")"
	}
	return c.handleQuery("select * from contacts where " + q + " order by updated desc;")
}

func (c *ContactDB) Count(query string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := "select Count(*) from contacts"
	if query != "" {
		q += " where " + query
	}
	q += ";"
	row := c.db.QueryRow(q)
	var count int
	row.Scan(&count)
	return count
}

func (c *ContactDB) UpdateName(id string, name string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update contacts set username=?, updated=? where id=?", name, time.Now().UnixNano(), id)
	return err
}

func (c *ContactDB) UpdateAvatar(id string, avatar string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update contacts set avatar=?, updated=? where id=?", avatar, time.Now().UnixNano(), id)
	return err
}

func (c *ContactDB) UpdateInboxes(id string, inboxes []*pb.Cafe) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	inboxesb, err := json.Marshal(inboxes)
	if err != nil {
		return err
	}
	_, err = c.db.Exec("update contacts set inboxes=?, updated=? where id=?", inboxesb, time.Now().UnixNano(), id)
	return err
}

func (c *ContactDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from contacts where id=?", id)
	return err
}

func (c *ContactDB) DeleteByAddress(address string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from contacts where address=?", address)
	return err
}

func (c *ContactDB) handleQuery(stm string) []*pb.Contact {
	list := make([]*pb.Contact, 0)
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, address, name, avatar string
		var inboxes []byte
		var createdInt, updatedInt int64
		if err := rows.Scan(&id, &address, &name, &avatar, &inboxes, &createdInt, &updatedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		row := c.handleRow(id, address, name, avatar, inboxes, createdInt, updatedInt)
		if row != nil {
			list = append(list, row)
		}
	}
	return list
}

func (c *ContactDB) handleRow(id string, address string, name string, avatar string, inboxes []byte, createdInt int64, updatedInt int64) *pb.Contact {
	cafes := make([]*pb.Cafe, 0)
	if err := json.Unmarshal(inboxes, &cafes); err != nil {
		log.Errorf("error unmarshaling cafes: %s", err)
		return nil
	}

	return &pb.Contact{
		Id:      id,
		Address: address,
		Name:    name,
		Avatar:  avatar,
		Inboxes: cafes,
		Created: util.ProtoTs(createdInt),
		Updated: util.ProtoTs(updatedInt),
	}
}
