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

type groupCount struct {
	total    int
	complete int
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
	stm := `insert into cafe_requests(id, peerId, targetId, cafeId, cafe, type, date, size, groupId, status) values(?,?,?,?,?,?,?,?,?,?)`
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
		int32(req.Status),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeRequestDB) Get(id string) *pb.CafeRequest {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from cafe_requests where id='" + id + "';")
	if len(res.Items) == 0 {
		return nil
	}
	return res.Items[0]
}

func (c *CafeRequestDB) List(offset string, limit int) *pb.CafeRequestList {
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

func (c *CafeRequestDB) ListCompletedGroups() []string {
	c.lock.Lock()
	defer c.lock.Unlock()
	groups := make(map[string]*groupCount)

	total, err := c.db.Query("select Count(*), groupId from cafe_requests group by groupId;")
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for total.Next() {
		var count int
		var groupId string
		if err := total.Scan(&count, &groupId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		if groups[groupId] == nil {
			groups[groupId] = &groupCount{}
		}
		groups[groupId].total = count
	}

	complete, err := c.db.Query("select Count(*), groupId from cafe_requests where status=2 group by groupId;")
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for complete.Next() {
		var count int
		var groupId string
		if err := complete.Scan(&count, &groupId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		groups[groupId].complete = count
	}

	var list []string
	for g, counts := range groups {
		if counts.complete == counts.total {
			list = append(list, g)
		}
	}

	return list
}

func (c *CafeRequestDB) CountByGroup(groupId string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from cafe_requests where groupId='" + groupId + "';")
	var count int
	row.Scan(&count)
	return count
}

func (c *CafeRequestDB) GroupStatus(groupId string) *pb.CafeRequestGroupStatus {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.handleStatQuery("select cafeId, size, status from cafe_requests where groupId='" + groupId + "' order by date asc;")
}

func (c *CafeRequestDB) UpdateStatus(id string, status pb.CafeRequest_Status) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update cafe_requests set status=? where id=?", int32(status), id)
	return err
}

func (c *CafeRequestDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_requests where id=?", id)
	return err
}

func (c *CafeRequestDB) DeleteByGroup(groupId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_requests where groupId=?", groupId)
	return err
}

func (c *CafeRequestDB) DeleteByCafe(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_requests where cafeId=?", cafeId)
	return err
}

func (c *CafeRequestDB) handleQuery(stm string) *pb.CafeRequestList {
	list := &pb.CafeRequestList{Items: make([]*pb.CafeRequest, 0)}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, peerId, targetId, cafeId, groupId string
		var typeInt, statusInt int
		var dateInt, size int64
		var cafe []byte
		if err := rows.Scan(&id, &peerId, &targetId, &cafeId, &cafe, &typeInt, &dateInt, &size, &groupId, &statusInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		mod := new(pb.Cafe)
		if err := jsonpb.Unmarshal(bytes.NewReader(cafe), mod); err != nil {
			log.Errorf("error unmarshaling cafe: %s", err)
			continue
		}

		list.Items = append(list.Items, &pb.CafeRequest{
			Id:     id,
			Peer:   peerId,
			Target: targetId,
			Cafe:   mod,
			Type:   pb.CafeRequest_Type(typeInt),
			Date:   util.ProtoTs(dateInt),
			Size:   size,
			Group:  groupId,
			Status: pb.CafeRequest_Status(statusInt),
		})
	}
	return list
}

func (c *CafeRequestDB) handleStatQuery(stm string) *pb.CafeRequestGroupStatus {
	group := &pb.CafeRequestGroupStatus{}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return group
	}
	for rows.Next() {
		var cafeId string
		var size int64
		var status int
		if err := rows.Scan(&cafeId, &size, &status); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		group.NumTotal += 1
		group.SizeTotal += size
		switch status {
		case 1:
			group.NumPending += 1
			group.SizePending += size
		case 2:
			group.NumComplete += 1
			group.SizeComplete += size
		}
	}
	return group
}
