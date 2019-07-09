package mill

import (
	"bytes"
	"fmt"

	"github.com/textileio/go-textile/keypair"

	"github.com/golang/protobuf/jsonpb"
	"github.com/textileio/go-textile/pb"
)

type Roles struct{}

func (m *Roles) ID() string {
	return "/roles"
}

func (m *Roles) Encrypt() bool {
	return true
}

func (m *Roles) Pin() bool {
	return true
}

func (m *Roles) AcceptMedia(media string) error {
	return accepts([]string{"application/json"}, media)
}

func (m *Roles) Options(add map[string]interface{}) (string, error) {
	return hashOpts(make(map[string]string), add)
}

func (m *Roles) Mill(input []byte, name string) (*Result, error) {
	var roles pb.Thread2_Roles
	err := jsonpb.Unmarshal(bytes.NewReader(input), &roles)
	if err != nil {
		return nil, err
	}

	for account := range roles.Accounts {
		kp, err := keypair.Parse(account)
		if err != nil {
			return nil, fmt.Errorf("error parsing address: %s", err)
		}
		_, err = kp.Sign([]byte{0x00})
		if err == nil {
			// we don't want to handle account seeds, just addresses
			return nil, fmt.Errorf("entry is an account seed, not address")
		}
	}

	data, err := pbMarshaler.MarshalToString(&roles)
	if err != nil {
		return nil, err
	}

	return &Result{File: []byte(data)}, nil
}
