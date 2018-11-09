package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"strconv"
	"sync"
	"time"
)

type FileDB struct {
	modelStore
}

func NewFileStore(db *sql.DB, lock *sync.Mutex) repo.FileStore {
	return &FileDB{modelStore{db, lock}}
}

func (c *FileDB) Add(file *repo.File) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into files(id, hash, key, added) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		file.Id,
		file.Hash,
		file.Key,
		int(file.Added.Unix()),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *FileDB) Get(id string) *repo.File {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from files where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *FileDB) List(offset string, limit int) []repo.File {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from files where added<(select added from files where id='" + offset + "') order by added desc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from files order by added desc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *FileDB) ListByHash(hash string) []repo.File {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select * from files where hash='" + hash + "' order by added desc;"
	return c.handleQuery(stm)
}

func (c *FileDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from files;")
	var count int
	row.Scan(&count)
	return count
}

func (c *FileDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from files where id=?", id)
	return err
}

func (c *FileDB) handleQuery(stm string) []repo.File {
	var res []repo.File
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, hash, key string
		var addedInt int
		if err := rows.Scan(&id, &hash, &key, &addedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		res = append(res, repo.File{
			Id:    id,
			Hash:  hash,
			Key:   key,
			Added: time.Unix(int64(addedInt), 0),
		})
	}
	return res
}
