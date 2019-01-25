package core

import (
	"time"

	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// CreateCafeToken creates a single random developer access token to be used to register with a Cafe
func (t *Textile) CreateCafeToken() (*repo.CafeDevToken, error) {
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, err
	}

	token := &repo.CafeDevToken{
		Id:      ksuid.New().String(),
		Token:   key[:32],
		Created: time.Now(),
	}
	error = t.datastore.CafeDevTokens().Add(token)
	if err != nil {
		return nil, err
	}

	return token, nil

}

// CafeDevTokens lists all stored (salted and encrypted) dev tokens
func (t *Textile) CafeDevTokens() ([]*repo.CafeDevToken, error) {
	return t.datastore.CafeDevTokens().List(), nil
}

// CheckCafeDevToken checks whether a given dev token is valid
func (t *Textile) CheckCafeDevToken(id string) (bool, error) {
	token, err := t.datastore.CafeDevTokens().Get(id)
	if err != nil {
		return nil, err
	}

}

// protoCafeToRepo is a tmp method just converting proto cafe info to the repo version
func protoCafeToRepo(pro *pb.Cafe) repo.Cafe {
	return repo.Cafe{
		Peer:     pro.Peer,
		Address:  pro.Address,
		API:      pro.Api,
		Protocol: pro.Protocol,
		Node:     pro.Node,
		URL:      pro.Url,
		Swarm:    pro.Swarm,
	}
}

// repoCafeToProto is a tmp method just converting repo cafe info to the proto version
func repoCafeToProto(rep repo.Cafe) *pb.Cafe {
	return &pb.Cafe{
		Peer:     rep.Peer,
		Address:  rep.Address,
		Api:      rep.API,
		Protocol: rep.Protocol,
		Node:     rep.Node,
		Url:      rep.URL,
		Swarm:    rep.Swarm,
	}
}
