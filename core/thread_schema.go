package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"

	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
	"github.com/xeipuuv/gojsonschema"
)

// processNode walks a file node, validating and applying a dag schema
func (t *Thread) processNode(node *schema.Node, inode ipld.Node, index int, keys Keys, inbound bool) error {
	hash := inode.Cid().Hash().B58String()
	t.cafeOutbox.Add(hash, repo.CafeStoreRequest)

	if len(node.Links) == 0 {
		key := keys["/"+strconv.Itoa(index)+"/"]
		return t.processLink(inode, node.Pin, key, inbound)
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
		if err := t.processLink(n, l.Pin, key, inbound); err != nil {
			return err
		}
	}

	// pin link directory
	if node.Pin && inbound {
		if err := ipfs.PinNode(t.node(), inode, false); err != nil {
			return err
		}
	}

	go t.cafeOutbox.Flush()

	return nil
}

// processLink validates and pins file nodes
func (t *Thread) processLink(inode ipld.Node, pin bool, key string, inbound bool) error {
	hash := inode.Cid().Hash().B58String()
	t.cafeOutbox.Add(hash, repo.CafeStoreRequest)

	if schema.LinkByName(inode.Links(), FileLinkName) == nil {
		return schema.ErrFileValidationFailed
	}

	dlink := schema.LinkByName(inode.Links(), DataLinkName)
	if dlink == nil {
		return schema.ErrFileValidationFailed
	}

	if err := t.validateJsonNode(inode, key); err != nil {
		return err
	}

	// pin leaf nodes if schema dictates or files originate locally
	if pin || !inbound {
		if err := ipfs.PinNode(t.node(), inode, true); err != nil {
			return err
		}
	}

	// remote pin leaf nodes if files originate locally
	if !inbound {
		t.cafeOutbox.Add(hash+"/"+FileLinkName, repo.CafeStoreRequest)

		if !t.config.IsMobile || dlink.Size <= uint64(t.config.Cafe.Client.Mobile.P2PWireLimit) {
			t.cafeOutbox.Add(hash+"/"+DataLinkName, repo.CafeStoreRequest)
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

	jschema, err := json.Marshal(&t.Schema.JsonSchema)
	if err != nil {
		return err
	}

	sch := gojsonschema.NewStringLoader(string(jschema))
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
		return errors.New(errs)
	}

	return nil
}
