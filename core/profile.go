package core

import (
	"crypto/rand"
	"io/ioutil"

	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"

	"github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/schema/textile"
)

// Profile returns this node's own contact info
func (t *Textile) Profile() *pb.Contact {
	return t.datastore.Contacts().Get(t.node.Identity.Pretty())
}

// Username returns profile username
func (t *Textile) Username() string {
	self := t.Profile()
	if self == nil {
		return ""
	}
	return self.Username
}

// SetUsername updates profile with a new username
func (t *Textile) SetUsername(username string) error {
	if username == t.Username() {
		return nil
	}
	if err := t.datastore.Contacts().UpdateUsername(t.node.Identity.Pretty(), username); err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(); err != nil {
			return err
		}
	}

	return t.PublishContact()
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
func (t *Textile) SetAvatar(hash string) error {
	if hash == t.Avatar() {
		return nil
	}

	data, file, err := t.FileData(hash)
	if err != nil {
		return err
	}
	input, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	// create a plaintext files thread for tracking avatars
	thrd := t.ThreadByKey("avatars")
	if thrd == nil {
		sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return err
		}

		sf, err := t.AddSchema(textile.Avatars, "avatars")
		if err != nil {
			return err
		}

		thrd, err = t.AddThread(pb.AddThreadConfig{
			Key:  "avatars",
			Name: "avatars",
			Schema: &pb.AddThreadConfig_Schema{
				Id: sf.Hash,
			},
			Type:    pb.Thread_Private,
			Sharing: pb.Thread_NotShared,
		}, sk, t.account.Address(), true)
		if err != nil {
			return err
		}
	}

	large, err := t.AddFileIndex(&mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   thrd.Schema.Links["large"].Opts["width"],
			Quality: thrd.Schema.Links["large"].Opts["quality"],
		},
	}, AddFileConfig{
		Input:     input,
		Media:     file.Media,
		Plaintext: thrd.Schema.Links["large"].Plaintext,
	})
	if err != nil {
		return err
	}

	small, err := t.AddFileIndex(&mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   thrd.Schema.Links["small"].Opts["width"],
			Quality: thrd.Schema.Links["small"].Opts["quality"],
		},
	}, AddFileConfig{
		Input:     input,
		Media:     file.Media,
		Plaintext: thrd.Schema.Links["small"].Plaintext,
	})
	if err != nil {
		return err
	}

	dir := map[string]*pb.FileIndex{"large": large, "small": small}
	dirs := &pb.DirectoryList{Items: []*pb.Directory{{Files: dir}}}
	node, keys, err := t.AddNodeFromDirs(dirs)
	if err != nil {
		return err
	}

	if _, err := thrd.AddFiles(node, "", keys.Files); err != nil {
		return err
	}

	avatar := node.Cid().Hash().B58String()
	if err := t.datastore.Contacts().UpdateAvatar(t.node.Identity.Pretty(), avatar); err != nil {
		return err
	}

	return t.PublishContact()
}
