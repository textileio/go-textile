package core

import (
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
)

type Directory struct {
	dir  uio.Directory
	node *Textile
}

func (t *Textile) NewDirectory() *Directory {
	return &Directory{dir: uio.NewDirectory(t.node.DAG), node: t}
}

func (d *Directory) AddFile(fileId string) error {
	file := d.node.datastore.Files().Get(fileId)
	if file == nil {
		return ErrFileNotFound
	}
	return ipfs.AddDirectoryLink(d.node.node, d.dir, file.Name, file.Hash)
}

func (d *Directory) Pin() (mh.Multihash, error) {
	node, err := d.dir.GetNode()
	if err != nil {
		return nil, err
	}

	// local pin
	if err := ipfs.PinDirectory(d.node.node, node); err != nil {
		return nil, err
	}

	// cafe pins
	hash := node.Cid().Hash().B58String()
	if err := d.node.cafeOutbox.Add(hash, repo.CafeStoreRequest); err != nil {
		return nil, err
	}

	return node.Cid().Hash(), nil
}
