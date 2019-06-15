package db

import (
	"database/sql"
	"strconv"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type NotificationDB struct {
	modelStore
}

func NewNotificationStore(db *sql.DB, lock *sync.Mutex) repo.NotificationStore {
	return &NotificationDB{modelStore{db, lock}}
}

func (c *NotificationDB) Add(notification *pb.Notification) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into notifications(id, date, actorId, subject, subjectId, blockId, target, type, body, read) values(?,?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		notification.Id,
		util.ProtoNanos(notification.Date),
		notification.Actor,
		notification.SubjectDesc,
		notification.Subject,
		notification.Block,
		notification.Target,
		int32(notification.Type),
		notification.Body,
		false,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *NotificationDB) Get(id string) *pb.Notification {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from notifications where id='" + id + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
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

func (c *NotificationDB) List(offset string, limit int) *pb.NotificationList {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from notifications where date<(select date from notifications where id='" + offset + "') order by date desc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from notifications order by date desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *NotificationDB) CountUnread() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from notifications where read=0;")
	var count int
	_ = row.Scan(&count)
	return count
}

func (c *NotificationDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where id=?", id)
	return err
}

func (c *NotificationDB) DeleteByActor(actorId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where actorId=?", actorId)
	return err
}

func (c *NotificationDB) DeleteBySubject(subjectId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where subjectId=?", subjectId)
	return err
}

func (c *NotificationDB) DeleteByBlock(blockId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from notifications where blockId=?", blockId)
	return err
}

func (c *NotificationDB) handleQuery(stm string) *pb.NotificationList {
	list := &pb.NotificationList{Items: make([]*pb.Notification, 0)}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, actorId, subject, subjectId, blockId, target, body string
		var dateInt int64
		var typeInt, readInt int
		if err := rows.Scan(&id, &dateInt, &actorId, &subject, &subjectId, &blockId, &target, &typeInt, &body, &readInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		read := false
		if readInt == 1 {
			read = true
		}
		list.Items = append(list.Items, &pb.Notification{
			Id:          id,
			Date:        util.ProtoTs(dateInt),
			Actor:       actorId,
			SubjectDesc: subject,
			Subject:     subjectId,
			Block:       blockId,
			Target:      target,
			Type:        pb.Notification_Type(typeInt),
			Body:        body,
			Read:        read,
		})
	}
	return list
}
