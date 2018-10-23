package db

import (
	"database/sql"
	"github.com/gogo/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type ThreadMessageDB struct {
	modelStore
}

func NewThreadMessageStore(db *sql.DB, lock *sync.Mutex) repo.ThreadMessageStore {
	return &ThreadMessageDB{modelStore{db, lock}}
}

func (c *ThreadMessageDB) Add(msg *repo.ThreadMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into thread_messages(id, peerId, envelope, date) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	// marshal envelope
	env, err := proto.Marshal(msg.Envelope)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		msg.Id,
		msg.PeerId,
		env,
		int(msg.Date.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *ThreadMessageDB) List(offset string, limit int) []repo.ThreadMessage {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from thread_messages where date>(select date from thread_messages where id='" + offset + "') order by date asc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from thread_messages order by date asc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *ThreadMessageDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from thread_messages where id=?", id)
	return err
}

func (c *ThreadMessageDB) handleQuery(stm string) []repo.ThreadMessage {
	var ret []repo.ThreadMessage
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId string
		var dateInt int
		var envelopeb []byte
		if err := rows.Scan(&id, &peerId, &envelopeb, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		// unmarshal envelope
		env := new(pb.Envelope)
		if err := proto.Unmarshal(envelopeb, env); err != nil {
			log.Errorf("error unmarshaling envelope: %s", err)
			continue
		}
		msg := repo.ThreadMessage{
			Id:       id,
			PeerId:   peerId,
			Envelope: env,
			Date:     time.Unix(int64(dateInt), 0),
		}
		ret = append(ret, msg)
	}
	return ret
}
