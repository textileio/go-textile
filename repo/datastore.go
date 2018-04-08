package repo

import (
	"database/sql"
	"time"
)

type Datastore interface {
	Config() Config
	Settings() ConfigurationStore
	Photos() PhotoStore
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
	// Create the database and tables
	Init(password string) error

	// Configure the database with the node's mnemonic seed and identity key.
	Configure(mnemonic string, identityKey []byte, creationDate time.Time) error

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

type PhotoStore interface {
	Queryable

	// Put a new photo to the database
	Put(cid string, timestamp time.Time) error

	// A list of photos
	GetPhotos(offsetId string, limit int) []PhotoSet

	// Delete a photos
	DeletePhoto(cid string) error
}
