package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type NotificationDB struct {
	modelStore
}

func NewNotificationStore(db *sql.DB, lock *sync.Mutex) repo.NotificationStore {
	return &NotificationDB{modelStore{db, lock}}
}

func (c *NotificationDB) Add(notification *repo.Notification) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into notifications(id, date, actorId, actorUsername, subject, subjectId, blockId, dataId, type, body, read) values(?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		notification.Id,
		int(notification.Date.Unix()),
		notification.ActorId,
		notification.ActorUsername,
		notification.Subject,
		notification.SubjectId,
		notification.BlockId,
		notification.DataId,
		int(notification.Type),
		notification.Body,
		false,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *NotificationDB) Get(id string) *repo.Notification {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from notifications where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *NotificationDB) Read(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update notifications set read=1 where id=?", id)
	return err
}

func (c *NotificationDB) ReadAll() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update notifications set read=1")
	return err
}

func (c *NotificationDB) List(offset string, limit int, query string) []repo.Notification {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm, q string
	if offset != "" {
		if query != "" {
			q = query + " and "
		}
		stm = "select * from notifications where " + q + "date<(select date from notifications where id='" + offset + "') order by date desc limit " + strconv.Itoa(limit) + " ;"
	} else {
		if query != "" {
			q = "where " + query + " "
		}
		stm = "select * from notifications " + q + "order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *NotificationDB) CountUnread() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from notifications where read=0;")
	var count int
	row.Scan(&count)
	return count
}

func (c *NotificationDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where id=?", id)
	return err
}

func (c *NotificationDB) DeleteByActorId(actorId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where actorId=?", actorId)
	return err
}

func (c *NotificationDB) DeleteBySubjectId(subjectId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where subjectId=?", subjectId)
	return err
}

func (c *NotificationDB) DeleteByBlockId(blockId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where blockId=?", blockId)
	return err
}

func (c *NotificationDB) handleQuery(stm string) []repo.Notification {
	var ret []repo.Notification
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, actorId, actorUsername, subject, subjectId, blockId, dataId, body string
		var dateInt, typeInt, readInt int
		if err := rows.Scan(&id, &dateInt, &actorId, &actorUsername, &subject, &subjectId, &blockId, &dataId, &typeInt, &body, &readInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		read := false
		if readInt == 1 {
			read = true
		}
		notif := repo.Notification{
			Id:            id,
			Date:          time.Unix(int64(dateInt), 0),
			ActorId:       actorId,
			ActorUsername: actorUsername,
			Subject:       subject,
			SubjectId:     subjectId,
			BlockId:       blockId,
			DataId:        dataId,
			Type:          repo.NotificationType(typeInt),
			Body:          body,
			Read:          read,
		}
		ret = append(ret, notif)
	}
	return ret
}
