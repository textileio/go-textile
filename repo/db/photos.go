package db

import (
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type PhotoDB struct {
	modelStore
}

func NewPhotoStore(db *sql.DB, lock *sync.Mutex) repo.PhotoStore {
	return &PhotoDB{modelStore{db, lock}}
}

func (c *PhotoDB) Put(cid string, thumb string, timestamp time.Time) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into photos(cid, thumb, timestamp) values(?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(
		cid,
		thumb,
		int(timestamp.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *PhotoDB) GetPhotos(offsetId string, limit int) []repo.PhotoSet {
	c.lock.Lock()
	defer c.lock.Unlock()
	var ret []repo.PhotoSet

	var stm string
	if offsetId != "" {
		stm = "select cid, thumb, timestamp from photos where timestamp<(select timestamp from photos where cid='" + offsetId + "') order by timestamp desc limit " + strconv.Itoa(limit) + " ;"
	} else {
		stm = "select cid, thumb, timestamp from photos order by timestamp desc limit " + strconv.Itoa(limit) + ";"
	}
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Error(err)
		return ret
	}
	for rows.Next() {
		var cid string
		var thumb string
		var timestampInt int
		if err := rows.Scan(&cid, &thumb, &timestampInt); err != nil {
			continue
		}
		timestamp := time.Unix(int64(timestampInt), 0)
		photo := repo.PhotoSet{
			Cid:       cid,
			Thumb:     thumb,
			Timestamp: timestamp,
		}
		ret = append(ret, photo)
	}
	return ret
}

func (c *PhotoDB) DeletePhoto(cid string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.db.Exec("delete from photos where cid=?", cid)
	return nil
}
