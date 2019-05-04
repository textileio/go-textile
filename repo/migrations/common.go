package migrations

import (
	"fmt"
	"os"
	"strings"
)

// ErrorCannotMigrateDown is thrown if migrate down is called on a major migration
var ErrorCannotMigrateDown = fmt.Errorf("cannot migrate down major")

// blastRepo repo sans logs
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

// conflictError checks if a db error is from a conflict
func conflictError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
