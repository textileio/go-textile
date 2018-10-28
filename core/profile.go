package core

import (
	"bytes"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"time"
)

// Profile is an account-wide public profile
// NOTE: any account peer can publish profile entries to the same IPNS key
type Profile struct {
	Address   string `json:"address"`
	Username  string `json:"username,omitempty"`
	AvatarUri string `json:"avatar_uri,omitempty"`
}

// profileLifetime is the duration the ipns profile record will be considered valid
var profileLifetime = time.Hour * 24 * 7

// profileTTL is the duration the ipns profile record will be locally cached
var profileTTL = time.Hour

// GetUsername returns profile username
func (t *Textile) GetUsername() (*string, error) {
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

	// update
	if err := t.datastore.Profile().SetUsername(username); err != nil {
		return err
	}

	go func() {
		<-t.OnlineCh()

		// publish
		prof, err := t.GetProfile(t.ipfs.Identity)
		if err != nil {
			log.Errorf("error getting profile (set username): %s", err)
			return
		}
		if _, err := t.PublishProfile(prof); err != nil {
			log.Errorf("error publishing profile (set username): %s", err)
			return
		}
	}()
	return nil
}

// GetAvatar returns profile avatar
func (t *Textile) GetAvatar() (*string, error) {
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
	key, err := t.GetPhotoKey(id)
	if err != nil {
		return err
	}

	// build a public uri
	uri := fmt.Sprintf("/ipfs/%s/thumb?key=%s", id, key)

	// update
	if err := t.datastore.Profile().SetAvatar(uri); err != nil {
		return err
	}

	go func() {
		<-t.OnlineCh()

		// publish
		prof, err := t.GetProfile(t.ipfs.Identity)
		if err != nil {
			log.Errorf("error getting profile (set avatar): %s", err)
			return
		}
		if _, err := t.PublishProfile(prof); err != nil {
			log.Errorf("error publishing profile (set avatar): %s", err)
			return
		}
	}()
	return nil
}

// GetProfile return a model representation of an ipns profile
func (t *Textile) GetProfile(pid peer.ID) (*Profile, error) {
	profile := &Profile{}

	// if peer id is local, return profile from db
	if t.ipfs.Identity.Pretty() == pid.Pretty() {
		addr, err := t.Address()
		if err != nil {
			return nil, err
		}
		profile.Address = addr
		username, err := t.GetUsername()
		if err != nil {
			return nil, err
		}
		if username != nil {
			profile.Username = *username
		}
		avatar, err := t.GetAvatar()
		if err != nil {
			return nil, err
		}
		if avatar != nil {
			profile.AvatarUri = *avatar
		}
		return profile, nil
	}

	// resolve profile at peer id
	entry, err := t.ResolveProfile(pid)
	if err != nil {
		return nil, err
	}
	root := entry.String()

	// get components from entry
	var addrb, usernameb, avatarb []byte
	addrb, _ = ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s/%s", root, "address"))
	if addrb != nil {
		profile.Address = string(addrb)
	}
	usernameb, _ = ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s/%s", root, "username"))
	if usernameb != nil {
		profile.Username = string(usernameb)
	}
	avatarb, _ = ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s/%s", root, "avatar_uri"))
	if avatarb != nil {
		profile.AvatarUri = string(avatarb)
	}
	return profile, nil
}

// PublishProfile publishes profile to ipns
func (t *Textile) PublishProfile(prof *Profile) (*ipfs.IpnsEntry, error) {
	if !t.Online() {
		return nil, ErrOffline
	}

	// if nil profile, use current
	if prof == nil {
		var err error
		prof, err = t.GetProfile(t.ipfs.Identity)
		if err != nil {
			return nil, err
		}
	}

	// create a virtual directory for the profile
	dir := uio.NewDirectory(t.ipfs.DAG)

	// add public components
	addressId, err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.Address)), "address")
	if err != nil {
		return nil, err
	}
	usernameId, err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.Username)), "username")
	if err != nil {
		return nil, err
	}
	avatarId, err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.AvatarUri)), "avatar_uri")
	if err != nil {
		return nil, err
	}

	// pin the directory locally
	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}
	if err := ipfs.PinDirectory(t.ipfs, node, []string{}); err != nil {
		return nil, err
	}

	// add store requests
	t.cafeOutbox.Add(addressId.Hash().B58String(), repo.CafeStoreRequest)
	t.cafeOutbox.Add(usernameId.Hash().B58String(), repo.CafeStoreRequest)
	t.cafeOutbox.Add(avatarId.Hash().B58String(), repo.CafeStoreRequest)
	t.cafeOutbox.Add(node.Cid().Hash().B58String(), repo.CafeStoreRequest)
	go t.cafeOutbox.Flush()

	// finish
	value := node.Cid().Hash().B58String()
	entry, err := ipfs.Publish(t.ipfs, t.ipfs.PrivateKey, value, profileLifetime, profileTTL)
	if err != nil {
		return nil, err
	}
	log.Debugf("published: %s -> %s", entry.Name, entry.Value)
	return entry, nil
}

// ResolveProfile looks up a profile on ipns
func (t *Textile) ResolveProfile(name peer.ID) (*path.Path, error) {
	if !t.Online() {
		return nil, ErrOffline
	}
	return ipfs.Resolve(t.ipfs, name)
}
