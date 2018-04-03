package core

import (
	"github.com/textileio/textile-go/repo"

	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core"
	oldcmds "gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/commands"
)

var (
	VERSION = "0.0.1"
)

var Node *TextileNode

type TextileNode struct {
	// Context for issuing IPFS commands
	Context oldcmds.Context

	// IPFS node object
	IpfsNode *core.IpfsNode

	// The path to the openbazaar repo in the file system
	RepoPath string

	// Database for storing node specific data
	Datastore repo.Datastore
}
