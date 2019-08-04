package mobile

import (
	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/pb"
)

// SetLogLevel calls core SetLogLevel
func (m *Mobile) SetLogLevel(level []byte) error {
	mlevel := new(pb.LogLevel)
	if err := proto.Unmarshal(level, mlevel); err != nil {
		return err
	}

	return m.node.SetLogLevel(mlevel, false)
}
