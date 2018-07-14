package db

import (
	"database/sql"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/repo"
	"path"
	"sync"
)

var log = logging.MustGetLogger("db")

type SQLiteDatastore struct {
	config          repo.ConfigStore
	profile         repo.ProfileStore
	threads         repo.ThreadStore
	devices         repo.DeviceStore
	peers           repo.PeerStore
	blocks          repo.BlockStore
	offlineMessages repo.OfflineMessageStore
	pointers        repo.PointerStore
	db              *sql.DB
	lock            *sync.Mutex
}

func Create(repoPath, password string) (*SQLiteDatastore, error) {
	var dbPath string
	dbPath = path.Join(repoPath, "datastore", "mainnet.db")
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if password != "" {
		p := "pragma key='" + password + "';"
		conn.Exec(p)
	}
	mux := new(sync.Mutex)
	sqliteDB := &SQLiteDatastore{
		config:          NewConfigStore(conn, mux, dbPath),
		profile:         NewProfileStore(conn, mux),
		threads:         NewThreadStore(conn, mux),
		devices:         NewDeviceStore(conn, mux),
		peers:           NewPeerStore(conn, mux),
		blocks:          NewBlockStore(conn, mux),
		offlineMessages: NewOfflineMessageStore(conn, mux),
		pointers:        NewPointerStore(conn, mux),
		db:              conn,
		lock:            mux,
	}

	return sqliteDB, nil
}

func (d *SQLiteDatastore) Ping() error {
	return d.db.Ping()
}

func (d *SQLiteDatastore) Close() {
	d.db.Close()
}

func (d *SQLiteDatastore) Config() repo.ConfigStore {
	return d.config
}

func (d *SQLiteDatastore) Profile() repo.ProfileStore {
	return d.profile
}

func (d *SQLiteDatastore) Threads() repo.ThreadStore {
	return d.threads
}

func (d *SQLiteDatastore) Devices() repo.DeviceStore {
	return d.devices
}

func (d *SQLiteDatastore) Peers() repo.PeerStore {
	return d.peers
}

func (d *SQLiteDatastore) Blocks() repo.BlockStore {
	return d.blocks
}

func (d *SQLiteDatastore) OfflineMessages() repo.OfflineMessageStore {
	return d.offlineMessages
}

func (d *SQLiteDatastore) Pointers() repo.PointerStore {
	return d.pointers
}

func (d *SQLiteDatastore) Copy(dbPath string, password string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	var cp string
	stmt := "select name from sqlite_master where type='table'"
	rows, err := d.db.Query(stmt)
	if err != nil {
		log.Errorf("error in copy: %s", err)
		return err
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables = append(tables, name)
	}
	if password == "" {
		cp = `attach database '` + dbPath + `' as plaintext key '';`
		for _, name := range tables {
			cp = cp + "insert into plaintext." + name + " select * from main." + name + ";"
		}
	} else {
		cp = `attach database '` + dbPath + `' as encrypted key '` + password + `';`
		for _, name := range tables {
			cp = cp + "insert into encrypted." + name + " select * from main." + name + ";"
		}
	}

	_, err = d.db.Exec(cp)
	if err != nil {
		return err
	}

	return nil
}

func (d *SQLiteDatastore) InitTables(password string) error {
	return initDatabaseTables(d.db, password)
}

func initDatabaseTables(db *sql.DB, password string) error {
	var sqlStmt string
	if password != "" {
		sqlStmt = "PRAGMA key = '" + password + "';"
	}
	sqlStmt += `
	create table config (key text primary key not null, value blob);
    create table profile (key text primary key not null, value blob);
    create table threads (id text primary key not null, name text not null, sk blob not null, head text not null);
    create unique index thread_name on threads (name);
    create table devices (id text primary key not null, name text not null);
    create unique index device_name on devices (name);
    create table peers (row text primary key not null, id text not null, thread text not null, pk blob not null);
    create unique index peer_thread_id on peers (thread, id);
    create table blocks (id text primary key not null, target text not null, parents text not null, key blob not null, pk text not null, ppk text not null, type integer not null, date integer not null);
    create index block_target on blocks (target);
    create index block_pk_type_date on blocks (pk, type, date);
    create table offlinemessages (url text primary key not null, date integer, message blob);
	create table pointers (id text primary key not null, key text, address text, cancelId text, purpose integer, date integer);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}
