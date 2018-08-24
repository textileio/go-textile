package repo

import (
	"time"
)

type Thread struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	PrivKey []byte `json:"sk"`
	Head    string `json:"head"`
}

type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Peer struct {
	Row      string `json:"row"`
	Id       string `json:"id"`
	PubKey   []byte `json:"pk"`
	ThreadId string `json:"thread_id"`
}

type Block struct {
	Id                   string    `json:"id"`
	Date                 time.Time `json:"date"`
	Parents              []string  `json:"parents"`
	ThreadId             string    `json:"thread_id"`
	AuthorPk             string    `json:"author_pk"`
	AuthorUsernameCipher []byte    `json:"author_username_cipher"`
	Type                 BlockType `json:"type"`

	DataId             string `json:"data_id"`
	DataKeyCipher      []byte `json:"data_key_cipher"`
	DataCaptionCipher  []byte `json:"data_caption_cipher"`
	DataMetadataCipher []byte `json:"data_metadata_cipher"`
}

type DataBlockConfig struct {
	DataId             string `json:"data_id"`
	DataKeyCipher      []byte `json:"data_key_cipher"`
	DataCaptionCipher  []byte `json:"data_caption_cipher"`
	DataMetadataCipher []byte `json:"data_metadata_cipher"`
}

type BlockType int

const (
	InviteBlock         BlockType = iota // no longer used
	ExternalInviteBlock                  // no longer used
	JoinBlock
	LeaveBlock
	PhotoBlock
	CommentBlock
	LikeBlock

	IgnoreBlock = 200
	MergeBlock  = 201
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
	ReceivedInviteNotification NotificationType = iota // peerA invited you
	DeviceAddedNotification                            // new device added
	PhotoAddedNotification                             // peerA added a photo
	CommentAddedNotification                           // peerA commented on peerB's photo, video, comment, etc.
	LikeAddedNotification                              // peerA liked peerB's photo, video, comment, etc.
	PeerJoinedNotification                             // peerA joined
	PeerLeftNotification                               // peerA left
	TextAddedNotification                              // peerA added a message
)

func (n NotificationType) Description() string {
	switch n {
	case ReceivedInviteNotification:
		return "RECEIVED_INVITE"
	case DeviceAddedNotification:
		return "DEVICE_ADDED"
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

type PinRequest struct {
	Id   string    `json:"id"`
	Date time.Time `json:"date"`
}

type CafeTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
