package core

import (
	"crypto/rand"
	"io/ioutil"

	libp2pc "github.com/libp2p/go-libp2p-crypto"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema/textile"
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
			Type:    pb.Thread_PRIVATE,
			Sharing: pb.Thread_NOT_SHARED,
		}, sk, t.account.Address(), true, false)
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
	if err := t.datastore.Peers().UpdateAvatar(t.node.Identity.Pretty(), avatar); err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(nil); err != nil {
			return err
		}
	}

	return t.publishPeer()
}
