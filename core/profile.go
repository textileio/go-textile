package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/namesys/opts"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"time"
)

// Profile is an account-wide public profile
// NOTE: any account peer can publish profile entries to the same IPNS key
type Profile struct {
	Address   string   `json:"address"`
	Peers     []string `json:"peers"`
	Username  string   `json:"username,omitempty"`
	AvatarUri string   `json:"avatar_uri,omitempty"`
}

// profileTTL is the duration the ipns profile record will be valid
var profileTTL = time.Hour * 24 * 7 * 4

// profileCacheTTL is the duration the ipns profile record will be cached
var profileCacheTTL = time.Hour * 24 * 7

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
		pid, err := t.Id()
		if err != nil {
			log.Errorf("error getting id (set username): %s", err)
			return
		}
		prof, err := t.GetProfile(pid.Pretty())
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
		pid, err := t.Id()
		if err != nil {
			log.Errorf("error getting id (set avatar): %s", err)
			return
		}
		prof, err := t.GetProfile(pid.Pretty())
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
func (t *Textile) GetProfile(peerId string) (*Profile, error) {
	profile := &Profile{}

	// if peer id is local, return profile from db
	pid, err := t.Id()
	if err != nil {
		return nil, err
	}
	if pid.Pretty() == peerId {
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
	entry, err := t.ResolveProfile(peerId)
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
	if t.ipfs.Mounts.Ipns != nil && t.ipfs.Mounts.Ipns.IsActive() {
		return nil, errors.New("cannot manually publish while IPNS is mounted")
	}

	// if nil profile, use current
	if prof == nil {
		pid, err := t.Id()
		if err != nil {
			return nil, err
		}
		prof, err = t.GetProfile(pid.Pretty())
		if err != nil {
			return nil, err
		}
	}

	// create a virtual directory for the profile
	dir := uio.NewDirectory(t.ipfs.DAG)

	// add public components
	if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.Address)), "address"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.Username)), "username"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.AvatarUri)), "avatar_uri"); err != nil {
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
	id := node.Cid().Hash().B58String()

	// request cafe store
	t.cafeRequestQueue.Add(id, repo.CafeStoreRequest)

	// load our private key
	accnt, err := t.Account()
	if err != nil {
		return nil, err
	}
	sk, err := accnt.LibP2PPrivKey()
	if err != nil {
		return nil, err
	}

	// finish
	return t.publish(id, sk)
}

// ResolveProfile looks up a profile on ipns
func (t *Textile) ResolveProfile(name string) (*path.Path, error) {
	if !t.Online() {
		return nil, ErrOffline
	}

	// query options
	name = fmt.Sprintf("/ipns/%s", name)
	var ropts []nsopts.ResolveOpt
	ropts = append(ropts, nsopts.Depth(1))
	ropts = append(ropts, nsopts.DhtRecordCount(4))
	ropts = append(ropts, nsopts.DhtTimeout(5))

	// resolve w/ ipns
	pth, err := t.ipfs.Namesys.Resolve(t.ipfs.Context(), name, ropts...)
	if err != nil {
		return nil, err
	}
	return &pth, nil
}

func (t *Textile) publish(cid string, sk libp2pc.PrivKey) (*ipfs.IpnsEntry, error) {
	pth, err := path.ParsePath(cid)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(t.ipfs.Context(), "ipns-publish-ttl", profileTTL)
	entry, err := ipfs.Publish(ctx, t.ipfs, sk, pth, profileCacheTTL)
	if err != nil {
		return nil, err
	}

	log.Debugf("published: %s -> %s", entry.Name, entry.Value)

	return entry, nil
}
