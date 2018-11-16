package mill

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"

	"github.com/mr-tron/base58/base58"
)

var log = logging.Logger("tex-mill")

var ErrMediaTypeNotSupported = errors.New("media type not supported")

type Result struct {
	File []byte
	Meta map[string]interface{}
	Opts string // hash of opts, could be empty
}

type Mill interface {
	ID() string
	Encrypt() bool
	Pin() bool
	AcceptMedia(media string) error
	Options() (string, error)
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

func hashOpts(opts interface{}) (string, error) {
	data, err := json.Marshal(opts)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return base58.FastBase58Encoding(sum[:]), nil
}
