package repo

import (
	"database/sql"
	"time"
)

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
	Init(password string) error
	Configure(created time.Time, version string) error
	GetCreationDate() (time.Time, error)
	GetVersion() (string, error)
	IsEncrypted() bool
}

type ProfileStore interface {
	Init(id string, secret []byte) error
	SignIn(username string, accessToken string, refreshToken string) error
	SignOut() error
	GetID() (string, error)
	GetSecret() ([]byte, error)
	GetUsername() (string, error)
	GetTokens() (accessToken string, refreshToken string, err error)
}

type ThreadStore interface {
	Queryable
	Add(thread *Thread) error
	Get(id string) *Thread
	GetByName(name string) *Thread
	List(query string) []Thread
	UpdateHead(id string, head string) error
	Delete(id string) error
	DeleteByName(name string) error
}

type BlockStore interface {
	Queryable
	Add(block *Block) error
	Get(id string) *Block
	GetByTarget(target string) *Block
	List(offsetId string, limit int, query string) []Block
	Delete(id string) error
}
