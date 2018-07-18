package model

import (
	"github.com/textileio/textile-go/net"
)

type AddResult struct {
	Id         string          `json:"id"`
	Key        string          `json:"key"`
	PinRequest *net.PinRequest `json:"pin_request"`
}
