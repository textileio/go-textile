package db

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
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
		contact.Username,
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
		contact.Username,
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
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
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

func (c *ContactDB) List(query string) *pb.ContactList {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := "select * from contacts"
	if query != "" {
		q += " where " + query
	}
	q += " order by username asc;"
	return c.handleQuery(q)
}

func (c *ContactDB) Find(id string, address string, username string, exclude []string) *pb.ContactList {
	c.lock.Lock()
	defer c.lock.Unlock()
	if id != "" {
		if util.ListContainsString(exclude, id) {
			return nil
		}
		return c.handleQuery("select * from contacts where id='" + id + "';")
	}
	if address == "" && username == "" {
		return nil
	}
	var q string
	if address != "" {
		q += "address='" + address + "'"
	}
	if username != "" {
		if len(q) > 0 {
			q += " and "
		}
		q += "username like '%" + username + "%'"
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

func (c *ContactDB) UpdateUsername(id string, username string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update contacts set username=?, updated=? where id=?", username, time.Now().UnixNano(), id)
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

func (c *ContactDB) handleQuery(stm string) *pb.ContactList {
	list := &pb.ContactList{Items: make([]*pb.Contact, 0)}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, address, username, avatar string
		var inboxes []byte
		var createdInt, updatedInt int64
		if err := rows.Scan(&id, &address, &username, &avatar, &inboxes, &createdInt, &updatedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		row := c.handleRow(id, address, username, avatar, inboxes, createdInt, updatedInt)
		if row != nil {
			list.Items = append(list.Items, row)
		}
	}
	return list
}

func (c *ContactDB) handleRow(id string, address string, username string, avatar string, inboxes []byte, createdInt int64, updatedInt int64) *pb.Contact {
	cafes := make([]*pb.Cafe, 0)
	if err := json.Unmarshal(inboxes, &cafes); err != nil {
		log.Errorf("error unmarshaling cafes: %s", err)
		return nil
	}

	return &pb.Contact{
		Id:       id,
		Address:  address,
		Username: username,
		Avatar:   avatar,
		Inboxes:  cafes,
		Created:  util.ProtoTs(createdInt),
		Updated:  util.ProtoTs(updatedInt),
	}
}
