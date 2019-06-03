package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type FileDB struct {
	modelStore
}

func NewFileStore(db *sql.DB, lock *sync.Mutex) repo.FileStore {
	return &FileDB{modelStore{db, lock}}
}

func (c *FileDB) Add(file *pb.FileIndex) error {
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

	var meta string
	if file.Meta != nil {
		var err error
		meta, err = pbMarshaler.MarshalToString(file.Meta)
		if err != nil {
			return err
		}
	}

	var targets *string
	if len(file.Targets) > 0 {
		tmp := strings.Join(file.Targets, ",")
		targets = &tmp
	}

	var added int64
	if file.Added == nil {
		added = time.Now().UnixNano()
	} else {
		added = util.ProtoNanos(file.Added)
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
		added,
		[]byte(meta),
		targets,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *FileDB) Get(hash string) *pb.FileIndex {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from files where hash='" + hash + "';")
	if len(res) == 0 {
		return nil
	}
	return &res[0]
}

func (c *FileDB) GetByPrimary(mill string, checksum string) *pb.FileIndex {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from files where mill='" + mill + "' and checksum='" + checksum + "';")
	if len(res) == 0 {
		return nil
	}
	return &res[0]
}

func (c *FileDB) GetBySource(mill string, source string, opts string) *pb.FileIndex {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from files where mill='" + mill + "' and source='" + source + "' and opts='" + opts + "';")
	if len(res) == 0 {
		return nil
	}
	return &res[0]
}

func (c *FileDB) AddTarget(hash string, target string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	res := c.handleTargetsQuery("select targets from files where hash='" + hash + "';")
	if len(res) == 0 {
		return fmt.Errorf("file not found")
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
		return fmt.Errorf("file not found")
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
	_ = row.Scan(&count)
	return count
}

func (c *FileDB) Delete(hash string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from files where hash=?", hash)
	return err
}

func (c *FileDB) handleQuery(stm string) []pb.FileIndex {
	var list []pb.FileIndex
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var mill, checksum, source, opts, hash, key, media, name string
		var size int64
		var addedInt int64
		var metab []byte
		var targets *string

		if err := rows.Scan(&mill, &checksum, &source, &opts, &hash, &key, &media, &name, &size, &addedInt, &metab, &targets); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		meta := &structpb.Struct{}
		if metab != nil {
			if err := jsonpb.Unmarshal(bytes.NewReader(metab), meta); err != nil {
				log.Errorf("failed to unmarshal file meta: %s", err)
				continue
			}
		}

		tlist := make([]string, 0)
		if targets != nil {
			tlist = util.SplitString(*targets, ",")
		}

		list = append(list, pb.FileIndex{
			Mill:     mill,
			Checksum: checksum,
			Source:   source,
			Opts:     opts,
			Hash:     hash,
			Key:      key,
			Media:    media,
			Name:     name,
			Size:     size,
			Added:    util.ProtoTs(addedInt),
			Meta:     meta,
			Targets:  tlist,
		})
	}

	return list
}

func (c *FileDB) handleTargetsQuery(stm string) [][]string {
	var list [][]string
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
			tlist = util.SplitString(*targets, ",")
		}

		list = append(list, tlist)
	}

	return list
}

func targetExists(t string, list []string) bool {
	for _, i := range list {
		if t == i {
			return true
		}
	}
	return false
}
