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
	ThreadId string `json:"thread_id"`
	PubKey   []byte `json:"pk"`
}

type Block struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	Parents  []string  `json:"parents"`
	ThreadId string    `json:"thread_pk"`
	AuthorPk string    `json:"author_pk"`

	Type      BlockType `json:"type"`
	Target    string    `json:"target"`
	TargetKey []byte    `json:"target_key"`
}

type BlockType int

const (
	InviteBlock BlockType = iota
	ExternalInviteBlock
	JoinBlock
	LeaveBlock
	PhotoBlock
)
