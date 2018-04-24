package util

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mitchellh/go-homedir"
)

type textileSchemaManager struct {
	os       string
	dataPath string
}

// SchemaContext are the parameters which the SchemaManager derive its source of
// truth. When their zero values are provided, a reasonable default will be
// assumed during runtime.
type SchemaContext struct {
	DataPath string
	OS       string
}

// DefaultPathTransform accepts a string path representing the location where
// application data can be stored and returns a string representing the location
// where OpenBazaar prefers to store its schema on the filesystem relative to that
// path. If the path cannot be transformed, an error will be returned
func TextilePathTransform(basePath string) (path string, err error) {
	path, err = homedir.Expand(filepath.Join(basePath, directoryName()))
	if err == nil {
		path = filepath.Clean(path)
	}
	return
}

// GenerateTempPath returns a string path representing the location where
// it is okay to store temporary data. No structure or created or deleted as
// part of this operation.
func GenerateTempPath() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return filepath.Join(os.TempDir(), fmt.Sprintf("tt_tempdir_%d", r.Intn(999)))
}

// NewSchemaManager returns a service that handles the data storage directory
// required during runtime. This service also ensures no errors can be produced
// at runtime after initial creation. An error may be produced if the SchemaManager
// is unable to verify the availability of the data storage directory.
func NewSchemaManager() (*textileSchemaManager, error) {
	transformedPath, err := TextilePathTransform(defaultDataPath())
	if err != nil {
		return nil, err
	}
	return NewCustomSchemaManager(SchemaContext{
		DataPath: transformedPath,
		OS:       runtime.GOOS,
	})
}

// NewCustomSchemaManger allows a custom SchemaContext to be provided to change
func NewCustomSchemaManager(ctx SchemaContext) (*textileSchemaManager, error) {
	if len(ctx.DataPath) == 0 {
		path, err := TextilePathTransform(defaultDataPath())
		if err != nil {
			return nil, err
		}
		ctx.DataPath = path
	}
	if len(ctx.OS) == 0 {
		ctx.OS = runtime.GOOS
	}

	return &textileSchemaManager{
		dataPath: ctx.DataPath,
		os:       ctx.OS,
	}, nil
}

func defaultDataPath() (path string) {
	if runtime.GOOS == "darwin" {
		return "~/Library/Application Support"
	}
	return "~"
}

func directoryName() (directoryName string) {
	if runtime.GOOS == "linux" {
		directoryName = ".textile"
	} else {
		directoryName = "Textile"
	}

	return
}

// DataPath returns the expected location of the data storage directory
func (m *textileSchemaManager) DataPath() string { return m.dataPath }

// DatastorePath returns the expected location of the datastore file
func (m *textileSchemaManager) DatastorePath() string {
	return m.DataPathJoin("datastore", "mainnet.db")
}

// DataPathJoin is a helper function which joins the pathArgs to the service's
// dataPath and returns the result
func (m *textileSchemaManager) DataPathJoin(pathArgs ...string) string {
	allPathArgs := append([]string{m.dataPath}, pathArgs...)
	return filepath.Join(allPathArgs...)
}

// VerifySchemaVersion will ensure that the schema is currently
// the same as the expectedVersion otherwise returning an error. If the
// schema is exactly the same, nil will be returned.
func (m *textileSchemaManager) VerifySchemaVersion(expectedVersion string) error {
	schemaVersion, err := ioutil.ReadFile(m.DataPathJoin("repover"))
	if err != nil {
		return fmt.Errorf("accessing schema version: %s", err.Error())
	}
	if string(schemaVersion) != expectedVersion {
		return fmt.Errorf("schema does not match expected version %s", expectedVersion)
	}
	return nil
}

func (m *textileSchemaManager) BuildSchemaDirectories() error {
	if err := os.MkdirAll(m.DataPathJoin("datastore"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(m.DataPathJoin("logs"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(m.DataPathJoin("tmp"), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (m *textileSchemaManager) DestroySchemaDirectories() {
	os.RemoveAll(m.dataPath)
}
