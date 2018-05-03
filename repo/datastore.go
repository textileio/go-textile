package repo

import (
	"database/sql"
	"time"
)

type Datastore interface {
	Config() Config
	Settings() ConfigurationStore
	Photos() PhotoStore
	Albums() AlbumStore
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

	// Configure the database
	Configure(created time.Time) error

	// Saves username, access token, and refresh token
	SignIn(username string, at string, rt string) error

	// Deletes username and jwt
	SignOut() error

	// Get username
	GetUsername() (string, error)

	// Retrieve JSON web tokens
	GetTokens() (at string, rt string, err error)

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
	Put(set *PhotoSet) error

	// Get a single photo set
	GetPhoto(cid string) *PhotoSet

	// A list of photos
	GetPhotos(offsetId string, limit int, query string) []PhotoSet

	// Delete a photo
	DeletePhoto(cid string) error
}

type AlbumStore interface {
	Queryable

	// Put a new album to the database
	Put(album *PhotoAlbum) error

	// Get a single album
	GetAlbum(id string) *PhotoAlbum

	// Get a single album by name
	GetAlbumByName(name string) *PhotoAlbum

	// A list of albums
	GetAlbums(query string) []PhotoAlbum

	// Delete an album
	DeleteAlbum(id string) error
}
