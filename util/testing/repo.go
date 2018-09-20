package testing

import (
	"crypto/rand"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/repo/db"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"os"
	"path"
	"time"
)

// Repository represents a test (temporary/volitile) repository
type Repository struct {
	Path     string
	Password string
	DB       *db.SQLiteDatastore
}

// NewRepository creates and initializes a new temporary repository for tests
func NewRepository() (*Repository, error) {
	// Create repository object
	repository := &Repository{
		Path:     GetRepoPath(),
		Password: GetPassword(),
	}

	// Create database
	var err error
	repository.DB, err = db.Create(repository.Path, "")
	if err != nil {
		return nil, err
	}

	return repository, nil
}

// ConfigFile returns the path to the test configuration file
func (r *Repository) ConfigFile() string {
	return path.Join(r.Path, "config")
}

// RemoveRepo removes the test repository
func (r *Repository) RemoveRepo() error {
	return deleteDirectory(r.Path)
}

// RemoveRoot removes the profile json from the repository
func (r *Repository) RemoveRoot() error {
	return deleteDirectory(path.Join(r.Path, "root"))
}

// Reset sets the repo state back to a blank slate but retains keys
// Initialize the IPFS repo if it does not already exist
func (r *Repository) Reset() error {
	// Clear old root
	err := r.RemoveRoot()
	if err != nil {
		return err
	}

	// Rebuild any necessary structure
	err = repo.DoInit(r.Path, "boom", func() error {
		sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return err
		}
		if err := r.DB.Config().Init(""); err != nil {
			return err
		}
		return r.DB.Config().Configure(sk, time.Now())
	})
	if err != nil && err != repo.ErrRepoExists {
		return err
	}

	return nil
}

// MustReset sets the repo state back to a blank slate but retains keys
// It panics upon failure instead of allowing tests to continue
func (r *Repository) MustReset() {
	err := r.Reset()
	if err != nil {
		panic(err)
	}
}

func deleteDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
