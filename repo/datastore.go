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
	ThreadPeers() ThreadPeerStore
	AccountPeers() AccountPeerStore
	Blocks() BlockStore
	Notifications() NotificationStore
	CafeSessions() CafeSessionStore
	CafeRequests() CafeRequestStore
	CafeNonces() CafeNonceStore
	CafeClients() CafeClientStore
	CafeClientThreads() CafeClientThreadStore
	CafeClientMessages() CafeClientMessageStore
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
	List() []Thread
	Count(query string) int
	UpdateHead(id string, head string) error
	Delete(id string) error
}

type ThreadPeerStore interface {
	Queryable
	Add(peer *ThreadPeer) error
	List() []ThreadPeer
	ListById(id string) []ThreadPeer
	ListByThread(threadId string) []ThreadPeer
	Count(distinct bool) int
	Delete(id string, thread string) error
	DeleteById(id string) error
	DeleteByThread(thread string) error
}

type AccountPeerStore interface {
	Queryable
	Add(device *AccountPeer) error
	Get(id string) *AccountPeer
	List(query string) []AccountPeer
	Count(query string) int
	Delete(id string) error
}

type BlockStore interface {
	Queryable
	Add(block *Block) error
	Get(id string) *Block
	GetByDataId(dataId string) *Block
	List(offset string, limit int, query string) []Block
	Count(query string) int
	Delete(id string) error
	DeleteByThread(threadId string) error
}

type NotificationStore interface {
	Queryable
	Add(notification *Notification) error
	Get(id string) *Notification
	Read(id string) error
	ReadAll() error
	List(offset string, limit int) []Notification
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
	Add(req *CafeRequest) error
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

type CafeClientStore interface {
	Add(account *CafeClient) error
	Get(id string) *CafeClient
	Count() int
	List() []CafeClient
	ListByAddress(address string) []CafeClient
	UpdateLastSeen(id string, date time.Time) error
	Delete(id string) error
}

type CafeClientThreadStore interface {
	AddOrUpdate(thrd *CafeClientThread) error
	ListByClient(clientId string) []CafeClientThread
	Delete(id string, clientId string) error
	DeleteByClient(clientId string) error
}

type CafeClientMessageStore interface {
	AddOrUpdate(message *CafeClientMessage) error
	ListByClient(clientId string, offset string, limit int) []CafeClientMessage
	Read(id string, clientId string) error
	Delete(id string, clientId string) error
	DeleteByClient(clientId string) error
}
