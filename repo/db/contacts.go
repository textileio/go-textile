package db

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type ContactDB struct {
	modelStore
}

func NewContactStore(db *sql.DB, lock *sync.Mutex) repo.ContactStore {
	return &ContactDB{modelStore{db, lock}}
}

func (c *ContactDB) Add(contact *repo.Contact) error {
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

func (c *ContactDB) AddOrUpdate(contact *repo.Contact) error {
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

	_, err = stmt.Exec(
		contact.Id,
		contact.Address,
		contact.Username,
		contact.Avatar,
		inboxes,
		contact.Id,
		contact.Created.UnixNano(),
		time.Now().UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ContactDB) Get(id string) *repo.Contact {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from contacts where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *ContactDB) List() []repo.Contact {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.handleQuery("select * from contacts order by username asc;")
}

func (c *ContactDB) Find(id string, address string, username string) []repo.Contact {
	c.lock.Lock()
	defer c.lock.Unlock()
	if id != "" {
		return c.handleQuery("select * from contacts where id='" + id + "';")
	}
	if address == "" && username == "" {
		return nil
	}
	stm := "select * from contacts where"
	if address != "" {
		stm += " address='" + address + "'"
		if username != "" {
			stm += " and"
		}
	}
	if username != "" {
		stm += " username='" + username + "'"
	}
	stm += " order by updated desc;"
	return c.handleQuery(stm)
}

func (c *ContactDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from contacts;")
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

func (c *ContactDB) UpdateInboxes(id string, inboxes []repo.Cafe) error {
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

func (c *ContactDB) handleQuery(stm string) []repo.Contact {
	var ret []repo.Contact
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, address, username, avatar string
		var inboxes []byte
		var createdInt, updatedInt int64
		if err := rows.Scan(&id, &address, &username, &avatar, &inboxes, &createdInt, &updatedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		ilist := make([]repo.Cafe, 0)
		if err := json.Unmarshal(inboxes, &ilist); err != nil {
			log.Errorf("error unmarshaling cafes: %s", err)
			continue
		}

		ret = append(ret, repo.Contact{
			Id:       id,
			Address:  address,
			Username: username,
			Avatar:   avatar,
			Inboxes:  ilist,
			Created:  time.Unix(0, createdInt),
			Updated:  time.Unix(0, updatedInt),
		})
	}
	return ret
}
