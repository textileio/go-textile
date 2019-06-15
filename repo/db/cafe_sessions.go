package db

import (
	"bytes"
	"database/sql"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeSessionDB struct {
	modelStore
}

func NewCafeSessionStore(db *sql.DB, lock *sync.Mutex) repo.CafeSessionStore {
	return &CafeSessionDB{modelStore{db, lock}}
}

func (c *CafeSessionDB) AddOrUpdate(session *pb.CafeSession) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into cafe_sessions(cafeId, access, refresh, expiry, cafe) values(?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	cafe, err := pbMarshaler.MarshalToString(session.Cafe)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		session.Id,
		session.Access,
		session.Refresh,
		util.ProtoNanos(session.Exp),
		[]byte(cafe),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *CafeSessionDB) Get(cafeId string) *pb.CafeSession {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from cafe_sessions where cafeId='" + cafeId + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
}

func (c *CafeSessionDB) List() *pb.CafeSessionList {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from cafe_sessions order by expiry desc;"
	return c.handleQuery(stm)
}

func (c *CafeSessionDB) Delete(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_sessions where cafeId=?", cafeId)
	return err
}

func (c *CafeSessionDB) handleQuery(stm string) *pb.CafeSessionList {
	list := &pb.CafeSessionList{Items: make([]*pb.CafeSession, 0)}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var cafeId, access, refresh string
		var expiryInt int64
		var cafe []byte
		if err := rows.Scan(&cafeId, &access, &refresh, &expiryInt, &cafe); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		rcafe := new(pb.Cafe)
		if err := pbUnmarshaler.Unmarshal(bytes.NewReader(cafe), rcafe); err != nil {
			log.Errorf("error unmarshaling cafe: %s", err)
			continue
		}

		list.Items = append(list.Items, &pb.CafeSession{
			Id:      cafeId,
			Access:  access,
			Refresh: refresh,
			Exp:     util.ProtoTs(expiryInt),
			Cafe:    rcafe,
		})
	}
	return list
}
