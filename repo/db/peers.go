package db

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
)

type PeerDB struct {
	modelStore
}

func NewPeerStore(db *sql.DB, lock *sync.Mutex) repo.PeerStore {
	return &PeerDB{modelStore{db, lock}}
}

func (c *PeerDB) Add(peer *pb.Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into peers(id, address, username, avatar, inboxes, created, updated) values(?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	inboxes, err := json.Marshal(peer.Inboxes)
	if err != nil {
		return err
	}

	var created, updated int64
	if peer.Created == nil {
		created = time.Now().UnixNano()
	} else {
		created = util.ProtoNanos(peer.Created)
	}
	if peer.Updated == nil {
		updated = time.Now().UnixNano()
	} else {
		updated = util.ProtoNanos(peer.Updated)
	}

	_, err = stmt.Exec(
		peer.Id,
		peer.Address,
		peer.Name,
		peer.Avatar,
		inboxes,
		created,
		updated,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *PeerDB) AddOrUpdate(peer *pb.Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert or replace into peers(id, address, username, avatar, inboxes, created, updated) values(?,?,?,?,?,coalesce((select created from peers where id=?),?),?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	inboxes, err := json.Marshal(peer.Inboxes)
	if err != nil {
		return err
	}

	var created int64
	if peer.Created == nil {
		created = time.Now().UnixNano()
	} else {
		created = util.ProtoNanos(peer.Created)
	}

	_, err = stmt.Exec(
		peer.Id,
		peer.Address,
		peer.Name,
		peer.Avatar,
		inboxes,
		peer.Id,
		created,
		time.Now().UnixNano(),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (c *PeerDB) Get(id string) *pb.Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := c.handleQuery("select * from peers where id='" + id + "';")
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

func (c *PeerDB) GetBestUser(id string) *pb.User {
	c.lock.Lock()
	defer c.lock.Unlock()
	stm := "select username, avatar, (select address from peers where id='" + id + "') as addr from peers where address=addr order by updated desc;"
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	var latest *pb.User
	var i int
	for rows.Next() {
		var name, avatar, addr string
		if err := rows.Scan(&name, &avatar, &addr); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		if i == 0 {
			latest = &pb.User{Address: addr, Name: name, Avatar: avatar}
		} else if latest != nil {
			if name != "" && latest.Name == "" {
				latest.Name = name
			}
			if avatar != "" && latest.Avatar == "" {
				latest.Avatar = avatar
			}
			if latest.Name != "" && latest.Avatar != "" {
				_ = rows.Close()
				return ensureName(latest)
			}
		}
		i++
	}
	return ensureName(latest)
}

func (c *PeerDB) List(query string) []*pb.Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := "select * from peers"
	if query != "" {
		q += " where " + query
	}
	q += " order by updated desc;"
	return c.handleQuery(q)
}

func (c *PeerDB) Find(address string, name string, exclude []string) []*pb.Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	if address == "" && name == "" {
		return nil
	}
	var q string
	if address != "" {
		q += "address='" + address + "'"
	}
	if name != "" {
		if len(q) > 0 {
			q += " and "
		}
		q += "username like '%" + name + "%'"
	}
	if len(exclude) > 0 {
		q += " and id not in ("
		for i, e := range exclude {
			q += "'" + e + "'"
			if i != len(exclude)-1 {
				q += ","
			}
		}
		q += ")"
	}
	return c.handleQuery("select * from peers where " + q + " order by updated desc;")
}

func (c *PeerDB) Count(query string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := "select Count(*) from peers"
	if query != "" {
		q += " where " + query
	}
	q += ";"
	row := c.db.QueryRow(q)
	var count int
	_ = row.Scan(&count)
	return count
}

func (c *PeerDB) UpdateName(id string, name string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update peers set username=?, updated=? where id=?", name, time.Now().UnixNano(), id)
	return err
}

func (c *PeerDB) UpdateAvatar(id string, avatar string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("update peers set avatar=?, updated=? where id=?", avatar, time.Now().UnixNano(), id)
	return err
}

func (c *PeerDB) UpdateInboxes(id string, inboxes []*pb.Cafe) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	inboxesb, err := json.Marshal(inboxes)
	if err != nil {
		return err
	}
	_, err = c.db.Exec("update peers set inboxes=?, updated=? where id=?", inboxesb, time.Now().UnixNano(), id)
	return err
}

func (c *PeerDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from peers where id=?", id)
	return err
}

func (c *PeerDB) DeleteByAddress(address string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from peers where address=?", address)
	return err
}

func (c *PeerDB) handleQuery(stm string) []*pb.Peer {
	list := make([]*pb.Peer, 0)
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return list
	}
	for rows.Next() {
		var id, address, name, avatar string
		var inboxes []byte
		var createdInt, updatedInt int64
		if err := rows.Scan(&id, &address, &name, &avatar, &inboxes, &createdInt, &updatedInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}
		row := c.handleRow(id, address, name, avatar, inboxes, createdInt, updatedInt)
		if row != nil {
			list = append(list, row)
		}
	}
	return list
}

func (c *PeerDB) handleRow(id string, address string, name string, avatar string, inboxes []byte, createdInt int64, updatedInt int64) *pb.Peer {
	cafes := make([]*pb.Cafe, 0)
	if err := json.Unmarshal(inboxes, &cafes); err != nil {
		log.Errorf("error unmarshaling cafes: %s", err)
		return nil
	}

	return &pb.Peer{
		Id:      id,
		Address: address,
		Name:    name,
		Avatar:  avatar,
		Inboxes: cafes,
		Created: util.ProtoTs(createdInt),
		Updated: util.ProtoTs(updatedInt),
	}
}

func ensureName(user *pb.User) *pb.User {
	if user == nil || user.Address == "" || user.Name != "" {
		return user
	}
	if len(user.Address) >= 7 {
		user.Name = user.Address[:7]
	}
	return user
}
