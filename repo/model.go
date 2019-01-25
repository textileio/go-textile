package repo

import (
	"errors"
	"strings"
	"time"

	"github.com/textileio/textile-go/pb"
)

type File struct {
	Mill     string                 `json:"mill"`
	Checksum string                 `json:"checksum"`
	Source   string                 `json:"source"`
	Opts     string                 `json:"opts,omitempty"`
	Hash     string                 `json:"hash"`
	Key      string                 `json:"key,omitempty"`
	Media    string                 `json:"media"`
	Name     string                 `json:"name,omitempty"`
	Size     int                    `json:"size"`
	Added    time.Time              `json:"added"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
	Targets  []string               `json:"targets,omitempty"`
}

type Thread struct {
	Id        string        `json:"id"`
	Key       string        `json:"key"`
	PrivKey   []byte        `json:"sk"`
	Name      string        `json:"name"`
	Schema    string        `json:"schema"`
	Initiator string        `json:"initiator"`
	Type      ThreadType    `json:"type"`
	Sharing   ThreadSharing `json:"sharing"`
	Members   []string      `json:"members"` // if empty, _everyone_ is a member
	State     ThreadState   `json:"state"`
	Head      string        `json:"head"`
}

// ThreadType controls read (R), annotate (A), and write (W) access
type ThreadType int

// in order of decreasing privacy
const (
	PrivateThread  ThreadType = iota // initiator: RAW, members:
	ReadOnlyThread                   // initiator: RAW, members: R
	PublicThread                     // initiator: RAW, members: RA
	OpenThread                       // initiator: RAW, members: RAW
)

func (tt ThreadType) Description() string {
	switch tt {
	case PrivateThread:
		return "PRIVATE"
	case ReadOnlyThread:
		return "READ_ONLY"
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
	case "READ_ONLY":
		return ReadOnlyThread, nil
	case "PUBLIC":
		return PublicThread, nil
	case "OPEN":
		return OpenThread, nil
	default:
		return -1, errors.New("could not parse thread type")
	}
}

// ThreadSharing controls if (Y/N) a thread can be shared
type ThreadSharing int

const (
	NotSharedThread  ThreadSharing = iota // initiator: N, members: N
	InviteOnlyThread                      // initiator: Y, members: N
	SharedThread                          // initiator: Y, members: Y
)

func (ts ThreadSharing) Description() string {
	switch ts {
	case NotSharedThread:
		return "NOT_SHARED"
	case InviteOnlyThread:
		return "INVITE_ONLY"
	case SharedThread:
		return "SHARED"
	default:
		return "INVALID"
	}
}

func ThreadSharingFromString(desc string) (ThreadSharing, error) {
	switch strings.ToUpper(strings.TrimSpace(desc)) {
	case "NOT_SHARED":
		return NotSharedThread, nil
	case "INVITE_ONLY":
		return InviteOnlyThread, nil
	case "SHARED":
		return SharedThread, nil
	default:
		return -1, errors.New("could not parse thread sharing")
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

type ThreadInvite struct {
	Id      string    `json:"id"`
	Block   []byte    `json:"block"`
	Name    string    `json:"name"`
	Contact *Contact  `json:"contact"`
	Date    time.Time `json:"date"`
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

type Contact struct {
	Id       string    `json:"id"`
	Address  string    `json:"address"`
	Username string    `json:"username,omitempty"`
	Avatar   string    `json:"avatar,omitempty"`
	Inboxes  []Cafe    `json:"inboxes,omitempty"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
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

type Cafe struct {
	Peer     string   `json:"peer"`
	Address  string   `json:"address"`
	API      string   `json:"api"`
	Protocol string   `json:"protocol"`
	Node     string   `json:"node"`
	URL      string   `json:"url"`
	Swarm    []string `json:"swarm"`
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
	Cafe     pb.Cafe         `json:"cafe"`
	Type     CafeRequestType `json:"type"`
	Date     time.Time       `json:"date"`
}

type CafeMessage struct {
	Id       string    `json:"id"`
	PeerId   string    `json:"peer_id"`
	Date     time.Time `json:"date"`
	Attempts int       `json:"attempts"`
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
