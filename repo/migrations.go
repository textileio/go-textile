package repo

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	m "github.com/textileio/go-textile/repo/migrations"
)

// Migration performs minor up and down migrations
type Migration interface {
	Up(repoPath string, pinCode string, testnet bool) error
	Down(repoPath string, pinCode string, testnet bool) error
	Major() bool
}

// minors are current minor migrations that need to be run for lower repovers
var migrations = []Migration{
	m.Minor000{},
	m.Minor001{},
	m.Minor002{},
	m.Minor003{},
	m.Minor004{},
	m.Major005{},
	m.Minor006{},
	m.Minor007{},
	m.Minor008{},
	m.Minor009{},
	m.Minor010{},
	m.Minor011{},
	m.Minor012{},
	m.Minor013{},
	m.Minor014{},
	m.Minor015{},
	m.Minor016{},
	m.Minor017{},
}

// Stat returns whether or not there's a major migration ahead of the current repover
func Stat(repoPath string) error {
	repover, err := version(repoPath)
	if err != nil {
		return err
	}
	if len(migrations) < repover {
		return ErrRepoCorrupted
	}
	for _, migration := range migrations[repover:] {
		if migration.Major() {
			return ErrMigrationRequired
		}
	}
	return nil
}

// MigrateUp applies minor migrations all the way up to current
func MigrateUp(repoPath string, pinCode string, testnet bool) error {
	repover, err := version(repoPath)
	if err != nil {
		return err
	}
	if len(migrations) < repover {
		return ErrRepoCorrupted
	}
	x := repover
	for _, migration := range migrations[repover:] {
		log.Infof("migrating repo to version %d...", x+1)
		err := migration.Up(repoPath, pinCode, testnet)
		if err != nil {
			log.Errorf("error migrating repo to version %d: %s", x+1, err)
			return err
		}
		x++
	}
	return nil
}

// version returns repo at path's version int
func version(repoPath string) (int, error) {
	version, err := ioutil.ReadFile(path.Join(repoPath, "repover"))
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	} else if err != nil && os.IsNotExist(err) {
		version = []byte("0")
	}
	v, err := strconv.Atoi(string(version[:]))
	if err != nil {
		return 0, err
	}
	return v, nil
}
