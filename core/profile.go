package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/namesys/opts"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/path"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"time"
)

type Profile struct {
	Id        string                 `json:"id"`
	Username  string                 `json:"username,omitempty"`
	AvatarUri string                 `json:"avatar_uri,omitempty"`
	Threads   map[string]ThreadState `json:"threads,omitempty"`
}

type ThreadState struct {
	Head string `json:"head"`
	Sk   string `json:"sk"`
}

var profileTTL = time.Hour * 24 * 7 * 4
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
		<-t.Online()

		// publish
		pid, err := t.ID()
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
		<-t.Online()

		// publish
		pid, err := t.ID()
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
	profile := &Profile{Id: peerId}

	// if peer id is local, return profile from db
	pid, err := t.ID()
	if err != nil {
		return nil, err
	}
	if pid.Pretty() == peerId {
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
	var usernameb, avatarb []byte
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

// ResolveProfile looks up a profile on ipns
func (t *Textile) ResolveProfile(name string) (*path.Path, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}

	// setup query
	name = fmt.Sprintf("/ipns/%s", name)
	var ropts []nsopts.ResolveOpt
	ropts = append(ropts, nsopts.Depth(1))
	ropts = append(ropts, nsopts.DhtRecordCount(4))
	ropts = append(ropts, nsopts.DhtTimeout(5))

	pth, err := t.ipfs.Namesys.Resolve(t.ipfs.Context(), name, ropts...)
	if err != nil {
		return nil, err
	}
	return &pth, nil
}

// PublishProfile publishes profile to ipns
func (t *Textile) PublishProfile(prof *Profile) (*ipfs.IpnsEntry, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}
	if t.ipfs.Mounts.Ipns != nil && t.ipfs.Mounts.Ipns.IsActive() {
		return nil, errors.New("cannot manually publish while IPNS is mounted")
	}

	// if nil profile, use current
	if prof == nil {
		pid, err := t.ID()
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
	if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.Id)), "id"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.Username)), "username"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader([]byte(prof.AvatarUri)), "avatar_uri"); err != nil {
		return nil, err
	}

	// add private encrypted threads state
	threadsDir := uio.NewDirectory(t.ipfs.DAG)
	for _, thrd := range t.threads {
		dir := uio.NewDirectory(t.ipfs.DAG)
		head, err := thrd.GetHead()
		if err != nil {
			return nil, err
		}
		if head != "" {
			headc, err := t.Encrypt([]byte(head))
			if err != nil {
				return nil, err
			}
			if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader(headc), "head"); err != nil {
				return nil, err
			}
		}
		skb, err := thrd.PrivKey.Bytes()
		if err != nil {
			return nil, err
		}
		skc, err := t.Encrypt(skb)
		if err != nil {
			return nil, err
		}
		if err := ipfs.AddFileToDirectory(t.ipfs, dir, bytes.NewReader(skc), "sk"); err != nil {
			return nil, err
		}
		id, err := thrd.Base58Id()
		if err != nil {
			return nil, err
		}
		node, err := dir.GetNode()
		if err != nil {
			return nil, err
		}
		if err := ipfs.PinDirectory(t.ipfs, node, []string{}); err != nil {
			return nil, err
		}
		if err := threadsDir.AddChild(t.ipfs.Context(), id, node); err != nil {
			return nil, err
		}
	}
	threadsNode, err := threadsDir.GetNode()
	if err != nil {
		return nil, err
	}
	if err := ipfs.PinDirectory(t.ipfs, threadsNode, []string{}); err != nil {
		return nil, err
	}
	if err := dir.AddChild(t.ipfs.Context(), "threads", threadsNode); err != nil {
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

	// request cafe pin
	go func() {
		if err := t.putPinRequest(id); err != nil {
			// TODO: #202 (Properly handle database/sql errors)
			log.Warningf("pin request exists: %s", id)
		}
	}()

	// extract path
	pth, err := path.ParsePath(id)
	if err != nil {
		return nil, err
	}

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
	popts := &ipfs.PublishOpts{
		VerifyExists: true,
		PubValidTime: profileCacheTTL,
	}
	ctx := context.WithValue(t.ipfs.Context(), "ipns-publish-ttl", profileTTL)
	entry, err := ipfs.Publish(ctx, t.ipfs, sk, pth, popts)
	if err != nil {
		return nil, err
	}

	log.Debugf("updated profile: %s -> %s", entry.Name, entry.Value)

	return entry, nil
}
