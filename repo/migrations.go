package repo

import (
	migs "github.com/textileio/textile-go/repo/migrations"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type Migration interface {
	Up(repoPath string, dbPassword string, testnet bool) error
	Down(repoPath string, dbPassword string, testnet bool) error
}

var migrations = []Migration{
	migs.Migration000{},
	migs.Migration001{},
	migs.Migration002{},
}

// migrateUp looks at the currently active migration version and will migrate all the way up (applying all up migrations)
func migrateUp(repoPath, dbPassword string, testnet bool) error {
	version, err := ioutil.ReadFile(path.Join(repoPath, "repover"))
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil && os.IsNotExist(err) {
		version = []byte("0")
	}
	v, err := strconv.Atoi(string(version[0]))
	if err != nil {
		return err
	}
	x := v
	for _, m := range migrations[v:] {
		log.Infof("migrating repo to version %d...", x+1)
		err := m.Up(repoPath, dbPassword, testnet)
		if err != nil {
			log.Errorf("error migrating repo to version %d: %s", x+1, err)
			return err
		}
		x++
	}
	return nil
}
