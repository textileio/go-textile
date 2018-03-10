package repo

import (
	"database/sql"
	"time"
)

type Datastore interface {
	Config() Config
	Settings() ConfigurationStore
	Ping() error
	Close()
}

type Queryable interface {
	BeginTransaction() (*sql.Tx, error)
	PrepareQuery(string) (*sql.Stmt, error)
	PrepareAndExecuteQuery(string, ...interface{}) (*sql.Rows, error)
	ExecuteQuery(string, ...interface{}) (sql.Result, error)
}

type Config interface {
	/* Initialize the database with the node's mnemonic seed and
	   identity key. This will be called during repo init. */
	Init(mnemonic string, identityKey []byte, password string, creationDate time.Time) error

	// Return the mnemonic string
	GetMnemonic() (string, error)

	// Return the identity key
	GetIdentityKey() ([]byte, error)

	// Returns the date the seed was created
	GetCreationDate() (time.Time, error)

	// Returns true if the database has failed to decrypt properly ex) wrong pw
	IsEncrypted() bool
}

type ConfigurationStore interface {
	Queryable

	// Put settings to the database, overriding all fields
	Put(settings SettingsData) error

	// Update all non-nil fields
	Update(settings SettingsData) error

	// Return the settings object
	Get() (SettingsData, error)

	// Delete all settings data
	Delete() error
}
