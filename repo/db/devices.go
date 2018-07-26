package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	"sync"
)

type DeviceDB struct {
	modelStore
}

func NewDeviceStore(db *sql.DB, lock *sync.Mutex) repo.DeviceStore {
	return &DeviceDB{modelStore{db, lock}}
}

func (c *DeviceDB) Add(device *repo.Device) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into devices(id, name) values(?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		device.Id,
		device.Name,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *DeviceDB) Get(id string) *repo.Device {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from devices where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *DeviceDB) List(query string) []repo.Device {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := ""
	if query != "" {
		q = " where " + query
	}
	return c.handleQuery("select * from devices" + q + ";")
}

func (c *DeviceDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from devices where id=?", id)
	return err
}

func (c *DeviceDB) handleQuery(stm string) []repo.Device {
	var ret []repo.Device
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		device := repo.Device{
			Id:   id,
			Name: name,
		}
		ret = append(ret, device)
	}
	return ret
}
