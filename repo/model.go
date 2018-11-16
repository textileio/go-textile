package repo

import (
	"errors"
	"strings"
	"time"

	"github.com/textileio/textile-go/pb"
)

type Contact struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	Inboxes  []string  `json:"inboxes"`
	Added    time.Time `json:"added"`
}

type File struct {
	Mill     string                 `json:"mill"`
	Checksum string                 `json:"checksum"`
	Source   string                 `json:"source"`
	Hash     string                 `json:"hash"`
	Key      string                 `json:"key,omitempty"`
	Media    string                 `json:"media"`
	Name     string                 `json:"name"`
	Size     int                    `json:"size"`
	Added    time.Time              `json:"added"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

type Thread struct {
	Id        string      `json:"id"`
	Key       string      `json:"key"`
	PrivKey   []byte      `json:"sk"`
	Name      string      `json:"name"`
	Schema    string      `json:"schema"`
	Initiator string      `json:"initiator"`
	Type      ThreadType  `json:"type"`
	State     ThreadState `json:"state"`
	Head      string      `json:"head"`
}

type ThreadType int

// in order of decreasing privacy
const (
	PrivateThread  ThreadType = iota // invites not allowed
	ReadOnlyThread                   // all non-initiator writes ignored
	PublicThread                     // only non-initiator file writes ignored (annotations allowed)
	OpenThread                       // all writes allowed
)

func (tt ThreadType) Description() string {
	switch tt {
	case PrivateThread:
		return "PRIVATE"
	case ReadOnlyThread:
		return "READONLY"
	case PublicThread:
		return "PUBLIC"
	case OpenThread:
		return "OPEN"
	default:
		return "INVALID"
	}
}

func ThreadTypeFromString(desc string) (ThreadType, error) {
	switch strings.ToUpper(strings.TrimSpace(desc)) {
	case "PRIVATE":
		return PrivateThread, nil
	case "OPEN":
		return OpenThread, nil
	default:
		return -1, errors.New("could not parse thread type")
	}
}

type ThreadState int

const (
	ThreadLoading ThreadState = iota
	ThreadLoaded
)

func (ts ThreadState) Description() string {
	switch ts {
	case ThreadLoading:
		return "LOADING"
	case ThreadLoaded:
		return "LOADED"
	default:
		return "INVALID"
	}
}

type ThreadPeer struct {
	Id       string `json:"id"`
	ThreadId string `json:"thread_id"`
	Welcomed bool   `json:"welcomed"`
}

type ThreadMessage struct {
	Id       string       `json:"id"`
	PeerId   string       `json:"peer_id"`
	Envelope *pb.Envelope `json:"envelope"`
	Date     time.Time    `json:"date"`
}

type Block struct {
	Id       string    `json:"id"`
	ThreadId string    `json:"thread_id"`
	AuthorId string    `json:"author_id"`
	Type     BlockType `json:"type"`
	Date     time.Time `json:"date"`
	Parents  []string  `json:"parents"`
	Target   string    `json:"target,omitempty"`
	Body     string    `json:"body,omitempty"`
}

type BlockType int

const (
	MergeBlock BlockType = iota
	IgnoreBlock
	FlagBlock
	JoinBlock
	AnnounceBlock
	LeaveBlock
	MessageBlock
	FilesBlock
	CommentBlock
	LikeBlock
)

func (b BlockType) Description() string {
	switch b {
	case MergeBlock:
		return "MERGE"
	case IgnoreBlock:
		return "IGNORE"
	case FlagBlock:
		return "FLAG"
	case JoinBlock:
		return "JOIN"
	case AnnounceBlock:
		return "ANNOUNCE"
	case LeaveBlock:
		return "LEAVE"
	case MessageBlock:
		return "MESSAGE"
	case FilesBlock:
		return "FILES"
	case CommentBlock:
		return "COMMENT"
	case LikeBlock:
		return "LIKE"
	default:
		return "INVALID"
	}
}

type Notification struct {
	Id        string           `json:"id"`
	Date      time.Time        `json:"date"`
	ActorId   string           `json:"actor_id"`
	Subject   string           `json:"subject"`
	SubjectId string           `json:"subject_id"`
	BlockId   string           `json:"block_id,omitempty"`
	Target    string           `json:"target,omitempty"`
	Type      NotificationType `json:"type"`
	Body      string           `json:"body"`
	Read      bool             `json:"read"`
}

type NotificationType int

const (
	InviteReceivedNotification NotificationType = iota
	AccountPeerJoinedNotification
	PeerJoinedNotification
	PeerLeftNotification
	MessageAddedNotification
	FilesAddedNotification
	CommentAddedNotification
	LikeAddedNotification
)

func (n NotificationType) Description() string {
	switch n {
	case InviteReceivedNotification:
		return "INVITE_RECEIVED"
	case AccountPeerJoinedNotification:
		return "ACCOUNT_PEER_JOINED"
	case PeerJoinedNotification:
		return "PEER_JOINED"
	case PeerLeftNotification:
		return "PEER_LEFT"
	case MessageAddedNotification:
		return "MESSAGE_ADDED"
	case FilesAddedNotification:
		return "FILES_ADDED"
	case CommentAddedNotification:
		return "COMMENT_ADDED"
	case LikeAddedNotification:
		return "LIKE_ADDED"
	default:
		return "INVALID"
	}
}

type CafeSession struct {
	CafeId     string    `json:"cafe_id"`
	Access     string    `json:"access"`
	Refresh    string    `json:"refresh"`
	Expiry     time.Time `json:"expiry"`
	HttpAddr   string    `json:"http_addr"`
	SwarmAddrs []string  `json:"swarm_addrs"`
}

type CafeRequestType int

const (
	CafeStoreRequest CafeRequestType = iota
	CafeStoreThreadRequest
	CafePeerInboxRequest
)

func (rt CafeRequestType) Description() string {
	switch rt {
	case CafeStoreRequest:
		return "STORE"
	case CafeStoreThreadRequest:
		return "STORE_THREAD"
	case CafePeerInboxRequest:
		return "INBOX"
	default:
		return "INVALID"
	}
}

type CafeRequest struct {
	Id       string          `json:"id"`
	PeerId   string          `json:"peer_id"`
	TargetId string          `json:"target_id"`
	CafeId   string          `json:"cafe_id"`
	Type     CafeRequestType `json:"type"`
	Date     time.Time       `json:"date"`
}

type CafeMessage struct {
	Id     string    `json:"id"`
	PeerId string    `json:"peer_id"`
	Date   time.Time `json:"date"`
}

type CafeClientNonce struct {
	Value   string    `json:"value"`
	Address string    `json:"address"`
	Date    time.Time `json:"date"`
}

type CafeClient struct {
	Id       string    `json:"id"`
	Address  string    `json:"address"`
	Created  time.Time `json:"created"`
	LastSeen time.Time `json:"last_seen"`
}

type CafeClientThread struct {
	Id         string `json:"id"`
	ClientId   string `json:"client_id"`
	Ciphertext []byte `json:"ciphertext"`
}

type CafeClientMessage struct {
	Id       string    `json:"id"`
	PeerId   string    `json:"peer_id"`
	ClientId string    `json:"client_id"`
	Date     time.Time `json:"date"`
}
