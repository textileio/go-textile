package core

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/mr-tron/base58/base58"
	mh "github.com/multiformats/go-multihash"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo/db"
	"github.com/textileio/go-textile/schema"
	"github.com/textileio/go-textile/util"
	"github.com/xeipuuv/gojsonschema"
)

// AddFile adds an outgoing files block
func (t *Thread) AddFiles(node ipld.Node, caption string, keys map[string]string) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.writable(t.config.Account.Address) {
		return nil, ErrNotWritable
	}

	if t.Schema == nil {
		return nil, ErrThreadSchemaRequired
	}
	if node == nil {
		return nil, ErrInvalidFileNode
	}

	caption = strings.TrimSpace(caption)
	msg := &pb.ThreadFiles{
		Target: node.Cid().Hash().B58String(),
		Body:   caption,
		Keys:   keys,
	}

	// pre-hash the block, we only want to add it if validation passes,
	// but we need the hash for the sync group
	res, err := t.commitBlock(msg, pb.Block_FILES, false, nil)
	if err != nil {
		return nil, err
	}

	// validate and apply schema directives
	err = t.processFileTarget(t.Schema, node, keys, false)
	if err != nil {
		return nil, err
	}

	// add cafe store requests for the entire graph
	err = t.cafeReqFileTarget(node, res.hash.B58String(), "")
	if err != nil {
		return nil, err
	}

	// finish adding the block
	_, err = t.addBlock(res.ciphertext, false)
	if err != nil {
		return nil, err
	}

	err = t.indexBlock(res, pb.Block_FILES, msg.Target, msg.Body)
	if err != nil {
		return nil, err
	}

	err = t.indexFileTarget(node, msg.Target)
	if err != nil {
		return nil, err
	}

	err = t.updateHead(res.hash)
	if err != nil {
		return nil, err
	}

	err = t.post(res, t.Peers())
	if err != nil {
		return nil, err
	}

	log.Debugf("added FILES to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleFilesBlock handles an incoming files block
func (t *Thread) handleFilesBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadFiles, error) {
	msg := new(pb.ThreadFiles)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return nil, err
	}

	if !t.readable(t.config.Account.Address) {
		return nil, ErrNotReadable
	}
	if !t.writable(block.Header.Address) {
		return nil, ErrNotWritable
	}

	if t.Schema == nil {
		return nil, ErrThreadSchemaRequired
	}

	var node ipld.Node

	var ignore bool
	ignored := t.datastore.Blocks().List("", -1, "target='ignore-"+hash.B58String()+"'").Items
	if len(ignored) > 0 {
		// ignore if the first (latest) ignore came after (could happen during back prop)
		if util.ProtoTsIsNewer(ignored[0].Date, block.Header.Date) {
			ignore = true
		}
	}
	if !ignore {
		target, err := cid.Parse(msg.Target)
		if err != nil {
			return nil, err
		}
		node, err = ipfs.NodeAtCid(t.node(), target)
		if err != nil {
			return nil, err
		}
		err = ipfs.PinNode(t.node(), node, false)
		if err != nil {
			return nil, err
		}

		// validate and apply schema directives
		err = t.processFileTarget(t.Schema, node, msg.Keys, true)
		if err != nil {
			return nil, err
		}

		// use msg keys to decrypt each file
		for pth, key := range msg.Keys {
			fd, err := ipfs.DataAtPath(t.node(), msg.Target+pth+MetaLinkName)
			if err != nil {
				return nil, err
			}

			var plaintext []byte
			if key != "" {
				keyb, err := base58.Decode(key)
				if err != nil {
					return nil, err
				}
				plaintext, err = crypto.DecryptAES(fd, keyb)
				if err != nil {
					return nil, err
				}
			} else {
				plaintext = fd
			}

			var file pb.FileIndex
			err = jsonpb.Unmarshal(bytes.NewReader(plaintext), &file)
			if err != nil {
				return nil, err
			}

			log.Debugf("received file: %s", file.Hash)

			err = t.datastore.Files().Add(&file)
			if err != nil {
				if !db.ConflictError(err) {
					return nil, err
				}
				log.Debugf("file exists: %s", file.Hash)
			}
		}
	}

	err = t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_FILES, msg.Target, msg.Body)
	if err != nil {
		return nil, err
	}

	if !ignore {
		err = t.indexFileTarget(node, msg.Target)
		if err != nil {
			return nil, err
		}
	}

	return msg, nil
}

// removeFiles unpins and removes target files unless they are used by another target,
// and unpins the target itself if not used by another block.
func (t *Thread) removeFiles(node ipld.Node) error {
	if node == nil {
		return ErrInvalidFileNode
	}

	target := node.Cid().Hash().B58String()
	blocks := t.datastore.Blocks().List("", -1, "target='"+target+"'").Items
	if len(blocks) == 1 { // safe to unpin target node
		err := ipfs.UnpinNode(t.node(), node, false)
		if err != nil {
			return err
		}

		// unstore on cafes
		err = t.cafeOutbox.Add(target, pb.CafeRequest_UNSTORE)
		if err != nil {
			return err
		}

		// safe to dig deeper, check for other targets which contain the files
		for _, link := range node.Links() {
			nd, err := ipfs.NodeAtLink(t.node(), link)
			if err != nil {
				return err
			}
			err = t.deIndexFileNode(nd, target)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// processFileTarget ensures each link points to a dag described by the thread schema
func (t *Thread) processFileTarget(node *pb.Node, inode ipld.Node, keys map[string]string, inbound bool) error {
	for i, link := range inode.Links() {
		nd, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}
		err = t.processFileNode(t.Schema, nd, i, keys, inbound)
		if err != nil {
			return err
		}
	}

	return nil
}

// processFileNode walks a file node, validating and applying a dag schema
func (t *Thread) processFileNode(node *pb.Node, inode ipld.Node, index int, keys map[string]string, inbound bool) error {
	if len(node.Links) == 0 {
		key := keys["/"+strconv.Itoa(index)+"/"]
		return t.processFileLink(inode, node.Pin, node.Mill, key, inbound)
	}

	for name, l := range node.Links {
		// ensure link is present
		link := schema.LinkByName(inode.Links(), []string{name})
		if link == nil {
			return schema.ErrFileValidationFailed
		}

		n, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}

		key := keys["/"+strconv.Itoa(index)+"/"+name+"/"]
		err = t.processFileLink(n, l.Pin, l.Mill, key, inbound)
		if err != nil {
			return err
		}
	}

	// pin link directory
	if node.Pin && inbound {
		err := ipfs.PinNode(t.node(), inode, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// processFileLink validates and pins file nodes
func (t *Thread) processFileLink(inode ipld.Node, pin bool, mil string, key string, inbound bool) error {
	flink := schema.LinkByName(inode.Links(), ValidMetaLinkNames)
	if flink == nil {
		return ErrMissingMetaLink
	}
	dlink := schema.LinkByName(inode.Links(), ValidContentLinkNames)
	if dlink == nil {
		return ErrMissingContentLink
	}

	if mil == "/json" {
		err := t.validateJsonNode(inode, key)
		if err != nil {
			return err
		}
	}

	// pin leaf nodes if schema dictates
	if pin {
		err := ipfs.PinNode(t.node(), inode, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateJsonNode validates the node against schema's json schema
func (t *Thread) validateJsonNode(inode ipld.Node, key string) error {
	if t.Schema.JsonSchema == nil {
		return ErrJsonSchemaRequired
	}

	hash := inode.Cid().Hash().B58String()

	data, err := ipfs.DataAtPath(t.node(), hash+"/"+ContentLinkName)
	if err != nil {
		return err
	}

	var plaintext []byte
	if key != "" {
		keyb, err := base58.Decode(key)
		if err != nil {
			return err
		}
		plaintext, err = crypto.DecryptAES(data, keyb)
		if err != nil {
			return err
		}
	} else {
		plaintext = data
	}

	jschema, err := pbMarshaler.MarshalToString(t.Schema.JsonSchema)
	if err != nil {
		return err
	}

	sch := gojsonschema.NewStringLoader(jschema)
	doc := gojsonschema.NewStringLoader(string(plaintext))

	result, err := gojsonschema.Validate(sch, doc)
	if err != nil {
		return err
	}

	if !result.Valid() {
		var errs string
		for _, err := range result.Errors() {
			errs += fmt.Sprintf("- %s\n", err)
		}
		return fmt.Errorf(errs)
	}

	return nil
}

// indexFileTarget walks a file target node, indexing file links
func (t *Thread) indexFileTarget(inode ipld.Node, target string) error {
	for _, link := range inode.Links() {
		nd, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}
		err = t.indexFileNode(nd, target)
		if err != nil {
			return err
		}
	}

	return nil
}

// indexFileNode walks a file node, indexing file links
func (t *Thread) indexFileNode(inode ipld.Node, target string) error {
	links := inode.Links()

	if looksLikeFileNode(inode) {
		return t.indexFileLink(inode, target)
	}

	for _, link := range links {
		n, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}

		err = t.indexFileLink(n, target)
		if err != nil {
			return err
		}
	}

	return nil
}

// indexFileLink indexes a file link
func (t *Thread) indexFileLink(inode ipld.Node, target string) error {
	dlink := schema.LinkByName(inode.Links(), ValidContentLinkNames)
	if dlink == nil {
		return ErrMissingContentLink
	}

	return t.datastore.Files().AddTarget(dlink.Cid.Hash().B58String(), target)
}

// deIndexFileNode walks a file node, de-indexing file links
func (t *Thread) deIndexFileNode(inode ipld.Node, target string) error {
	links := inode.Links()

	if looksLikeFileNode(inode) {
		return t.deIndexFileLink(inode, target)
	}

	for _, link := range links {
		n, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}

		err = t.deIndexFileLink(n, target)
		if err != nil {
			return err
		}
	}

	return nil
}

// deIndexFileLink de-indexes a file link
func (t *Thread) deIndexFileLink(inode ipld.Node, target string) error {
	dlink := schema.LinkByName(inode.Links(), ValidContentLinkNames)
	if dlink == nil {
		return ErrMissingContentLink
	}

	hash := dlink.Cid.Hash().B58String()

	err := t.datastore.Files().RemoveTarget(hash, target)
	if err != nil {
		return err
	}

	file := t.datastore.Files().Get(hash)
	if file != nil {
		if len(file.Targets) == 0 {
			// safe to unpin and de-index

			err = ipfs.UnpinNode(t.node(), inode, true)
			if err != nil {
				return err
			}
			err = t.datastore.Files().Delete(hash)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// cafeReqFileTarget adds a cafe requests for the target node and all children
func (t *Thread) cafeReqFileTarget(inode ipld.Node, syncGroup string, cafe string) error {
	target := inode.Cid().Hash().B58String()
	ng := cafeReqOpt.Group(ksuid.New().String())
	sg := cafeReqOpt.SyncGroup(syncGroup)
	settings := CafeRequestOptions(ng, sg, cafeReqOpt.Cafe(cafe))

	err := t.cafeOutbox.Add(target, pb.CafeRequest_STORE, settings.Options()...)
	if err != nil {
		return err
	}

	for _, link := range inode.Links() {
		nd, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}
		err = t.cafeReqFileNode(nd, settings)
		if err != nil {
			return err
		}
	}

	return nil
}

// cafeReqFileNode adds a cafe request for each link in the node
func (t *Thread) cafeReqFileNode(inode ipld.Node, settings *CafeRequestSettings) error {
	hash := inode.Cid().Hash().B58String()
	err := t.cafeOutbox.Add(hash, pb.CafeRequest_STORE, settings.Options()...)
	if err != nil {
		return err
	}

	links := inode.Links()
	if looksLikeFileNode(inode) {
		return t.cafeReqFileLink(inode, settings)
	}

	for _, l := range links {
		n, err := ipfs.NodeAtLink(t.node(), l)
		if err != nil {
			return err
		}

		err = t.cafeReqFileLink(n, settings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Thread) cafeReqFileLink(inode ipld.Node, settings *CafeRequestSettings) error {
	hash := inode.Cid().Hash().B58String()
	err := t.cafeOutbox.Add(hash, pb.CafeRequest_STORE, settings.Options()...)
	if err != nil {
		return err
	}

	links := inode.Links()
	flink := schema.LinkByName(links, ValidMetaLinkNames)
	if flink == nil {
		return ErrMissingMetaLink
	}
	dlink := schema.LinkByName(links, ValidContentLinkNames)
	if dlink == nil {
		return ErrMissingContentLink
	}

	opts := []CafeRequestOption{
		cafeReqOpt.Group(ksuid.New().String()),
		cafeReqOpt.SyncGroup(settings.SyncGroup),
		cafeReqOpt.Cafe(settings.Cafe),
	}
	err = t.cafeOutbox.Add(flink.Cid.Hash().B58String(), pb.CafeRequest_STORE, opts...)
	if err != nil {
		return err
	}
	err = t.cafeOutbox.Add(dlink.Cid.Hash().B58String(), pb.CafeRequest_STORE, opts...)
	if err != nil {
		return err
	}

	return nil
}
