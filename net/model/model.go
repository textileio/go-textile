package model

import "github.com/textileio/textile-go/net"

type AddResult struct {
	Id            string
	Key           []byte
	RemoteRequest *net.MultipartRequest
}
