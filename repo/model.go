package repo

import (
	"strconv"
	"time"
)

type Thread struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	PrivKey []byte `json:"priv_key"`
	Head    string `json:"head"`
}

type Block struct {
	Id           string    `json:"id"`
	Target       string    `json:"target"`
	Parents      []string  `json:"parents"`
	TargetKey    []byte    `json:"target_key"`
	ThreadPubKey string    `json:"thread_pub_key"`
	Type         BlockType `json:"type"`
	Date         time.Time `json:"date"`
}

type BlockType int

const (
	InviteBlock BlockType = iota
	PhotoBlock
	CommentBlock
	LikeBlock
)

func (bt BlockType) Bytes() []byte {
	return []byte(strconv.Itoa(int(bt)))
}
