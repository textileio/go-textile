package db

import (
	"database/sql"
	"sync"

	"github.com/textileio/textile-go/repo"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

type AlbumDB struct {
	modelStore
}

func NewAlbumStore(db *sql.DB, lock *sync.Mutex) repo.AlbumStore {
	return &AlbumDB{modelStore{db, lock}}
}

func (c *AlbumDB) Put(album *repo.PhotoAlbum) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into albums(id, key, mnemonic, name) values(?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}

	skb, err := album.Key.Bytes()
	if err != nil {
		log.Errorf("error getting key bytes: %s", err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(
		album.Id,
		skb,
		album.Mnemonic,
		album.Name,
	)
	if err != nil {
		tx.Rollback()
		log.Errorf("error in db exec: %s", err)
		return err
	}
	tx.Commit()
	return nil
}

func (c *AlbumDB) GetAlbum(id string) *repo.PhotoAlbum {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from albums where id='" + id + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *AlbumDB) GetAlbumByName(name string) *repo.PhotoAlbum {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := c.handleQuery("select * from albums where name='" + name + "';")
	if len(ret) == 0 {
		return nil
	}
	return &ret[0]
}

func (c *AlbumDB) GetAlbums(query string) []repo.PhotoAlbum {
	c.lock.Lock()
	defer c.lock.Unlock()
	q := ""
	if query != "" {
		q = " where " + query
	}
	return c.handleQuery("select * from albums" + q + ";")
}

func (c *AlbumDB) DeleteAlbum(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from albums where id=?", id)
	return err
}

func (c *AlbumDB) DeleteAlbumByName(name string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from albums where name=?", name)
	return err
}

func (c *AlbumDB) handleQuery(stm string) []repo.PhotoAlbum {
	var ret []repo.PhotoAlbum
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, mnemonic, name string
		var key []byte
		if err := rows.Scan(&id, &key, &mnemonic, &name); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		sk, err := libp2p.UnmarshalPrivateKey(key)
		if err != nil {
			log.Errorf("error unmarshaling private key: %s", err)
			continue
		}

		album := repo.PhotoAlbum{
			Id:       id,
			Key:      sk,
			Mnemonic: mnemonic,
			Name:     name,
		}
		ret = append(ret, album)
	}
	return ret
}
