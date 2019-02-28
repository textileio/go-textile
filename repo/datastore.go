package repo

import (
	"database/sql"
	"strings"
	"time"

	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

type Datastore interface {
	Config() ConfigStore
	Contacts() ContactStore
	Files() FileStore
	Threads() ThreadStore
	ThreadPeers() ThreadPeerStore
	Blocks() BlockStore
	BlockMessages() BlockMessageStore
	Invites() InviteStore
	Notifications() NotificationStore
	CafeSessions() CafeSessionStore
	CafeRequests() CafeRequestStore
	CafeMessages() CafeMessageStore
	CafeClientNonces() CafeClientNonceStore
	CafeClients() CafeClientStore
	CafeTokens() CafeTokenStore
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
	Configure(accnt *keypair.Full, created time.Time) error
	GetAccount() (*keypair.Full, error)
	GetCreationDate() (time.Time, error)
	IsEncrypted() bool
}

type ContactStore interface {
	Queryable
	Add(contact *pb.Contact) error
	AddOrUpdate(contact *pb.Contact) error
	Get(id string) *pb.Contact
	GetBest(id string) *pb.Contact
	List(query string) *pb.ContactList
	Find(id string, address string, username string, exclude []string) *pb.ContactList
	Count(query string) int
	UpdateUsername(id string, username string) error
	UpdateAvatar(id string, avatar string) error
	UpdateInboxes(id string, inboxes []*pb.Cafe) error
	Delete(id string) error
}

type FileStore interface {
	Queryable
	Add(file *pb.FileIndex) error
	Get(hash string) *pb.FileIndex
	GetByPrimary(mill string, checksum string) *pb.FileIndex
	GetBySource(mill string, source string, opts string) *pb.FileIndex
	AddTarget(hash string, target string) error
	RemoveTarget(hash string, target string) error
	Count() int
	Delete(hash string) error
}

type ThreadStore interface {
	Queryable
	Add(thread *pb.Thread) error
	Get(id string) *pb.Thread
	GetByKey(key string) *pb.Thread
	List() *pb.ThreadList
	Count() int
	UpdateHead(id string, head string) error
	Delete(id string) error
}

type ThreadPeerStore interface {
	Queryable
	Add(peer *pb.ThreadPeer) error
	List() []pb.ThreadPeer
	ListById(id string) []pb.ThreadPeer
	ListByThread(threadId string) []pb.ThreadPeer
	ListUnwelcomedByThread(threadId string) []pb.ThreadPeer
	WelcomeByThread(thread string) error
	Count(distinct bool) int
	Delete(id string, thread string) error
	DeleteById(id string) error
	DeleteByThread(thread string) error
}

type BlockStore interface {
	Queryable
	Add(block *pb.Block) error
	Get(id string) *pb.Block
	List(offset string, limit int, query string) *pb.BlockList
	Count(query string) int
	Delete(id string) error
	DeleteByThread(threadId string) error
}

type BlockMessageStore interface {
	Queryable
	Add(msg *pb.BlockMessage) error
	List(offset string, limit int) []pb.BlockMessage
	Delete(id string) error
}

type InviteStore interface {
	Queryable
	Add(invite *pb.Invite) error
	Get(id string) *pb.Invite
	List() *pb.InviteList
	Delete(id string) error
}

type NotificationStore interface {
	Queryable
	Add(notification *pb.Notification) error
	Get(id string) *pb.Notification
	Read(id string) error
	ReadAll() error
	List(offset string, limit int) *pb.NotificationList
	CountUnread() int
	Delete(id string) error
	DeleteByActor(actorId string) error
	DeleteBySubject(subjectId string) error
	DeleteByBlock(blockId string) error
}

// Cafe user-side stores

type CafeSessionStore interface {
	AddOrUpdate(session *pb.CafeSession) error
	Get(cafeId string) *pb.CafeSession
	List() *pb.CafeSessionList
	Delete(cafeId string) error
}

type CafeRequestStore interface {
	Queryable
	Add(req *pb.CafeRequest) error
	List(offset string, limit int) []pb.CafeRequest
	Delete(id string) error
	DeleteByCafe(cafeId string) error
}

type CafeMessageStore interface {
	Queryable
	Add(msg *pb.CafeMessage) error
	List(offset string, limit int) []pb.CafeMessage
	AddAttempt(id string) error
	Delete(id string) error
}

// Cafe host-side stores

type CafeClientNonceStore interface {
	Add(nonce *pb.CafeClientNonce) error
	Get(value string) *pb.CafeClientNonce
	Delete(value string) error
}

type CafeClientStore interface {
	Add(account *pb.CafeClient) error
	Get(id string) *pb.CafeClient
	Count() int
	List() []pb.CafeClient
	ListByAddress(address string) []pb.CafeClient
	UpdateLastSeen(id string, date time.Time) error
	Delete(id string) error
}

type CafeClientThreadStore interface {
	AddOrUpdate(thrd *pb.CafeClientThread) error
	ListByClient(clientId string) []pb.CafeClientThread
	Delete(id string, clientId string) error
	DeleteByClient(clientId string) error
}

type CafeClientMessageStore interface {
	AddOrUpdate(message *pb.CafeClientMessage) error
	ListByClient(clientId string, limit int) []pb.CafeClientMessage
	CountByClient(clientId string) int
	Delete(id string, clientId string) error
	DeleteByClient(clientId string, limit int) error
}

type CafeTokenStore interface {
	Add(token *pb.CafeToken) error
	Get(id string) *pb.CafeToken
	List() []pb.CafeToken
	Delete(id string) error
}

func ConflictError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
