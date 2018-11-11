package mill

import (
	"errors"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
)

var log = logging.Logger("tex-mill")

var ErrMediaTypeNotSupported = errors.New("media type not supported")

type Result struct {
	File []byte
	Meta map[string]interface{}
}

type Mill interface {
	ID() string
	Encrypt() bool
	Pin() bool
	AcceptMedia(media string) error
	Mill(input []byte, name string) (*Result, error)
}

func accepts(list []string, media string) error {
	for _, m := range list {
		if media == m {
			return nil
		}
	}
	return ErrMediaTypeNotSupported
}
