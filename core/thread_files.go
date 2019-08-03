package core

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	icid "github.com/ipfs/go-cid"
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
func (t *Thread) AddFiles(node ipld.Node, target string, caption string, keys map[string]string) (mh.Multihash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

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
		Body: caption,
		Keys: keys,
	}

	// pre-hash the block, we only want to add it if validation passes,
	// but we need the hash for the sync group
	res, err := t.commitBlock(msg, pb.Block_FILES, false, nil)
	if err != nil {
		return nil, err
	}

	// validate and apply schema directives
	err = t.processFileData(t.Schema, node, keys, false)
	if err != nil {
		return nil, err
	}

	// add cafe store requests for the entire graph
	err = t.cafeReqFileData(node, res.hash.B58String(), "")
	if err != nil {
		return nil, err
	}

	// finish adding the block
	_, err = t.addBlock(res.ciphertext, false)
	if err != nil {
		return nil, err
	}

	data := node.Cid().Hash().B58String()
	err = t.indexBlock(&pb.Block{
		Id:     res.hash.B58String(),
		Thread: t.Id,
		Author: res.header.Author,
		Type:   pb.Block_FILES,
		Date:   res.header.Date,
		Target: target,
		Data:   data,
		Body:   msg.Body,
		Status: pb.Block_QUEUED,
	}, false)
	if err != nil {
		return nil, err
	}

	err = t.indexFileData(node, data)
	if err != nil {
		return nil, err
	}

	log.Debugf("added FILES to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleFilesBlock handles an incoming files block
func (t *Thread) handleFilesBlock(bnode *blockNode, block *pb.ThreadBlock) (handleResult, error) {
	var res handleResult

	msg := new(pb.ThreadFiles)
	err := ptypes.UnmarshalAny(block.Payload, msg)
	if err != nil {
		return res, err
	}

	if !t.readable(t.config.Account.Address) {
		return res, ErrNotReadable
	}
	if !t.writable(block.Header.Address) {
		return res, ErrNotWritable
	}

	if t.Schema == nil {
		return res, ErrThreadSchemaRequired
	}

	var node ipld.Node
	var data string
	if msg.Target != "" {
		data = msg.Target
	} else {
		data = bnode.data
	}

	var ignore bool
	query := fmt.Sprintf("target='%s' and type=%d", bnode.hash, pb.Block_IGNORE)
	ignored := t.datastore.Blocks().List("", -1, query).Items
	if len(ignored) > 0 {
		// ignore if the first (latest) ignore came after (could happen during back prop)
		if util.ProtoTsIsNewer(ignored[0].Date, block.Header.Date) {
			ignore = true
		}
	}
	if !ignore {
		tcid, err := icid.Parse(data)
		if err != nil {
			return res, err
		}
		node, err = ipfs.NodeAtCid(t.node(), tcid)
		if err != nil {
			return res, err
		}
		err = ipfs.PinNode(t.node(), node, false)
		if err != nil {
			return res, err
		}

		// validate and apply schema directives
		err = t.processFileData(t.Schema, node, msg.Keys, true)
		if err != nil {
			return res, err
		}

		// use msg keys to decrypt each file
		for pth, key := range msg.Keys {
			fd, err := ipfs.DataAtPath(t.node(), data+pth+MetaLinkName)
			if err != nil {
				return res, err
			}

			var plaintext []byte
			if key != "" {
				keyb, err := base58.Decode(key)
				if err != nil {
					return res, err
				}
				plaintext, err = crypto.DecryptAES(fd, keyb)
				if err != nil {
					return res, err
				}
			} else {
				plaintext = fd
			}

			var file pb.FileIndex
			err = jsonpb.Unmarshal(bytes.NewReader(plaintext), &file)
			if err != nil {
				return res, err
			}

			log.Debugf("received file: %s", file.Hash)

			err = t.datastore.Files().Add(&file)
			if err != nil {
				if !db.ConflictError(err) {
					return res, err
				}
				log.Debugf("file exists: %s", file.Hash)
			}
		}
	}

	if !ignore {
		err = t.indexFileData(node, data)
		if err != nil {
			return res, err
		}
	}

	res.oldData = msg.Target // not a typo, old target is now data
	res.body = msg.Body
	return res, nil
}

// removeFiles unpins and removes linked files unless they are used by another block
func (t *Thread) removeFiles(node ipld.Node) error {
	if node == nil {
		return ErrInvalidFileNode
	}

	data := node.Cid().Hash().B58String()
	blocks := t.datastore.Blocks().List("", -1, "data='"+data+"'").Items
	if len(blocks) == 1 { // safe to unpin data node
		err := ipfs.UnpinNode(t.node(), node, false)
		if err != nil {
			return err
		}

		// unstore on cafes
		err = t.cafeOutbox.Add(data, pb.CafeRequest_UNSTORE)
		if err != nil {
			return err
		}

		// safe to dig deeper, check for other blocks which contain the files
		for _, link := range node.Links() {
			nd, err := ipfs.NodeAtLink(t.node(), link)
			if err != nil {
				return err
			}
			err = t.deIndexFileNode(nd, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// processFileData ensures each link points to a dag described by the thread schema
func (t *Thread) processFileData(node *pb.Node, inode ipld.Node, keys map[string]string, inbound bool) error {
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

// indexFileData walks a file data node, indexing file links
func (t *Thread) indexFileData(inode ipld.Node, data string) error {
	for _, link := range inode.Links() {
		nd, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}
		err = t.indexFileNode(nd, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// indexFileNode walks a file node, indexing file links
func (t *Thread) indexFileNode(inode ipld.Node, data string) error {
	links := inode.Links()

	if looksLikeFileNode(inode) {
		return t.indexFileLink(inode, data)
	}

	for _, link := range links {
		n, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}

		err = t.indexFileLink(n, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// indexFileLink indexes a file link
func (t *Thread) indexFileLink(inode ipld.Node, data string) error {
	dlink := schema.LinkByName(inode.Links(), ValidContentLinkNames)
	if dlink == nil {
		return ErrMissingContentLink
	}

	return t.datastore.Files().AddTarget(dlink.Cid.Hash().B58String(), data)
}

// deIndexFileNode walks a file node, de-indexing file links
func (t *Thread) deIndexFileNode(inode ipld.Node, data string) error {
	links := inode.Links()

	if looksLikeFileNode(inode) {
		return t.deIndexFileLink(inode, data)
	}

	for _, link := range links {
		n, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}

		err = t.deIndexFileLink(n, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// deIndexFileLink de-indexes a file link
func (t *Thread) deIndexFileLink(inode ipld.Node, data string) error {
	dlink := schema.LinkByName(inode.Links(), ValidContentLinkNames)
	if dlink == nil {
		return ErrMissingContentLink
	}

	hash := dlink.Cid.Hash().B58String()

	err := t.datastore.Files().RemoveTarget(hash, data)
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

// cafeReqFileData adds a cafe requests for the linked node and all children
func (t *Thread) cafeReqFileData(inode ipld.Node, syncGroup string, cafe string) error {
	data := inode.Cid().Hash().B58String()
	ng := cafeReqOpt.Group(ksuid.New().String())
	sg := cafeReqOpt.SyncGroup(syncGroup)
	settings := CafeRequestOptions(ng, sg, cafeReqOpt.Cafe(cafe))

	err := t.cafeOutbox.Add(data, pb.CafeRequest_STORE, settings.Options()...)
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
