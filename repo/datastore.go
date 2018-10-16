package repo

import (
	"database/sql"
	"github.com/textileio/textile-go/keypair"
	"time"
)

type Datastore interface {
	Config() ConfigStore
	Profile() ProfileStore
	Threads() ThreadStore
	Devices() DeviceStore
	Peers() PeerStore
	Blocks() BlockStore
	Notifications() NotificationStore
	CafeNonces() CafeNonceStore
	CafeAccounts() CafeAccountStore
	CafeSessions() CafeSessionStore
	CafeRequests() CafeRequestStore
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
	Init(pin string) error
	Configure(kp *keypair.Full, mobile bool, created time.Time) error
	GetAccount() (*keypair.Full, error)
	GetMobile() (bool, error)
	GetCreationDate() (time.Time, error)
	IsEncrypted() bool
}

type ProfileStore interface {
	SetUsername(username string) error
	GetUsername() (*string, error)
	SetAvatar(uri string) error
	GetAvatar() (*string, error)
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

type NotificationStore interface {
	Queryable
	Add(notification *Notification) error
	Get(id string) *Notification
	Read(id string) error
	ReadAll() error
	List(offset string, limit int, query string) []Notification
	CountUnread() int
	Delete(id string) error
	DeleteByActorId(actorId string) error
	DeleteBySubjectId(subjectId string) error
	DeleteByBlockId(blockId string) error
}

// Cafe user-side stores

type CafeSessionStore interface {
	AddOrUpdate(session *CafeSession) error
	Get(cafeId string) *CafeSession
	List() []CafeSession
	Delete(cafeId string) error
}

type CafeRequestStore interface {
	Queryable
	Put(req *CafeRequest) error
	List(offset string, limit int) []CafeRequest
	Delete(id string) error
	DeleteByCafe(cafeId string) error
}

// Cafe host-side stores

type CafeNonceStore interface {
	Add(nonce *CafeNonce) error
	Get(value string) *CafeNonce
	Delete(value string) error
}

type CafeAccountStore interface {
	Add(account *CafeAccount) error
	Get(id string) *CafeAccount
	Count() int
	ListByAddress(address string) []CafeAccount
	UpdateLastSeen(id string, date time.Time) error
	Delete(id string) error
}
