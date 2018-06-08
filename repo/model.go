package repo

import (
	"strconv"
	"time"
)

type Thread struct {
	Id      string
	Name    string
	PrivKey []byte
	Head    string
}

type Block struct {
	Id           string
	Target       string
	Parents      []string
	TargetKey    []byte
	ThreadPubKey []byte
	Type         BlockType
	Date         time.Time
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
