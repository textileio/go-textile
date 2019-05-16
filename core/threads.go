package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	libp2pc "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/mill"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/schema/textile"
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

	var schema string
	if conf.Schema != nil {
		var sjson string

		if conf.Schema.Id != "" {
			// ensure schema id is a multi hash
			_, err = mh.FromB58String(conf.Schema.Id)
			if err != nil {
				return nil, err
			}
			schema = conf.Schema.Id
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
			schema = sfile.Hash
		}

		if schema != "" {
			err = t.cafeOutbox.Add(schema, pb.CafeRequest_STORE)
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
		Schema:    schema,
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

	thrd, err := t.loadThread(model)
	if err != nil {
		return nil, err
	}

	// we join here if we're the creator
	if join {
		_, err = thrd.joinInitial()
		if err != nil {
			return nil, err
		}
	}

	go t.cafeOutbox.Flush()

	t.sendUpdate(&pb.WalletUpdate{
		Id:   thrd.Id,
		Key:  thrd.Key,
		Type: pb.WalletUpdate_THREAD_ADDED,
	})

	// invite account peers if inviter is not an account peer
	if inviteAccount {
		for _, p := range t.accountPeers() {
			_, err = thrd.AddInvite(p)
			if err != nil {
				return nil, err
			}
		}
	}

	log.Debugf("added a new thread %s with name %s", thrd.Id, conf.Name)

	return thrd, nil
}

// AddOrUpdateThread add or updates a thread directly, usually from a backup
func (t *Textile) AddOrUpdateThread(thrd *pb.Thread) error {
	// check if we're allowed to get an invite
	// Note: just using a dummy thread here because having these access+sharing
	// methods on Thread is very nice elsewhere.
	dummy := &Thread{
		initiator: thrd.Initiator,
		ttype:     thrd.Type,
		sharing:   thrd.Sharing,
		whitelist: thrd.Whitelist,
	}
	if !dummy.shareable(t.config.Account.Address, t.config.Account.Address) {
		return ErrNotShareable
	}

	sk, err := ipfs.UnmarshalPrivateKey(thrd.Sk)
	if err != nil {
		return err
	}

	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return err
	}

	nthrd := t.Thread(id.Pretty())
	if nthrd == nil {
		config := pb.AddThreadConfig{
			Key:  thrd.Key,
			Name: thrd.Name,
			Schema: &pb.AddThreadConfig_Schema{
				Id: thrd.Schema,
			},
			Type:      thrd.Type,
			Sharing:   thrd.Sharing,
			Whitelist: thrd.Whitelist,
			Force:     true,
		}

		var err error
		nthrd, err = t.AddThread(config, sk, thrd.Initiator, false, false)
		if err != nil {
			return err
		}
	}

	index := t.datastore.Blocks().Get(thrd.Head)
	if index != nil {
		// exists, abort
		log.Debugf("%s exists, aborting", thrd.Head)
		return nil
	}

	parents, err := nthrd.followParents([]string{thrd.Head})
	if err != nil {
		return err
	}
	hash, err := mh.FromB58String(thrd.Head)
	if err != nil {
		return err
	}

	_, err = nthrd.handleHead(hash, parents)
	if err != nil {
		return err
	}

	// have we joined?
	query := fmt.Sprintf("threadId='%s' and type=%d and authorId='%s'",
		nthrd.Id, pb.Block_JOIN, t.node.Identity.Pretty())
	if t.datastore.Blocks().Count(query) == 0 {
		// go ahead, invite yourself
		_, err = nthrd.join(t.node.Identity)
		if err != nil {
			return err
		}
	} else {
		// handle newly discovered peers during back prop
		err = nthrd.sendWelcome()
		if err != nil {
			return err
		}
	}

	// flush cafe queue _at the very end_
	go nthrd.cafeOutbox.Flush()

	return nil
}

// RenameThread adds an announce block to the thread w/ a new name
// Note: Only thread initiators can update the thread's name
func (t *Textile) RenameThread(id string, name string) error {
	thrd := t.Thread(id)
	if thrd == nil {
		return ErrThreadNotFound
	}
	if thrd.initiator != t.account.Address() {
		return fmt.Errorf("thread name is not writable")
	}

	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil
	}

	thrd.Name = trimmed
	err := t.datastore.Threads().UpdateName(thrd.Id, trimmed)
	if err != nil {
		return err
	}

	_, err = thrd.annouce(&pb.ThreadAnnounce{Name: trimmed})
	return err
}

// Thread get a thread by id from loaded threads
func (t *Textile) Thread(id string) *Thread {
	for _, thrd := range t.loadedThreads {
		if thrd.Id == id {
			return thrd
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
	thrd := t.Thread(id)
	if thrd == nil {
		return nil, ErrThreadNotFound
	}

	peers := &pb.PeerList{Items: make([]*pb.Peer, 0)}
	for _, tp := range thrd.Peers() {
		p := t.datastore.Peers().Get(tp.Id)
		if p != nil {
			peers.Items = append(peers.Items, p)
		}
	}

	return peers, nil
}

// RemoveThread removes a thread
func (t *Textile) RemoveThread(id string) (mh.Multihash, error) {
	var thrd *Thread
	var index int
	for i, th := range t.loadedThreads {
		if th.Id == id {
			thrd = th
			index = i
			break
		}
	}
	if thrd == nil {
		return nil, ErrThreadNotFound
	}

	// notify peers
	addr, err := thrd.leave()
	if err != nil {
		return nil, err
	}

	// delete backups
	err = t.cafeOutbox.Add(thrd.Id, pb.CafeRequest_UNSTORE_THREAD)
	if err != nil {
		return nil, err
	}

	err = t.datastore.Threads().Delete(thrd.Id)
	if err != nil {
		return nil, err
	}

	copy(t.loadedThreads[index:], t.loadedThreads[index+1:])
	t.loadedThreads[len(t.loadedThreads)-1] = nil
	t.loadedThreads = t.loadedThreads[:len(t.loadedThreads)-1]

	t.sendUpdate(&pb.WalletUpdate{
		Id:   thrd.Id,
		Key:  thrd.Key,
		Type: pb.WalletUpdate_THREAD_REMOVED,
	})

	log.Infof("removed thread %s with name %s", thrd.Id, thrd.Name)

	return addr, nil
}

// ThreadByKey get a thread by key from loaded threads
func (t *Textile) ThreadByKey(key string) *Thread {
	for _, thrd := range t.loadedThreads {
		if thrd.Key == key {
			return thrd
		}
	}
	return nil
}

// ThreadView returns a thread with expanded view properties
func (t *Textile) ThreadView(id string) (*pb.Thread, error) {
	thrd := t.Thread(id)
	if thrd == nil {
		return nil, ErrThreadNotFound
	}

	mod := t.datastore.Threads().Get(thrd.Id)
	if mod == nil {
		return nil, errThreadReload
	}

	// add extra view info
	mod.SchemaNode = thrd.Schema
	if mod.Head != "" {
		mod.HeadBlock = t.datastore.Blocks().Get(mod.Head)
		if mod.HeadBlock != nil {
			mod.HeadBlock.User = t.PeerUser(mod.HeadBlock.Author)
		}
	}
	mod.BlockCount = int32(t.datastore.Blocks().Count(fmt.Sprintf("threadId='%s'", thrd.Id)))
	mod.PeerCount = int32(len(thrd.Peers()) + 1)

	return mod, nil
}

// SnapshotThreads creates a store thread request for all threads
func (t *Textile) SnapshotThreads() error {
	var err error
	for _, thrd := range t.loadedThreads {
		err = thrd.store()
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

				thrd := new(pb.Thread)
				err = proto.Unmarshal(plaintext, thrd)
				if err != nil {
					terrCh <- err
					break
				}

				res.Id += ":" + thrd.Head
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
	thrd, err := t.AddThread(config, sk, t.account.Address(), true, false)
	if err != nil {
		return err
	}

	// add existing contacts
	for _, p := range t.datastore.Peers().List(fmt.Sprintf("address!='%s'", t.account.Address())) {
		_, err = thrd.annouce(&pb.ThreadAnnounce{Peer: p})
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
