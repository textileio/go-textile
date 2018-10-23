package repo

import (
	"github.com/textileio/textile-go/photo"
	"time"
)

type Contact struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	Added    time.Time `json:"added"`
}

type Thread struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	PrivKey []byte `json:"sk"`
	Head    string `json:"head"`
}

type ThreadPeer struct {
	Id       string `json:"id"`
	ThreadId string `json:"thread_id"`
	Welcomed bool   `json:"welcomed"`
}

type Block struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	Parents  []string  `json:"parents"`
	ThreadId string    `json:"thread_id"`
	AuthorId string    `json:"author_id"`
	Type     BlockType `json:"type"`

	DataId       string          `json:"data_id"`
	DataKey      string          `json:"data_key"`
	DataCaption  string          `json:"data_caption"`
	DataMetadata *photo.Metadata `json:"data_metadata"`
}

type DataBlockConfig struct {
	DataId       string          `json:"data_id"`
	DataKey      string          `json:"data_key"`
	DataCaption  string          `json:"data_caption"`
	DataMetadata *photo.Metadata `json:"data_metadata"`
}

type BlockType int

const (
	MergeBlock BlockType = iota
	IgnoreBlock
	JoinBlock
	AnnounceBlock
	LeaveBlock
	PhotoBlock
	CommentBlock
	LikeBlock
)

func (b BlockType) Description() string {
	switch b {
	case JoinBlock:
		return "JOIN"
	case LeaveBlock:
		return "LEAVE"
	case PhotoBlock:
		return "PHOTO"
	case CommentBlock:
		return "COMMENT"
	case LikeBlock:
		return "LIKE"
	case IgnoreBlock:
		return "IGNORE"
	case MergeBlock:
		return "MERGE"
	default:
		return "INVALID"
	}
}

type Notification struct {
	Id            string           `json:"id"`
	Date          time.Time        `json:"date"`
	ActorId       string           `json:"actor_id"`                 // peer id
	ActorUsername string           `json:"actor_username,omitempty"` // peer username
	Subject       string           `json:"subject"`                  // thread name | device name
	SubjectId     string           `json:"subject_id"`               // thread id | device id
	BlockId       string           `json:"block_id,omitempty"`       // block id
	DataId        string           `json:"data_id,omitempty"`        // photo id, etc.
	Type          NotificationType `json:"type"`
	Body          string           `json:"body"`
	Read          bool             `json:"read"`
}

type NotificationType int

const (
	ReceivedInviteNotification   NotificationType = iota // peerA invited you
	AccountPeerAddedNotification                         // new account peer added
	PhotoAddedNotification                               // peerA added a photo
	CommentAddedNotification                             // peerA commented on peerB's photo, video, comment, etc.
	LikeAddedNotification                                // peerA liked peerB's photo, video, comment, etc.
	PeerJoinedNotification                               // peerA joined
	PeerLeftNotification                                 // peerA left
	TextAddedNotification                                // peerA added a message
)

func (n NotificationType) Description() string {
	switch n {
	case ReceivedInviteNotification:
		return "RECEIVED_INVITE"
	case AccountPeerAddedNotification:
		return "ACCOUNT_PEER_ADDED"
	case PhotoAddedNotification:
		return "PHOTO_ADDED"
	case CommentAddedNotification:
		return "COMMENT_ADDED"
	case LikeAddedNotification:
		return "LIKE_ADDED"
	case PeerJoinedNotification:
		return "PEER_JOINED"
	case PeerLeftNotification:
		return "PEER_LEFT"
	default:
		return "INVALID"
	}
}

type CafeSession struct {
	CafeId  string    `json:"cafe_id"`
	Access  string    `json:"access"`
	Refresh string    `json:"refresh"`
	Expiry  time.Time `json:"expiry"`
}

type CafeRequestType int

const (
	CafeStoreRequest CafeRequestType = iota
	CafeStoreThreadRequest
	CafeInboxRequest
)

func (rt CafeRequestType) Description() string {
	switch rt {
	case CafeStoreRequest:
		return "STORE"
	case CafeStoreThreadRequest:
		return "STORE_THREAD"
	case CafeInboxRequest:
		return "INBOX"
	default:
		return "INVALID"
	}
}

type CafeRequest struct {
	Id       string          `json:"id"`
	TargetId string          `json:"target_id"`
	CafeId   string          `json:"cafe_id"`
	Type     CafeRequestType `json:"type"`
	Date     time.Time       `json:"date"`
}

type CafeInbox struct {
	PeerId string `json:"peer_id"`
	CafeId string `json:"cafe_id"`
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
	SkCipher   []byte `json:"sk_cipher"`
	HeadCipher []byte `json:"head_cipher"`
	NameCipher []byte `json:"name_cipher"`
}

type CafeClientMessage struct {
	Id       string    `json:"id"`
	ClientId string    `json:"client_id"`
	Date     time.Time `json:"date"`
	Read     bool      `json:"read"`
}
