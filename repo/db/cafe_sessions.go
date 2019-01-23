package db

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
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

	cafe, err := json.Marshal(session.Cafe)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		session.Id,
		session.Access,
		session.Refresh,
		time.Unix(session.Exp.Seconds, int64(session.Exp.Nanos)).UnixNano(),
		cafe,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeSessionDB) Get(cafeId string) *pb.CafeSession {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from cafe_sessions where cafeId='" + cafeId + "';")
	if len(ret) == 0 {
		return nil
	}
	return ret[0]
}

func (c *CafeSessionDB) List() []*pb.CafeSession {
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

func (c *CafeSessionDB) handleQuery(stm string) []*pb.CafeSession {
	var ret []*pb.CafeSession
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var cafeId, access, refresh string
		var expiryInt int64
		var cafe []byte
		if err := rows.Scan(&cafeId, &access, &refresh, &expiryInt, &cafe); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		var rcafe *pb.Cafe
		if err := json.Unmarshal(cafe, &rcafe); err != nil {
			log.Errorf("error unmarshaling cafe: %s", err)
			continue
		}

		time := time.Unix(0, expiryInt)

		timestamp := timestamp.Timestamp{
			Seconds: time.Unix(),
			Nanos:   int32(time.Nanosecond()),
		}

		ret = append(ret, &pb.CafeSession{
			Id:      cafeId,
			Access:  access,
			Refresh: refresh,
			Exp:     &timestamp,
			Cafe:    rcafe,
		})
	}
	return ret
}
