package repo

import (
	"database/sql"
	"time"
	"github.com/textileio/textile-go/wallet"
)

type SettingsData struct {
	Version *string `json:"version"`
}

type Datastore interface {
	Config() ConfigStore
	Profile() ProfileStore
	Threads() ThreadStore
	Blocks() BlockStore
	Ping() error
	Close()
}

type Queryable interface {
	BeginTransaction() (*sql.Tx, error)
	PrepareQuery(string) (*sql.Stmt, error)
	PrepareAndExecuteQuery(string, ...interface{}) (*sql.Rows, error)
	ExecuteQuery(string, ...interface{}) (sql.Result, error)
}

type ConfigStore interface {
	// Create the database and tables
	Init(password string) error

	// Configure the database
	Configure(created time.Time) error

	// Returns the date the seed was created
	GetCreationDate() (time.Time, error)

	// Returns current version of the database
	GetVersion() (string, error)

	// Returns true if the database has failed to decrypt properly ex) wrong pw
	IsEncrypted() bool
}

type ProfileStore interface {
	// Saves username, access token, and refresh token
	SignIn(username string, at string, rt string) error

	// Deletes username and jwt
	SignOut() error

	// Get username
	GetUsername() (string, error)

	// Retrieve JSON web tokens
	GetTokens() (at string, rt string, err error)
}

type ThreadStore interface {
	Queryable

	// Add a new thread
	Add(thread *wallet.Thread) error

	// Get a single thread
	Get(id string) *wallet.Thread

	// Get a single thread by name
	GetByName(name string) *wallet.Thread

	// List threads
	List(query string) []wallet.Thread

	// Update a thread's head block
	UpdateHead(id string, head string) error

	// Delete a thread
	Delete(id string) error
}

type BlockStore interface {
	Queryable

	// Add a new block
	Add(block *wallet.Block) error

	// Get a single block
	Get(id string) *wallet.Block

	// List blocks
	List(offsetId string, limit int, query string) []wallet.Block

	// Delete a block
	Delete(id string) error
}
