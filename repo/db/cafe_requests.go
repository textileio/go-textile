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
	stmt, err := tx.Prepare(`
        INSERT INTO cafe_requests(
    	    id, peerId, targetId, cafeId, cafe, groupId, syncGroupId, type, date, size, status, attempts
        ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
    `)
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

	res := c.handleQuery("SELECT * FROM cafe_requests WHERE id='" + id + "';")
	if len(res.Items) == 0 {
		return nil
	}

	return res.Items[0]
}

func (c *CafeRequestDB) GetGroup(group string) *pb.CafeRequestList {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.handleQuery("SELECT * FROM cafe_requests WHERE groupId='" + group + "';")
}

func (c *CafeRequestDB) List(offset string, limit int) *pb.CafeRequestList {
	c.lock.Lock()
	defer c.lock.Unlock()

	stm := "SELECT * FROM cafe_requests WHERE status=0"
	if offset != "" {
		stm += " AND date>(SELECT date FROM cafe_requests WHERE id='" + offset + "')"
	}
	stm += " ORDER BY date ASC LIMIT " + strconv.Itoa(limit) + ";"

	return c.handleQuery(stm)
}

func (c *CafeRequestDB) ListGroups(offset string, limit int) []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	stm := "SELECT DISTINCT groupId FROM cafe_requests WHERE status=0"
	if offset != "" {
		stm += " AND date>(SELECT date FROM cafe_requests WHERE groupId='" + offset + "')"
	}
	stm += " ORDER BY date ASC LIMIT " + strconv.Itoa(limit) + ";"

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
	total, err := c.db.Query(`
        SELECT DISTINCT syncGroupId FROM cafe_requests WHERE status!=? ORDER BY date ASC;
    `, pb.CafeRequest_COMPLETE)
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

func (c *CafeRequestDB) SyncGroupStatus(groupId string) *pb.CafeRequestSyncGroupStatus {
	c.lock.Lock()
	defer c.lock.Unlock()
	group := &pb.CafeRequestSyncGroupStatus{}

	rows, err := c.db.Query(`
        SELECT cafeId, size, status FROM cafe_requests WHERE syncGroupId=(
            SELECT syncGroupId FROM cafe_requests WHERE groupId=?
        ) ORDER BY date ASC;
	`, groupId)
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

	_, err := c.db.Exec("UPDATE cafe_requests SET status=? WHERE id=?", int32(status), id)
	return err
}

func (c *CafeRequestDB) UpdateGroupStatus(groupId string, status pb.CafeRequest_Status) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("UPDATE cafe_requests SET status=? WHERE groupId=?", int32(status), groupId)
	return err
}

func (c *CafeRequestDB) AddAttempt(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("UPDATE cafe_requests SET attempts=attempts+1 WHERE id=?", id)
	return err
}

func (c *CafeRequestDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("DELETE FROM cafe_requests WHERE id=?", id)
	return err
}

func (c *CafeRequestDB) DeleteByGroup(groupId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("DELETE FROM cafe_requests WHERE groupId=?", groupId)
	return err
}

func (c *CafeRequestDB) DeleteBySyncGroup(syncGroupId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("DELETE FROM cafe_requests WHERE syncGroupId=?", syncGroupId)
	return err
}

func (c *CafeRequestDB) DeleteCompleteSyncGroups() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec(`
        DELETE FROM cafe_requests WHERE syncGroupId=(
		    SELECT a.syncGroupId
		    FROM   (SELECT syncGroupId, COUNT(*) as total
		            FROM   cafe_requests
		            GROUP BY syncGroupId) a
		    JOIN   (SELECT syncGroupId, COUNT(*) as total_complete
		            FROM   cafe_requests
    		        WHERE  status=?
	    	        GROUP BY syncGroupId) b
		    ON     a.syncGroupId = b.syncGroupId AND a.total = b.total_complete
        )
	`, pb.CafeRequest_COMPLETE)
	return err
}

func (c *CafeRequestDB) DeleteByCafe(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec("DELETE FROM cafe_requests WHERE cafeId=?", cafeId)
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
