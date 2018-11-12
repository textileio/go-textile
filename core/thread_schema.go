package core

import (
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
)

// process walks a file node, validating and applying a dag schema
func (t *Thread) process(dag *schema.Node, node ipld.Node, pin bool) error {
	// determine if we're at a leaf
	if len(dag.Nodes) > 0 {

		for name, ds := range dag.Nodes {
			// ensure link is present
			link := schema.LinkByName(node.Links(), name)
			if link == nil {
				return schema.ErrSchemaValidationFailed
			}

			nd, err := ipfs.LinkNode(t.node(), link)
			if err != nil {
				return err
			}

			// keep going
			if err := t.process(ds, nd, pin); err != nil {
				return err
			}
		}

		if dag.Pin && pin {
			if err := ipfs.PinNode(t.node(), node); err != nil {
				return err
			}
		}

	} else {
		hash := node.Cid().Hash().B58String()

		if dag.Pin {
			if err := ipfs.PinPath(t.node(), hash, true); err != nil {
				return err
			}
		}

		// if not mobile, remote pin the actual file data
		if !t.config.IsMobile {
			t.cafeOutbox.Add(hash, repo.CafeStoreRequest)
		}
	}

	go t.cafeOutbox.Flush()

	return nil
}
