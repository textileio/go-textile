package core

import (
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
)

// processNode walks a file node, validating and applying a dag schema
func (t *Thread) processNode(node *schema.Node, inode ipld.Node, inbound bool) error {
	hash := inode.Cid().Hash().B58String()
	t.cafeOutbox.Add(hash, repo.CafeStoreRequest)

	if len(node.Links) == 0 {
		return t.processLink(inode, node.Pin, inbound)
	}

	for name, l := range node.Links {
		// ensure link is present
		link := schema.LinkByName(inode.Links(), name)
		if link == nil {
			return schema.ErrSchemaValidationFailed
		}

		n, err := ipfs.NodeAtLink(t.node(), link)
		if err != nil {
			return err
		}

		if err := t.processLink(n, l.Pin, inbound); err != nil {
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
func (t *Thread) processLink(inode ipld.Node, pin bool, inbound bool) error {
	hash := inode.Cid().Hash().B58String()
	t.cafeOutbox.Add(hash, repo.CafeStoreRequest)

	if schema.LinkByName(inode.Links(), FileLinkName) == nil {
		return schema.ErrSchemaValidationFailed
	}
	if schema.LinkByName(inode.Links(), DataLinkName) == nil {
		return schema.ErrSchemaValidationFailed
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
		if !t.config.IsMobile {
			t.cafeOutbox.Add(hash+"/"+DataLinkName, repo.CafeStoreRequest)
		}
	}

	return nil
}
