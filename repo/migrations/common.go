package migrations

import (
	"errors"
	"fmt"
	"os"
)

// ErrorCannotMigrateDown is thrown if migrate down is called on a major migration
var ErrorCannotMigrateDown = errors.New("cannot migrate down major")

// blast repo sans logs
func blastRepo(repoPath string) error {
	paths := []string{
		"blocks",
		"config",
		"datastore",
		"datastore_spec",
		"keystore",
		"repo.lock",
		"tmp",
		"version",
	}
	for _, pth := range paths {
		err := os.RemoveAll(fmt.Sprintf("%s/%s", repoPath, pth))
		if err != nil {
			return err
		}
	}
	return nil
}
