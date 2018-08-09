package db

import (
	"database/sql"
	"github.com/textileio/textile-go/repo"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	ps "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"strconv"
	"sync"
	"time"
)

type PointersDB struct {
	modelStore
}

func NewPointerStore(db *sql.DB, lock *sync.Mutex) repo.PointerStore {
	return &PointersDB{modelStore{db, lock}}
}

func (p *PointersDB) Put(pointer repo.Pointer) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into pointers(id, key, address, cancelId, purpose, date) values(?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var cancelId string
	if pointer.CancelId != nil {
		cancelId = pointer.CancelId.Pretty()
	}
	_, err = stmt.Exec(pointer.Value.ID.Pretty(), pointer.Cid.String(), pointer.Value.Addrs[0].String(), cancelId, pointer.Purpose, int(time.Now().Unix()))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (p *PointersDB) Delete(id peer.ID) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, err := p.db.Exec("delete from pointers where id=?", id.Pretty())
	if err != nil {
		return err
	}
	return nil
}

func (p *PointersDB) DeleteAll(purpose repo.Purpose) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, err := p.db.Exec("delete from pointers where purpose=?", purpose)
	if err != nil {
		return err
	}
	return nil
}

func (p *PointersDB) GetAll() ([]repo.Pointer, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	stm := "select * from pointers"
	rows, err := p.db.Query(stm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []repo.Pointer
	for rows.Next() {
		var pointerId string
		var key string
		var address string
		var purpose int
		var date int
		var cancelId string
		if err := rows.Scan(&pointerId, &key, &address, &cancelId, &purpose, &date); err != nil {
			return ret, err
		}
		maAddr, err := ma.NewMultiaddr(address)
		if err != nil {
			return ret, err
		}
		pid, err := peer.IDB58Decode(pointerId)
		if err != nil {
			return ret, err
		}
		k, err := cid.Decode(key)
		if err != nil {
			return ret, err
		}
		var canId *peer.ID
		if cancelId != "" {
			c, err := peer.IDB58Decode(cancelId)
			if err != nil {
				return ret, err
			}
			canId = &c
		}
		pointer := repo.Pointer{
			Cid: k,
			Value: ps.PeerInfo{
				ID:    pid,
				Addrs: []ma.Multiaddr{maAddr},
			},
			CancelId: canId,
			Purpose:  repo.Purpose(purpose),
			Date:     time.Unix(int64(date), 0),
		}
		ret = append(ret, pointer)
	}
	return ret, nil
}

func (p *PointersDB) GetByPurpose(purpose repo.Purpose) ([]repo.Pointer, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	stm := "select * from pointers where purpose=" + strconv.Itoa(int(purpose))
	rows, err := p.db.Query(stm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []repo.Pointer
	for rows.Next() {
		var pointerId string
		var key string
		var address string
		var purpose int
		var date int
		var cancelId string
		if err := rows.Scan(&pointerId, &key, &address, &cancelId, &purpose, &date); err != nil {
			return ret, err
		}
		maAddr, err := ma.NewMultiaddr(address)
		if err != nil {
			return ret, err
		}
		pid, err := peer.IDB58Decode(pointerId)
		if err != nil {
			return ret, err
		}
		k, err := cid.Decode(key)
		if err != nil {
			return ret, err
		}
		var canID *peer.ID
		if cancelId != "" {
			c, err := peer.IDB58Decode(cancelId)
			if err != nil {
				return ret, err
			}
			canID = &c
		}
		pointer := repo.Pointer{
			Cid: k,
			Value: ps.PeerInfo{
				ID:    pid,
				Addrs: []ma.Multiaddr{maAddr},
			},
			CancelId: canID,
			Purpose:  repo.Purpose(purpose),
			Date:     time.Unix(int64(date), 0),
		}
		ret = append(ret, pointer)
	}
	return ret, nil
}

func (p *PointersDB) Get(id peer.ID) *repo.Pointer {
	p.lock.Lock()
	defer p.lock.Unlock()
	stm := "select * from pointers where id=?"
	row := p.db.QueryRow(stm, id.Pretty())

	var pointerId string
	var key string
	var address string
	var purpose int
	var date int
	var cancelId string
	if err := row.Scan(&pointerId, &key, &address, &cancelId, &purpose, &date); err != nil {
		return nil
	}
	maAddr, err := ma.NewMultiaddr(address)
	if err != nil {
		log.Errorf("error getting addr: %s", err)
		return nil
	}
	pid, err := peer.IDB58Decode(pointerId)
	if err != nil {
		log.Errorf("error getting id: %s", err)
		return nil
	}
	k, err := cid.Decode(key)
	if err != nil {
		log.Errorf("error decoding cid: %s", err)
		return nil
	}
	var canID *peer.ID
	if cancelId != "" {
		c, err := peer.IDB58Decode(cancelId)
		if err != nil {
			log.Errorf("error getting cancel id: %s", err)
			return nil
		}
		canID = &c
	}
	return &repo.Pointer{
		Cid: k,
		Value: ps.PeerInfo{
			ID:    pid,
			Addrs: []ma.Multiaddr{maAddr},
		},
		CancelId: canID,
		Purpose:  repo.Purpose(purpose),
		Date:     time.Unix(int64(date), 0),
	}
}
