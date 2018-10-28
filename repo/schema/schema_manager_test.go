package schema

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestNewSchemaManagerSetsReasonableDefaults(t *testing.T) {
	subject, err := NewSchemaManager()
	if err != nil {
		t.Fatal(err)
	}
	if subject.os != runtime.GOOS {
		t.Error("expected default OS to be set via runtime.GOOS constant")
	}

	expectedDataPath := "/foobarbaz"
	subject, err = NewCustomSchemaManager(Context{
		DataPath: expectedDataPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(subject.DataPath(), expectedDataPath) != true {
		t.Errorf("expected DataPath to start with %s", expectedDataPath)
	}
}

func TestNewSchemaManagerServiceDatastorePath(t *testing.T) {
	dataPath := "/root"
	subject, err := NewCustomSchemaManager(Context{
		DataPath: dataPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedDatastorePath := "/root/datastore/mainnet.db"
	actualPath := subject.DatastorePath()
	if actualPath != expectedDatastorePath {
		t.Errorf("datastore path for test disabled was incorrect\n\texpected: %s\n\tactual: %s",
			expectedDatastorePath,
			actualPath,
		)
	}
}

func TestVerifySchemaVersion(t *testing.T) {
	var (
		expectedVersion = "123"
	)
	paths, err := NewCustomSchemaManager(Context{})
	if err != nil {
		t.Fatal(err)
	}
	if err = os.MkdirAll(paths.DataPath(), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	versionFile, err := os.Create(paths.DataPathJoin("repover"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = versionFile.Write([]byte(expectedVersion))
	if err != nil {
		t.Fatal(err)
	}
	versionFile.Close()

	if err = paths.VerifySchemaVersion(expectedVersion); err != nil {
		t.Fatal("expected schema version verification to pass with expected version. error:", err)
	}
	if err = paths.VerifySchemaVersion("anotherversion"); err == nil {
		t.Fatal("expected schema version verification to fail with wrong version. error:", err)
	}

	if err = os.Remove(paths.DataPathJoin("repover")); err != nil {
		t.Fatal(err)
	}
	if err = paths.VerifySchemaVersion("missingfile!"); err == nil {
		t.Fatal("expected schema version verification to fail when file is missing. error:", err)
	}
}

func TestBuildSchemaDirectories(t *testing.T) {
	paths, err := NewCustomSchemaManager(Context{
		DataPath: GenerateTempPath(),
	})
	err = paths.BuildSchemaDirectories()
	if err != nil {
		t.Fatal(err)
	}
	defer paths.DestroySchemaDirectories()

	checkDirectoryCreation(t, paths.DataPathJoin("logs"))
	checkDirectoryCreation(t, paths.DataPathJoin("datastore"))
}

func checkDirectoryCreation(t *testing.T, directory string) {
	f, err := os.Open(directory)
	if err != nil {
		t.Errorf("created directory %s could not be opened", directory)
	}
	fi, _ := f.Stat()
	if fi.IsDir() == false {
		t.Errorf("maybeCreateOBDirectories did not create the directory %s", directory)
	}
	if fi.Mode().String()[1:3] != "rw" {
		t.Errorf("the created directory %s is not readable and writable for the owner", directory)
	}
}
