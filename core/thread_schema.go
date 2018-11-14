package core

import (
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
)

// process walks a file node, validating and applying a dag schema
func (t *Thread) process(dag *schema.Node, node ipld.Node, inbound bool) error {
	hash := node.Cid().Hash().B58String()
	t.cafeOutbox.Add(hash, repo.CafeStoreRequest)

	// determine if we're at a leaf
	if len(dag.Nodes) > 0 {

		for name, ds := range dag.Nodes {
			// ensure link is present
			link := schema.LinkByName(node.Links(), name)
			if link == nil {
				return schema.ErrSchemaValidationFailed
			}

			nd, err := ipfs.NodeAtLink(t.node(), link)
			if err != nil {
				return err
			}

			// keep going
			if err := t.process(ds, nd, inbound); err != nil {
				return err
			}
		}

		if dag.Pin && inbound {
			if err := ipfs.PinNode(t.node(), node, false); err != nil {
				return err
			}
		}

	} else {
		if schema.LinkByName(node.Links(), FileLinkName) == nil {
			return schema.ErrSchemaValidationFailed
		}
		if schema.LinkByName(node.Links(), DataLinkName) == nil {
			return schema.ErrSchemaValidationFailed
		}

		// pin leaf nodes if schema dictates or files originate locally
		if dag.Pin || !inbound {
			if err := ipfs.PinNode(t.node(), node, true); err != nil {
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
	}

	go t.cafeOutbox.Flush()

	return nil
}
