package repo

import (
	"database/sql"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"time"
)

type Datastore interface {
	Config() ConfigStore
	Profile() ProfileStore
	Threads() ThreadStore
	Devices() DeviceStore
	Peers() PeerStore
	Blocks() BlockStore
	OfflineMessages() OfflineMessageStore
	Pointers() PointerStore
	Ping() error
	Close()
}

type Queryable interface {
	BeginTransaction() (*sql.Tx, error)
	PrepareQuery(string) (*sql.Stmt, error)
	PrepareAndExecuteQuery(string, ...interface{}) (*sql.Rows, error)
	ExecuteQuery(string, ...interface{}) (sql.Result, error)
}

type ConfigStore interface {
	Init(password string) error
	Configure(created time.Time) error
	GetCreationDate() (time.Time, error)
	IsEncrypted() bool
}

type ProfileStore interface {
	SignIn(username string, accessToken string, refreshToken string) error
	SignOut() error
	GetUsername() (string, error)
	GetTokens() (accessToken string, refreshToken string, err error)
}

type ThreadStore interface {
	Queryable
	Add(thread *Thread) error
	Get(id string) *Thread
	GetByName(name string) *Thread
	List(query string) []Thread
	UpdateHead(id string, head string) error
	Delete(id string) error
	DeleteByName(name string) error
}

type DeviceStore interface {
	Queryable
	Add(device *Device) error
	Get(id string) *Device
	GetByName(name string) *Device
	List(query string) []Device
	Delete(id string) error
	DeleteByName(name string) error
}

type PeerStore interface {
	Queryable
	Add(peer *Peer) error
	Get(row string) *Peer
	GetByPubKey(pk string) *Peer
	List(offset string, limit int, query string) []Peer
	Delete(row string) error
}

type BlockStore interface {
	Queryable
	Add(block *Block) error
	Get(id string) *Block
	GetByTarget(target string) *Block
	List(offset string, limit int, query string) []Block
	Delete(id string) error
}

type OfflineMessageStore interface {
	Queryable
	Put(url string) error
	Has(url string) bool
	SetMessage(url string, message []byte) error
	GetMessages() (map[string][]byte, error)
	DeleteMessage(url string) error
}

type PointerStore interface {
	Queryable
	Put(p Pointer) error
	Delete(id peer.ID) error
	DeleteAll(purpose Purpose) error
	Get(id peer.ID) (Pointer, error)
	GetByPurpose(purpose Purpose) ([]Pointer, error)
	GetAll() ([]Pointer, error)
}
