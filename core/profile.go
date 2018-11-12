package core

import (
	"bytes"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"strings"
	"time"
)

// Profile is an account-wide public profile
// NOTE: any account peer can publish profile entries to the same IPNS key
type Profile struct {
	Address   string   `json:"address"`
	Inboxes   []string `json:"inboxes"`
	Username  string   `json:"username,omitempty"`
	AvatarUri string   `json:"avatar_uri,omitempty"`
}

// profileLifetime is the duration the ipns profile record will be considered valid
var profileLifetime = time.Hour * 24 * 7

// profileTTL is the duration the ipns profile record will be locally cached
var profileTTL = time.Hour

// Username returns profile username
func (t *Textile) Username() (*string, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.Profile().GetUsername()
}

// SetUsername updates profile with a new username
func (t *Textile) SetUsername(username string) error {
	if err := t.touchDatastore(); err != nil {
		return err
	}
	if err := t.datastore.Profile().SetUsername(username); err != nil {
		return err
	}

	// annouce to all threads
	for _, thrd := range t.threads {
		if _, err := thrd.annouce(); err != nil {
			return err
		}
	}

	return t.PublishProfile()
}

// Avatar returns profile avatar
func (t *Textile) Avatar() (*string, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.Profile().GetAvatar()
}

// SetAvatar updates profile with a new avatar at the given photo id
func (t *Textile) SetAvatar(id string) error {
	if err := t.touchDatastore(); err != nil {
		return err
	}

	// get the public key for this photo
	// TODO: fix
	//key, err := t.PhotoKey(id)
	//if err != nil {
	//	return err
	//}

	// build a public uri
	uri := fmt.Sprintf("/ipfs/%s/thumb?key=%s", id, "fixme")
	if err := t.datastore.Profile().SetAvatar(uri); err != nil {
		return err
	}
	return t.PublishProfile()
}

// Profile return a model representation of an ipns profile
func (t *Textile) Profile(pid peer.ID) (*Profile, error) {
	if !t.Started() {
		return nil, ErrStopped
	}
	profile := &Profile{}

	// if peer id is local, return profile from db
	if t.node.Identity.Pretty() == pid.Pretty() {
		profile.Address = t.account.Address()
		for _, ses := range t.datastore.CafeSessions().List() {
			profile.Inboxes = append(profile.Inboxes, ses.CafeId)
		}
		username, err := t.Username()
		if err != nil {
			return nil, err
		}
		if username != nil {
			profile.Username = *username
		}
		avatar, err := t.Avatar()
		if err != nil {
			return nil, err
		}
		if avatar != nil {
			profile.AvatarUri = *avatar
		}
		return profile, nil
	}

	entry, err := t.ResolveProfile(pid)
	if err != nil {
		return nil, err
	}
	root := entry.String()

	var addrb, inboxesb, usernameb, avatarb []byte
	addrb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "address"))
	if addrb != nil {
		profile.Address = string(addrb)
	}
	inboxesb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "inboxes"))
	if inboxesb != nil && string(inboxesb) != "" {
		profile.Inboxes = strings.Split(string(inboxesb), ",")
	}
	usernameb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "username"))
	if usernameb != nil {
		profile.Username = string(usernameb)
	}
	avatarb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "avatar_uri"))
	if avatarb != nil {
		profile.AvatarUri = string(avatarb)
	}
	return profile, nil
}

// PublishProfile publishes the current profile
func (t *Textile) PublishProfile() error {
	prof, err := t.Profile(t.node.Identity)
	if err != nil {
		return err
	}
	if prof == nil {
		return nil
	}
	go func() {
		<-t.OnlineCh()
		entry, err := t.publishProfile(*prof)
		if err != nil {
			log.Errorf("error publishing profile: %s", err)
			return
		}
		log.Debugf("published: %s -> %s", entry.Name, entry.Value)
	}()
	return nil
}

// ResolveProfile looks up a profile on ipns
func (t *Textile) ResolveProfile(name peer.ID) (*path.Path, error) {
	if !t.Online() {
		return nil, ErrOffline
	}
	return ipfs.Resolve(t.node, name)
}

// publishProfile publishes profile to ipns
func (t *Textile) publishProfile(prof Profile) (*ipfs.IpnsEntry, error) {
	if !t.Online() {
		return nil, ErrOffline
	}

	dir := uio.NewDirectory(t.node.DAG)

	// add public components
	addressId, err := ipfs.AddDataToDirectory(t.node, dir, "address", bytes.NewReader([]byte(prof.Address)))
	if err != nil {
		return nil, err
	}

	var inboxesId *cid.Cid
	sessions := t.datastore.CafeSessions().List()
	if len(sessions) > 0 {
		var inboxes []string
		for _, ses := range t.datastore.CafeSessions().List() {
			inboxes = append(inboxes, ses.CafeId)
		}
		inboxesStr := strings.Join(inboxes, ",")
		inboxesId, err = ipfs.AddDataToDirectory(t.node, dir, "inboxes", bytes.NewReader([]byte(inboxesStr)))
		if err != nil {
			return nil, err
		}
	} else {
		inboxesId, err = ipfs.AddDataToDirectory(t.node, dir, "inboxes", bytes.NewReader([]byte("")))
		if err != nil {
			return nil, err
		}
	}
	usernameId, err := ipfs.AddDataToDirectory(t.node, dir, "username", bytes.NewReader([]byte(prof.Username)))
	if err != nil {
		return nil, err
	}
	avatarId, err := ipfs.AddDataToDirectory(t.node, dir, "avatar_uri", bytes.NewReader([]byte(prof.AvatarUri)))
	if err != nil {
		return nil, err
	}

	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return nil, err
	}

	t.cafeOutbox.Add(addressId.Hash().B58String(), repo.CafeStoreRequest)
	t.cafeOutbox.Add(inboxesId.Hash().B58String(), repo.CafeStoreRequest)
	t.cafeOutbox.Add(usernameId.Hash().B58String(), repo.CafeStoreRequest)
	t.cafeOutbox.Add(avatarId.Hash().B58String(), repo.CafeStoreRequest)
	t.cafeOutbox.Add(node.Cid().Hash().B58String(), repo.CafeStoreRequest)
	go t.cafeOutbox.Flush()

	value := node.Cid().Hash().B58String()
	return ipfs.Publish(t.node, t.node.PrivateKey, value, profileLifetime, profileTTL)
}
