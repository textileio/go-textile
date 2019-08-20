package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/schema/textile"
	"github.com/textileio/go-textile/util"
)

// ErrThreadNotFound indicates thread is not found in the loaded list
var ErrThreadNotFound = fmt.Errorf("thread not found")

// ErrThreadLoaded indicates the thread is already loaded from the datastore
var ErrThreadLoaded = fmt.Errorf("thread is loaded")

// emptyThreadKey indicates "" was used for a thread key
var emptyThreadKey = fmt.Errorf("thread key cannot by empty")

// AddThread adds a thread with a given name and secret key
func (t *Textile) AddThread(conf pb.AddThreadConfig, sk libp2pc.PrivKey, initiator string, join bool, inviteAccount bool) (*Thread, error) {
	conf.Key = strings.TrimSpace(conf.Key)
	if conf.Key == "" {
		return nil, emptyThreadKey
	}

	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	skb, err := sk.Bytes()
	if err != nil {
		return nil, err
	}

	var sch string
	if conf.Schema != nil {
		var sjson string

		if conf.Schema.Id != "" {
			// ensure schema id is a multi hash
			_, err = mh.FromB58String(conf.Schema.Id)
			if err != nil {
				return nil, err
			}
			sch = conf.Schema.Id
		} else if conf.Schema.Json != "" {
			sjson = conf.Schema.Json
		} else {
			switch conf.Schema.Preset {
			case pb.AddThreadConfig_Schema_BLOB:
				sjson = textile.Blob
			case pb.AddThreadConfig_Schema_CAMERA_ROLL:
				sjson = textile.CameraRoll
			case pb.AddThreadConfig_Schema_MEDIA:
				sjson = textile.Media
			}
		}

		if sjson != "" {
			sfile, err := t.AddFileIndex(&mill.Schema{}, AddFileConfig{
				Input: []byte(sjson),
				Media: "application/json",
			})
			if err != nil {
				return nil, err
			}
			sch = sfile.Hash
		}

		if sch != "" {
			err = t.cafeOutbox.Add(sch, pb.CafeRequest_STORE)
			if err != nil {
				return nil, err
			}
		}
	}

	// ensure whitelist is unique
	set := make(map[string]struct{})
	var members []string
	for _, m := range conf.Whitelist {
		if _, ok := set[m]; !ok {
			kp, err := keypair.Parse(m)
			if err != nil {
				return nil, fmt.Errorf("error parsing address: %s", err)
			}
			_, err = kp.Sign([]byte{0x00})
			if err == nil {
				// we don't want to handle account seeds, just addresses
				return nil, fmt.Errorf("entry is an account seed, not address")
			}
			members = append(members, m)
		}
		set[m] = struct{}{}
	}

	model := &pb.Thread{
		Id:        id.Pretty(),
		Key:       conf.Key,
		Sk:        skb,
		Name:      strings.TrimSpace(conf.Name),
		Schema:    sch,
		Initiator: initiator,
		Type:      conf.Type,
		Sharing:   conf.Sharing,
		Whitelist: members,
		State:     pb.Thread_LOADED,
	}
	err = t.datastore.Threads().Add(model)
	if err != nil {
		if conf.Force && db.ConflictError(err) && strings.Contains(err.Error(), ".key") {
			conf.Key = incrementKey(conf.Key)
			return t.AddThread(conf, sk, initiator, join, inviteAccount)
		}
		return nil, err
	}

	thread, err := t.loadThread(model)
	if err != nil {
		return nil, err
	}

	// we join here if we're the creator
	if join {
		_, err = thread.join("")
		if err != nil {
			return nil, err
		}
	}

	t.sendUpdate(&pb.AccountUpdate{
		Id:   thread.Id,
		Type: pb.AccountUpdate_THREAD_ADDED,
	})

	// invite account peers if inviter is not an account peer
	if inviteAccount {
		for _, p := range t.accountPeers() {
			_, err = thread.AddInvite(p)
			if err != nil {
				return nil, err
			}
		}
	}

	log.Debugf("added a new thread %s with name %s", thread.Id, conf.Name)

	return thread, nil
}

// AddOrUpdateThread add or updates a thread directly, usually from a backup
func (t *Textile) AddOrUpdateThread(thread *pb.Thread) error {
	// check if we're allowed to get an invite
	// Note: just using a dummy thread here because having these access+sharing
	// methods on Thread is very nice elsewhere.
	dummy := &Thread{
		initiator: thread.Initiator,
		ttype:     thread.Type,
		sharing:   thread.Sharing,
		whitelist: thread.Whitelist,
	}
	if !dummy.shareable(t.config.Account.Address, t.config.Account.Address) {
		return ErrNotShareable
	}

	sk, err := ipfs.UnmarshalPrivateKey(thread.Sk)
	if err != nil {
		return err
	}

	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return err
	}

	heads := util.SplitString(thread.Head, ",")
	nthread := t.Thread(id.Pretty())
	if nthread == nil {
		config := pb.AddThreadConfig{
			Key:  thread.Key,
			Name: thread.Name,
			Schema: &pb.AddThreadConfig_Schema{
				Id: thread.Schema,
			},
			Type:      thread.Type,
			Sharing:   thread.Sharing,
			Whitelist: thread.Whitelist,
			Force:     true,
		}

		var err error
		nthread, err = t.AddThread(config, sk, thread.Initiator, false, false)
		if err != nil {
			return err
		}
		err = nthread.updateHead(heads, false)
		if err != nil {
			return err
		}
	}

	// have we joined?
	query := fmt.Sprintf("threadId='%s' and type=%d and authorId='%s'", nthread.Id, pb.Block_JOIN, t.node.Identity.Pretty())
	if t.datastore.Blocks().Count(query) == 0 {
		// go ahead, invite yourself
		_, err = nthread.join(t.node.Identity.Pretty())
		if err != nil {
			return err
		}
	}

	// compare heads to determine if we need to backtrack
	xheads, err := nthread.Heads()
	if err != nil {
		return err
	}
	if util.EqualStringSlices(xheads, heads) {
		t.FlushCafes()
		return nil
	}

	// handle the thread tail in the background
	stopGroup.Add(1, "AddOrUpdateThread")
	go func() {
		defer stopGroup.Done("AddOrUpdateThread")

		leaves := nthread.followParents(heads)
		err = nthread.handleHead(heads, leaves)
		if err != nil {
			log.Warningf("failed to handle head %s: %s", thread.Head, err)
			return
		}

		// handle newly discovered peers during back prop
		err = nthread.sendWelcome()
		if err != nil {
			log.Warningf("error sending welcome: %s", err)
			return
		}

		// flush cafe queue _at the very end_
		t.cafeOutbox.Flush(false)
	}()

	return nil
}

// RenameThread adds an announce block to the thread w/ a new name
// Note: Only thread initiators can update the thread's name
func (t *Textile) RenameThread(id string, name string) error {
	thread := t.Thread(id)
	if thread == nil {
		return ErrThreadNotFound
	}
	if thread.initiator != t.account.Address() {
		return fmt.Errorf("thread name is not writable")
	}

	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil
	}

	thread.Name = trimmed
	err := t.datastore.Threads().UpdateName(thread.Id, trimmed)
	if err != nil {
		return err
	}

	_, err = thread.Annouce(&pb.ThreadAnnounce{Name: trimmed})
	return err
}

// Thread get a thread by id from loaded threads
func (t *Textile) Thread(id string) *Thread {
	for _, thread := range t.loadedThreads {
		if thread.Id == id {
			return thread
		}
	}
	return nil
}

// Threads lists loaded threads
func (t *Textile) Threads() []*Thread {
	var threads []*Thread
	for _, i := range t.loadedThreads {
		if i == nil || i.Key == t.account.Address() {
			continue
		}
		threads = append(threads, i)
	}
	return threads
}

// ThreadPeers returns a list of thread peers
func (t *Textile) ThreadPeers(id string) (*pb.PeerList, error) {
	thread := t.Thread(id)
	if thread == nil {
		return nil, ErrThreadNotFound
	}

	peers := &pb.PeerList{Items: make([]*pb.Peer, 0)}
	for _, tp := range thread.Peers() {
		p := t.datastore.Peers().Get(tp.Id)
		if p != nil {
			peers.Items = append(peers.Items, p)
		}
	}

	return peers, nil
}

// RemoveThread removes a thread
// @todo rename to abandon to be consistent with CLI+API
func (t *Textile) RemoveThread(id string) (mh.Multihash, error) {
	var thread *Thread
	var index int
	for i, th := range t.loadedThreads {
		if th.Id == id {
			thread = th
			index = i
			break
		}
	}
	if thread == nil {
		return nil, ErrThreadNotFound
	}

	// notify peers
	addr, err := thread.leave()
	if err != nil {
		log.Errorf("error leaving thread %s: %s", id, err)
	}

	// delete backups
	err = t.cafeOutbox.Add(thread.Id, pb.CafeRequest_UNSTORE_THREAD)
	if err != nil {
		return nil, err
	}

	err = t.datastore.Threads().Delete(thread.Id)
	if err != nil {
		return nil, err
	}

	copy(t.loadedThreads[index:], t.loadedThreads[index+1:])
	t.loadedThreads[len(t.loadedThreads)-1] = nil
	t.loadedThreads = t.loadedThreads[:len(t.loadedThreads)-1]

	t.sendUpdate(&pb.AccountUpdate{
		Id:   thread.Id,
		Type: pb.AccountUpdate_THREAD_REMOVED,
	})

	log.Infof("removed thread %s with name %s", thread.Id, thread.Name)

	return addr, nil
}

// ThreadByKey get a thread by key from loaded threads
func (t *Textile) ThreadByKey(key string) *Thread {
	for _, thread := range t.loadedThreads {
		if thread.Key == key {
			return thread
		}
	}
	return nil
}

// ThreadView returns a thread with expanded view properties
func (t *Textile) ThreadView(id string) (*pb.Thread, error) {
	thread := t.Thread(id)
	if thread == nil {
		return nil, ErrThreadNotFound
	}

	mod := t.datastore.Threads().Get(thread.Id)
	if mod == nil {
		return nil, errThreadReload
	}

	// add extra view info
	mod.SchemaNode = thread.Schema
	for _, head := range util.SplitString(mod.Head, ",") {
		hid, err := blockCIDFromNode(t.node, head)
		if err == nil {
			block := t.datastore.Blocks().Get(hid)
			if block != nil {
				block.User = t.PeerUser(block.Author)
				mod.HeadBlocks = append(mod.HeadBlocks, block)
			}
		} else {
			log.Errorf("error getting node block %s: %s", head, err)
		}
	}
	mod.BlockCount = int32(t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s'", thread.Id)))
	mod.PeerCount = int32(len(thread.Peers()) + 1)

	return mod, nil
}

// SnapshotThreads creates a store thread request for all threads
func (t *Textile) SnapshotThreads() error {
	var err error
	for _, thread := range t.loadedThreads {
		err = thread.store()
		if err != nil {
			return err
		}
	}
	return nil
}

// SearchThreadSnapshots searches the network for snapshots
func (t *Textile) SearchThreadSnapshots(query *pb.ThreadSnapshotQuery, options *pb.QueryOptions) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, nil, nil, err
	}

	// settings required for sync
	options.RemoteOnly = true
	options.Limit = -1
	options.Filter = pb.QueryOptions_NO_FILTER

	resCh, errCh, cancel := t.search(&pb.Query{
		Type:    pb.Query_THREAD_SNAPSHOTS,
		Options: options,
		Payload: &any.Any{
			TypeUrl: "/ThreadSnapshotQuery",
			Value:   payload,
		},
	})

	// transform and filter results into plaintext
	backups := make(map[string]struct{})
	tresCh := make(chan *pb.QueryResult)
	terrCh := make(chan error)
	go func() {
		for {
			select {
			case res, ok := <-resCh:
				if !ok {
					close(tresCh)
					return
				}

				backup := new(pb.CafeClientThread)
				if err := ptypes.UnmarshalAny(res.Value, backup); err != nil {
					terrCh <- err
					break
				}

				plaintext, err := t.account.Decrypt(backup.Ciphertext)
				if err != nil {
					terrCh <- err
					break
				}

				thread := new(pb.Thread)
				err = proto.Unmarshal(plaintext, thread)
				if err != nil {
					terrCh <- err
					break
				}

				res.Id += ":" + thread.Head
				if _, ok := backups[res.Id]; ok {
					continue
				}
				backups[res.Id] = struct{}{}

				res.Value = &any.Any{
					TypeUrl: "/Thread",
					Value:   plaintext,
				}
				tresCh <- res

			case err := <-errCh:
				terrCh <- err
			}
		}
	}()

	return tresCh, terrCh, cancel, nil
}

// addAccountThread adds a thread with seed representing the state of the account
func (t *Textile) addAccountThread() error {
	x := t.AccountThread()
	if x != nil {
		aid, err := t.account.Id()
		if err != nil {
			return err
		}
		// catch malformed account threads from 0.1.10
		if x.Id == aid.Pretty() {

			// catch schema-less account threads from 0.1.11
			if x.Schema == nil {
				sf, err := t.AddSchema(textile.Avatars, "avatars")
				if err != nil {
					return err
				}
				return x.UpdateSchema(sf.Hash)
			}

			return nil
		}
		_, err = t.RemoveThread(x.Id)
		if err != nil {
			return err
		}
	}

	sf, err := t.AddSchema(textile.Avatars, "avatars")
	if err != nil {
		return err
	}

	config := pb.AddThreadConfig{
		Key:  t.account.Address(),
		Name: "account",
		Schema: &pb.AddThreadConfig_Schema{
			Id: sf.Hash,
		},
		Type:    pb.Thread_PRIVATE,
		Sharing: pb.Thread_NOT_SHARED,
	}
	sk, err := t.account.LibP2PPrivKey()
	if err != nil {
		return err
	}
	thread, err := t.AddThread(config, sk, t.account.Address(), true, false)
	if err != nil {
		return err
	}

	// add existing contacts
	for _, p := range t.datastore.Peers().List(fmt.Sprintf("address!='%s'", t.account.Address())) {
		_, err = thread.Annouce(&pb.ThreadAnnounce{Peer: p})
		if err != nil {
			return err
		}
	}

	return nil
}

// incrementKey add "_xxx" to the end of a key
func incrementKey(key string) string {
	_, err := strconv.Atoi(key)
	if err == nil {
		return key + "_1"
	}
	a := strings.Split(key, "_")
	var x string
	x, a = a[len(a)-1], a[:len(a)-1]
	i, err := strconv.Atoi(x)
	if err != nil {
		return key + "_1"
	}
	return strings.Join(append(a, strconv.Itoa(i+1)), "_")
}
