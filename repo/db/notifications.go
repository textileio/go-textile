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
	stm := `insert into notifications(id, date, actorId, targetId, type, read) values(?,?,?,?,?,?)`
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
		notification.TargetId,
		int(notification.Type),
		false,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
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

func (c *NotificationDB) DeleteByTargetId(targetId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where targetId=?", targetId)
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
		var id, actorId, targetId string
		var dateInt, typeInt, readInt int
		if err := rows.Scan(&id, &dateInt, &actorId, &targetId, &typeInt, &readInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		read := false
		if readInt == 1 {
			read = true
		}
		notif := repo.Notification{
			Id:       id,
			Date:     time.Unix(int64(dateInt), 0),
			ActorId:  actorId,
			TargetId: targetId,
			Type:     repo.NotificationType(typeInt),
			Read:     read,
		}
		ret = append(ret, notif)
	}
	return ret
}
