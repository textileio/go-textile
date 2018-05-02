package test

import (
	"os"
)

// GetRepoPath returns the repo path to use for tests
// It should be considered volitile and may be destroyed at any time
func GetRepoPath() string {
	return getEnvString("TEXTILE_TEST_REPO_PATH", "/tmp/textile-test")
}

// GetPassword returns a static mneumonic to use
func GetPassword() string {
	return getEnvString("TEXTILE_TEST_PASSWORD", "correct horse battery staple")
}

func getEnvString(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
