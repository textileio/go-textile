package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
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
	stm := `insert into files(mill, checksum, source, opts, hash, key, media, name, size, added, meta, targets) values(?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	var meta []byte
	if file.Meta != nil {
		meta, err = json.Marshal(file.Meta)
		if err != nil {
			return err
		}
	}

	var targets *string
	if len(file.Targets) > 0 {
		tmp := strings.Join(file.Targets, ",")
		targets = &tmp
	}

	_, err = stmt.Exec(
		file.Mill,
		file.Checksum,
		file.Source,
		file.Opts,
		file.Hash,
		file.Key,
		file.Media,
		file.Name,
		file.Size,
		file.Added.UnixNano(),
		meta,
		targets,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *FileDB) Get(hash string) *repo.File {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from files where hash='" + hash + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *FileDB) GetByPrimary(mill string, checksum string) *repo.File {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from files where mill='" + mill + "' and checksum='" + checksum + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *FileDB) GetBySource(mill string, source string, opts string) *repo.File {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from files where mill='" + mill + "' and source='" + source + "' and opts='" + opts + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *FileDB) AddTarget(hash string, target string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	res := c.handleTargetsQuery("select targets from files where hash='" + hash + "';")
	if len(res) == 0 {
		return errors.New("file not found")
	}
	etargets := res[0]

	if targetExists(target, etargets) {
		return nil
	}

	etargets = append(etargets, target)
	targets := strings.Join(etargets, ",")

	_, err := c.db.Exec("update files set targets=? where hash=?", targets, hash)
	return err
}

func (c *FileDB) RemoveTarget(hash string, target string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	res := c.handleTargetsQuery("select targets from files where hash='" + hash + "';")
	if len(res) == 0 {
		return errors.New("file not found")
	}
	etargets := res[0]

	if !targetExists(target, etargets) {
		return nil
	}

	var list []string
	for _, t := range etargets {
		if t != target {
			list = append(list, t)
		}
	}

	var targets *string
	if len(list) > 0 {
		tmp := strings.Join(list, ",")
		targets = &tmp
	}

	_, err := c.db.Exec("update files set targets=? where hash=?", targets, hash)
	return err
}

func (c *FileDB) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	row := c.db.QueryRow("select Count(*) from files;")
	var count int
	row.Scan(&count)
	return count
}

func (c *FileDB) Delete(hash string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from files where hash=?", hash)
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
		var mill, checksum, source, opts, hash, key, media, name string
		var size int
		var addedInt int64
		var metab []byte
		var targets *string

		if err := rows.Scan(&mill, &checksum, &source, &opts, &hash, &key, &media, &name, &size, &addedInt, &metab, &targets); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		var meta map[string]interface{}
		if metab != nil {
			if err := json.Unmarshal(metab, &meta); err != nil {
				log.Errorf("failed to unmarshal file meta: %s", err)
				continue
			}
		}

		tlist := make([]string, 0)
		if targets != nil {
			for _, t := range strings.Split(*targets, ",") {
				if t != "" {
					tlist = append(tlist, t)
				}
			}
		}

		res = append(res, repo.File{
			Mill:     mill,
			Checksum: checksum,
			Source:   source,
			Opts:     opts,
			Hash:     hash,
			Key:      key,
			Media:    media,
			Name:     name,
			Size:     size,
			Added:    time.Unix(0, addedInt),
			Meta:     meta,
			Targets:  tlist,
		})
	}

	return res
}

func (c *FileDB) handleTargetsQuery(stm string) [][]string {
	var res [][]string
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var targets *string

		if err := rows.Scan(&targets); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		tlist := make([]string, 0)
		if targets != nil {
			for _, t := range strings.Split(*targets, ",") {
				if t != "" {
					tlist = append(tlist, t)
				}
			}
		}

		res = append(res, tlist)
	}

	return res
}

func targetExists(t string, list []string) bool {
	for _, i := range list {
		if t == i {
			return true
		}
	}
	return false
}
