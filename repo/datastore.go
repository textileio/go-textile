package repo

import (
	"database/sql"
	"github.com/textileio/textile-go/keypair"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
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
	Init(pin string) error
	Configure(kp *keypair.Full, mobile bool, created time.Time) error
	GetAccount() (*keypair.Full, error)
	GetMobile() (bool, error)
	GetCreationDate() (time.Time, error)
	IsEncrypted() bool
}

type ProfileStore interface {
	CafeLogin(tokens *CafeTokens) error
	CafeLogout() error
	GetCafeTokens() (tokens *CafeTokens, err error)
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

// Cafe stores

type NonceStore interface {
	Add(nonce *Nonce) error
	Get(value string) *Nonce
	Delete(value string) error
}

type AccountStore interface {
	Add(account *Account) error
	Get(id string) *Account
	Count() int
	ListByAddress(address string) []Account
	UpdateLastSeen(id string, date time.Time) error
	Delete(id string) error
}
