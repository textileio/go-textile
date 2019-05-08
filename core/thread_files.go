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
	target := node.Cid().Hash().B58String()
	group := cafeReqOpt.Group(target)

	// each link should point to a dag described by the thread schema
	for i, link := range node.Links() {
		nd, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return nil, err
		}
		if err := t.processFileNode(t.Schema, nd, i, keys, group, false); err != nil {
			return nil, err
		}
	}

	if err := t.cafeOutbox.Add(target, pb.CafeRequest_STORE, group); err != nil {
		return nil, err
	}

	caption = strings.TrimSpace(caption)
	msg := &pb.ThreadFiles{
		Target: target,
		Body:   caption,
		Keys:   keys,
	}

	res, err := t.commitBlock(msg, pb.Block_FILES, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, pb.Block_FILES, msg.Target, msg.Body); err != nil {
		return nil, err
	}

	for _, link := range node.Links() {
		nd, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return nil, err
		}
		if err := t.indexFileNode(nd, msg.Target); err != nil {
			return nil, err
		}
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added FILES to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleFilesBlock handles an incoming files block
func (t *Thread) handleFilesBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadFiles, error) {
	msg := new(pb.ThreadFiles)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
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
		if err := ipfs.PinNode(t.node(), node, false); err != nil {
			return nil, err
		}
		group := cafeReqOpt.Group(msg.Target)

		// each link should point to a dag described by the thread schema
		for i, link := range node.Links() {
			nd, err := ipfs.NodeAtLink(t.node(), link)
			if err != nil {
				return nil, err
			}
			if err := t.processFileNode(t.Schema, nd, i, msg.Keys, group, true); err != nil {
				return nil, err
			}
		}

		if err := t.cafeOutbox.Add(msg.Target, pb.CafeRequest_STORE, group); err != nil {
			return nil, err
		}

		// use msg keys to decrypt each file
		for pth, key := range msg.Keys {
			fd, err := ipfs.DataAtPath(t.node(), msg.Target+pth+FileLinkName)
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
			if err := jsonpb.Unmarshal(bytes.NewReader(plaintext), &file); err != nil {
				return nil, err
			}

			log.Debugf("received file: %s", file.Hash)

			if err := t.datastore.Files().Add(&file); err != nil {
				if !db.ConflictError(err) {
					return nil, err
				}
				log.Debugf("file exists: %s", file.Hash)
			}
		}
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, pb.Block_FILES, msg.Target, msg.Body); err != nil {
		return nil, err
	}

	if !ignore {
		for _, link := range node.Links() {
			nd, err := ipfs.NodeAtLink(t.node(), link)
			if err != nil {
				return nil, err
			}
			if err := t.indexFileNode(nd, msg.Target); err != nil {
				return nil, err
			}
		}
	}

	return msg, nil
}

// removeFiles unpins and removes target files unless they are used by another target,
// and unpins the target itself if not used by another block.
// TODO: Un-store on cafe(s)?
func (t *Thread) removeFiles(node ipld.Node) error {
	if node == nil {
		return ErrInvalidFileNode
	}

	target := node.Cid().Hash().B58String()

	blocks := t.datastore.Blocks().List("", -1, "target='"+target+"'").Items
	if len(blocks) == 1 {
		// safe to unpin target node

		if err := ipfs.UnpinNode(t.node(), node, false); err != nil {
			return err
		}

		// safe to dig deeper, check for other targets which contain the files
		for _, link := range node.Links() {
			nd, err := ipfs.NodeAtLink(t.node(), link)
			if err != nil {
				return err
			}
			if err := t.deIndexFileNode(nd, target); err != nil {
				return err
			}
		}
	}

	return nil
}

// processFileNode walks a file node, validating and applying a dag schema
func (t *Thread) processFileNode(node *pb.Node, inode ipld.Node, index int, keys map[string]string, group CafeRequestOption, inbound bool) error {
	hash := inode.Cid().Hash().B58String()
	if err := t.cafeOutbox.Add(hash, pb.CafeRequest_STORE, group); err != nil {
		return err
	}

	if len(node.Links) == 0 {
		key := keys["/"+strconv.Itoa(index)+"/"]
		return t.processFileLink(inode, node.Pin, node.Mill, key, group, inbound)
	}

	for name, l := range node.Links {
		// ensure link is present
		link := schema.LinkByName(inode.Links(), name)
		if link == nil {
			return schema.ErrFileValidationFailed
		}

		n, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}

		key := keys["/"+strconv.Itoa(index)+"/"+name+"/"]
		if err := t.processFileLink(n, l.Pin, l.Mill, key, group, inbound); err != nil {
			return err
		}
	}

	// pin link directory
	if node.Pin && inbound {
		if err := ipfs.PinNode(t.node(), inode, false); err != nil {
			return err
		}
	}

	return nil
}

// processFileLink validates and pins file nodes
func (t *Thread) processFileLink(inode ipld.Node, pin bool, mil string, key string, group CafeRequestOption, inbound bool) error {
	hash := inode.Cid().Hash().B58String()
	if err := t.cafeOutbox.Add(hash, pb.CafeRequest_STORE, group); err != nil {
		return err
	}

	flink := schema.LinkByName(inode.Links(), FileLinkName)
	if flink == nil {
		return ErrMissingFileLink
	}

	dlink := schema.LinkByName(inode.Links(), DataLinkName)
	if dlink == nil {
		return ErrMissingDataLink
	}

	if mil == "/json" {
		if err := t.validateJsonNode(inode, key); err != nil {
			return err
		}
	}

	// pin leaf nodes if schema dictates
	if pin {
		if err := ipfs.PinNode(t.node(), inode, true); err != nil {
			return err
		}
	}

	// remote pin leaf nodes if files originate locally
	if !inbound {
		if err := t.cafeOutbox.Add(flink.Cid.Hash().B58String(), pb.CafeRequest_STORE, group); err != nil {
			return err
		}

		if !t.config.IsMobile || dlink.Size <= uint64(t.config.Cafe.Client.Mobile.P2PWireLimit) {
			if err := t.cafeOutbox.Add(dlink.Cid.Hash().B58String(), pb.CafeRequest_STORE, group); err != nil {
				return err
			}
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

	data, err := ipfs.DataAtPath(t.node(), hash+"/"+DataLinkName)
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

		if err := t.indexFileLink(n, target); err != nil {
			return err
		}
	}

	return nil
}

// indexFileLink indexes a file link
func (t *Thread) indexFileLink(inode ipld.Node, target string) error {
	dlink := schema.LinkByName(inode.Links(), DataLinkName)
	if dlink == nil {
		return ErrMissingDataLink
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

		if err := t.deIndexFileLink(n, target); err != nil {
			return err
		}
	}

	return nil
}

// deIndexFileLink de-indexes a file link
func (t *Thread) deIndexFileLink(inode ipld.Node, target string) error {
	dlink := schema.LinkByName(inode.Links(), DataLinkName)
	if dlink == nil {
		return ErrMissingDataLink
	}

	hash := dlink.Cid.Hash().B58String()

	if err := t.datastore.Files().RemoveTarget(hash, target); err != nil {
		return err
	}

	file := t.datastore.Files().Get(hash)
	if file != nil {
		if len(file.Targets) == 0 {
			// safe to unpin and de-index

			if err := ipfs.UnpinNode(t.node(), inode, true); err != nil {
				return err
			}
			if err := t.datastore.Files().Delete(hash); err != nil {
				return err
			}
		}
	}

	return nil
}
