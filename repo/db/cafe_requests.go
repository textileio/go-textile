package db

import (
	"bytes"
	"database/sql"
	"strconv"
	"sync"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type CafeRequestDB struct {
	modelStore
}

type syncGroupCount struct {
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
	stm := `insert into cafe_requests(
    	id, peerId, targetId, cafeId, cafe, groupId, syncGroupId, type, date, size, status, attempts
    ) values(?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
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
		req.Group,
		req.SyncGroup,
		int32(req.Type),
		util.ProtoNanos(req.Date),
		req.Size,
		int32(req.Status),
		req.Attempts,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
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

	stm := "select * from cafe_requests where status=0"
	if offset != "" {
		stm += " and date>(select date from cafe_requests where id='" + offset + "')"
	}
	stm += " order by date asc limit " + strconv.Itoa(limit) + ";"

	return c.handleQuery(stm)
}

func (c *CafeRequestDB) ListGroups(offset string, limit int) []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	stm := "select distinct groupId from cafe_requests where status=0"
	if offset != "" {
		stm += " and date>(select date from cafe_requests where groupId='" + offset + "')"
	}
	stm += " order by date asc limit " + strconv.Itoa(limit) + ";"

	var groups []string
	total, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for total.Next() {
		var groupId string
		if err := total.Scan(&groupId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		groups = append(groups, groupId)
	}

	return groups
}

func (c *CafeRequestDB) ListIncompleteSyncGroups() []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	var syncGroups []string
	total, err := c.db.Query("select distinct syncGroupId from cafe_requests where status!=2 order by date asc;")
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for total.Next() {
		var syncGroupId string
		if err := total.Scan(&syncGroupId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		syncGroups = append(syncGroups, syncGroupId)
	}

	return syncGroups
}

// a new
// a pending
// a complete

// b complete
// b complete
// b complete

// c new
// c new
// c new

// -> b

func (c *CafeRequestDB) ListCompleteSyncGroups() []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	syncGroups := make(map[string]*syncGroupCount)
	total, err := c.db.Query("select Count(*), syncGroupId from cafe_requests group by syncGroupId;")
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for total.Next() {
		var count int
		var syncGroupId string
		if err := total.Scan(&count, &syncGroupId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		if syncGroups[syncGroupId] == nil {
			syncGroups[syncGroupId] = &syncGroupCount{}
		}
		syncGroups[syncGroupId].total = count
	}

	complete, err := c.db.Query("select Count(*), syncGroupId from cafe_requests where status=2 group by syncGroupId;")
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for complete.Next() {
		var count int
		var syncGroupId string
		if err := complete.Scan(&count, &syncGroupId); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		syncGroups[syncGroupId].complete = count
	}

	var list []string
	for g, counts := range syncGroups {
		if counts.complete == counts.total {
			list = append(list, g)
		}
	}

	return list
}

func (c *CafeRequestDB) SyncGroupStatus(syncGroupId string) *pb.CafeRequestSyncGroupStatus {
	c.lock.Lock()
	defer c.lock.Unlock()
	group := &pb.CafeRequestSyncGroupStatus{}

	stm := "select cafeId, size, status from cafe_requests where syncGroupId='" + syncGroupId + "' order by date asc;"
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

func (c *CafeRequestDB) UpdateStatus(id string, status pb.CafeRequest_Status) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("update cafe_requests set status=? where id=?", int32(status), id)
	return err
}

func (c *CafeRequestDB) UpdateGroupStatus(groupId string, status pb.CafeRequest_Status) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("update cafe_requests set status=? where groupId=?", int32(status), groupId)
	return err
}

func (c *CafeRequestDB) AddAttempt(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("update cafe_requests set attempts=attempts+1 where id=?", id)
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

func (c *CafeRequestDB) DeleteBySyncGroup(syncGroupId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("delete from cafe_requests where syncGroupId=?", syncGroupId)
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
		var id, peerId, targetId, cafeId, groupId, syncGroupId string
		var typeInt, statusInt, attempts int
		var dateInt, size int64
		var cafe []byte

		err := rows.Scan(&id, &peerId, &targetId, &cafeId, &cafe, &groupId, &syncGroupId, &typeInt, &dateInt, &size, &statusInt, &attempts)
		if err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		mod := new(pb.Cafe)
		err = pbUnmarshaler.Unmarshal(bytes.NewReader(cafe), mod)
		if err != nil {
			log.Errorf("error unmarshaling cafe: %s", err)
			continue
		}

		list.Items = append(list.Items, &pb.CafeRequest{
			Id:        id,
			Peer:      peerId,
			Target:    targetId,
			Cafe:      mod,
			Group:     groupId,
			SyncGroup: syncGroupId,
			Type:      pb.CafeRequest_Type(typeInt),
			Date:      util.ProtoTs(dateInt),
			Size:      size,
			Status:    pb.CafeRequest_Status(statusInt),
			Attempts:  int32(attempts),
		})
	}

	return list
}
