package mill

import (
	"errors"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"mime/multipart"
)

var log = logging.Logger("tex-mill")

var ErrMediaTypeNotSupported = errors.New("media type not supported")

type Result struct {
	File []byte
	Meta map[string]interface{}
}

type Mill interface {
	ID() string
	AcceptMedia(media string) error
	Mill(file multipart.File, name string) (*Result, error)
}

func accepts(list []string, media string) error {
	for _, m := range list {
		if media == m {
			return nil
		}
	}
	return ErrMediaTypeNotSupported
}
