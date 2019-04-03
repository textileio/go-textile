package db

import (
	"bytes"
	"database/sql"
	"strconv"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeRequestDB struct {
	modelStore
}

func NewCafeRequestStore(db *sql.DB, lock *sync.Mutex) repo.CafeRequestStore {
	return &CafeRequestDB{modelStore{db, lock}}
}

func (c *CafeRequestDB) Add(req *pb.CafeRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_requests(id, peerId, targetId, cafeId, cafe, type, date, size, groupId) values(?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	cafe, err := pbMarshaler.MarshalToString(req.Cafe)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		req.Id,
		req.Peer,
		req.Target,
		req.Cafe.Peer,
		[]byte(cafe),
		int32(req.Type),
		util.ProtoNanos(req.Date),
		req.Size,
		req.Group,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeRequestDB) List(offset string, limit int) []pb.CafeRequest {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from cafe_requests where date>(select date from cafe_requests where id='" + offset + "') order by date asc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from cafe_requests order by date asc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *CafeRequestDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_requests where id=?", id)
	return err
}

func (c *CafeRequestDB) DeleteByCafe(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_requests where cafeId=?", cafeId)
	return err
}

func (c *CafeRequestDB) handleQuery(stm string) []pb.CafeRequest {
	var list []pb.CafeRequest
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId, targetId, cafeId, groupId string
		var typeInt int
		var dateInt, size int64
		var cafe []byte
		if err := rows.Scan(&id, &peerId, &targetId, &cafeId, &cafe, &typeInt, &dateInt, &size, &groupId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		mod := new(pb.Cafe)
		if err := jsonpb.Unmarshal(bytes.NewReader(cafe), mod); err != nil {
			log.Errorf("error unmarshaling cafe: %s", err)
			continue
		}

		list = append(list, pb.CafeRequest{
			Id:     id,
			Peer:   peerId,
			Target: targetId,
			Cafe:   mod,
			Type:   pb.CafeRequest_Type(typeInt),
			Date:   util.ProtoTs(dateInt),
			Size:   size,
			Group:  groupId,
		})
	}
	return list
}
