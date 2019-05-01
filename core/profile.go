package core

import (
	"github.com/textileio/go-textile/pb"
)

// Profile returns this node's own peer
func (t *Textile) Profile() *pb.Peer {
	return t.datastore.Peers().Get(t.node.Identity.Pretty())
}

// Username returns profile username
func (t *Textile) Name() string {
	self := t.Profile()
	if self == nil {
		return ""
	}
	return self.Name
}

// SetName updates profile with a new username
func (t *Textile) SetName(name string) error {
	if name == t.Name() {
		return nil
	}
	if err := t.datastore.Peers().UpdateName(t.node.Identity.Pretty(), name); err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(nil); err != nil {
			return err
		}
	}

	return t.publishPeer()
}

// Avatar returns profile avatar
func (t *Textile) Avatar() string {
	self := t.Profile()
	if self == nil {
		return ""
	}
	return self.Avatar
}

// SetAvatar updates profile with a new avatar at the given file hash.
func (t *Textile) SetAvatar() error {

	t.AccountThread().Head()

	if hash == t.Avatar() {
		return nil
	}

	//avatar := node.Cid().Hash().B58String()
	//if err := t.datastore.Peers().UpdateAvatar(t.node.Identity.Pretty(), avatar); err != nil {
	//	return err
	//}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(nil); err != nil {
			return err
		}
	}

	return t.publishPeer()
}
