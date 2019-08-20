package core

import (
	"fmt"

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
	err := t.datastore.Peers().UpdateName(t.node.Identity.Pretty(), name)
	if err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		_, err = thrd.Annouce(nil)
		if err != nil {
			return err
		}
	}

	return t.PublishPeer()
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
	latest := t.AccountThread().LatestFiles()
	if latest == nil {
		return fmt.Errorf("account thread contains no files")
	}

	if latest.Data == t.Avatar() {
		return nil
	}

	err := t.datastore.Peers().UpdateAvatar(t.node.Identity.Pretty(), latest.Data)
	if err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		_, err = thrd.Annouce(nil)
		if err != nil {
			return err
		}
	}

	return t.PublishPeer()
}
