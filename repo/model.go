package repo

import (
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
	TokenId  string    `json:"token_id,omitempty"`
}

type CafeToken struct {
	Id    string    `json:"id"`
	Token []byte    `json:"token"`
	Date  time.Time `json:"date"`
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
