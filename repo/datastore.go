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
	PinRequests() PinRequestStore
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
	SignIn(username string, tokens *CafeTokens) error
	SignOut() error
	GetUsername() (string, error)
	SetAvatarId(id string) error
	GetAvatarId() (string, error)
	GetTokens() (tokens *CafeTokens, err error)
}

type ThreadStore interface {
	Queryable
	Add(thread *Thread) error
	Get(id string) *Thread
	List(query string) []Thread
	Count(query string) int
	UpdateHead(id string, head string) error
	Delete(id string) error
}

type DeviceStore interface {
	Queryable
	Add(device *Device) error
	Get(id string) *Device
	List(query string) []Device
	Count(query string) int
	Delete(id string) error
}

type PeerStore interface {
	Queryable
	Add(peer *Peer) error
	Get(row string) *Peer
	GetById(id string) *Peer
	List(offset string, limit int, query string) []Peer
	Count(query string, distinct bool) int
	Delete(id string, thread string) error
	DeleteByThreadId(thread string) error
}

type BlockStore interface {
	Queryable
	Add(block *Block) error
	Get(id string) *Block
	GetByDataId(dataId string) *Block
	List(offset string, limit int, query string) []Block
	Count(query string) int
	Delete(id string) error
	DeleteByThreadId(threadId string) error
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
	Get(id peer.ID) *Pointer
	GetByPurpose(purpose Purpose) ([]Pointer, error)
	GetAll() ([]Pointer, error)
}

type PinRequestStore interface {
	Queryable
	Put(pr *PinRequest) error
	List(offset string, limit int) []PinRequest
	Delete(id string) error
}
