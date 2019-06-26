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
    	    id, peerId, targetId, cafeId, cafe, groupId, syncGroupId, type, date, size, status, attempts, groupSize, groupTransferred
        ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
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
		req.GroupSize,
		req.GroupTransferred,
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

func (c *CafeRequestDB) GetSyncGroup(group string) string {
	c.lock.Lock()
	defer c.lock.Unlock()

	total, err := c.db.Query(`
		SELECT DISTINCT syncGroupId FROM cafe_requests WHERE groupId=?
	`, group)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return ""
	}

	for total.Next() {
		var syncGroupId string
		err = total.Scan(&syncGroupId)
		if err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		return syncGroupId
	}

	return ""
}

func (c *CafeRequestDB) Count(status pb.CafeRequest_Status) int {
	c.lock.Lock()
	defer c.lock.Unlock()

	stm := "SELECT COUNT(*) FROM cafe_requests"
	if status != -1 {
		stm += " WHERE status=" + strconv.Itoa(int(status))
	}
	stm += ";"

	row := c.db.QueryRow(stm)
	var count int
	_ = row.Scan(&count)

	return count
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
		err = total.Scan(&groupId)
		if err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		groups = append(groups, groupId)
	}

	return groups
}

func (c *CafeRequestDB) SyncGroupComplete(syncGroupId string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	row := c.db.QueryRow("SELECT COUNT(*) FROM cafe_requests WHERE syncGroupId=?;", syncGroupId)
	var count int
	_ = row.Scan(&count)
	if count == 0 {
		return true
	}

	var syncGroups []string
	total, err := c.db.Query(`
        SELECT a.syncGroupId
		FROM   (SELECT syncGroupId, COUNT(*) as total
		        FROM   cafe_requests
		        WHERE  syncGroupId=?
		        GROUP BY syncGroupId) a
		JOIN   (SELECT syncGroupId, COUNT(*) as total_complete
		        FROM   cafe_requests
    		    WHERE  syncGroupId=? AND status=?
	    	    GROUP BY syncGroupId) b
		ON     a.syncGroupId = b.syncGroupId AND a.total = b.total_complete
    `, syncGroupId, syncGroupId, pb.CafeRequest_COMPLETE)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return false
	}
	for total.Next() {
		var syncGroupId string
		err = total.Scan(&syncGroupId)
		if err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		syncGroups = append(syncGroups, syncGroupId)
	}

	return len(syncGroups) > 0
}

func (c *CafeRequestDB) SyncGroupStatus(groupId string) *pb.CafeSyncGroupStatus {
	c.lock.Lock()
	defer c.lock.Unlock()
	status := &pb.CafeSyncGroupStatus{}

	rows, err := c.db.Query(`
        SELECT cafeId, groupId, syncGroupId, size, status, groupSize, groupTransferred FROM cafe_requests WHERE syncGroupId=(
            SELECT syncGroupId FROM cafe_requests WHERE groupId=?
        ) ORDER BY date ASC;
	`, groupId)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return status
	}

	groups := make(map[string]struct{})
	for rows.Next() {
		var cafeId, groupId, syncGroupId string
		var size, groupSize, groupTransferred int64
		var stat int
		err = rows.Scan(&cafeId, &groupId, &syncGroupId, &size, &stat, &groupSize, &groupTransferred)
		if err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		status.Id = syncGroupId
		status.NumTotal += 1
		status.SizeTotal += size
		switch stat {
		case int(pb.CafeRequest_PENDING):
			status.NumPending += 1
			status.SizePending += size
		case int(pb.CafeRequest_COMPLETE):
			status.NumComplete += 1
			status.SizeComplete += size
		}

		if _, ok := groups[groupId]; !ok {
			groups[groupId] = struct{}{}
			status.GroupsSizeTotal += groupSize
			status.GroupsSizeComplete += groupTransferred
		}
	}

	return status
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

func (c *CafeRequestDB) UpdateGroupProgress(groupId string, transferred int64, total int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.db.Exec(
		"UPDATE cafe_requests SET groupSize=?, groupTransferred=? WHERE groupId=?",
		total, transferred, groupId)
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
        DELETE FROM cafe_requests WHERE syncGroupId IN (
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
	if err != nil {
		return err
	}

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
		var dateInt, size, groupSize, groupTransferred int64
		var cafe []byte

		err = rows.Scan(
			&id,
			&peerId,
			&targetId,
			&cafeId,
			&cafe,
			&groupId,
			&syncGroupId,
			&typeInt,
			&dateInt,
			&size,
			&statusInt,
			&attempts,
			&groupSize,
			&groupTransferred)
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
			Id:               id,
			Peer:             peerId,
			Target:           targetId,
			Cafe:             mod,
			Group:            groupId,
			SyncGroup:        syncGroupId,
			Type:             pb.CafeRequest_Type(typeInt),
			Date:             util.ProtoTs(dateInt),
			Size:             size,
			Status:           pb.CafeRequest_Status(statusInt),
			Attempts:         int32(attempts),
			GroupSize:        groupSize,
			GroupTransferred: groupTransferred,
		})
	}

	return list
}
