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
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	Parents  []string  `json:"parents"`
	ThreadId string    `json:"thread_id"`
	AuthorPk string    `json:"author_pk"`
	Type     BlockType `json:"type"`

	DataId            string `json:"data_id"`
	DataKeyCipher     []byte `json:"data_key_cipher"`
	DataCaptionCipher []byte `json:"data_caption_cipher"`
}

type DataBlockConfig struct {
	DataId            string `json:"data_id"`
	DataKeyCipher     []byte `json:"data_key_cipher"`
	DataCaptionCipher []byte `json:"data_caption_cipher"`
}

type BlockType int

const (
	InviteBlock BlockType = iota
	ExternalInviteBlock
	JoinBlock
	LeaveBlock
	PhotoBlock

	IgnoreBlock = 200
	MergeBlock  = 201
)

type PinRequest struct {
	Id   string    `json:"id"`
	Date time.Time `json:"date"`
}

type CafeTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
