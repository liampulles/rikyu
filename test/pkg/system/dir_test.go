package system_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/liampulles/rikyu/pkg/system"
)

func TestCreateDirRecursive_ShouldPassCorrectParamsToMkDirAll(t *testing.T) {
	// Setup fixture
	expected := fmt.Errorf("some error")
	expectedDir := "some dir"
	expectedPerm := os.ModePerm
	var passedDir string
	var passedPerm os.FileMode
	mkdirMock := func(path string, perm os.FileMode) error {
		passedDir = path
		passedPerm = perm
		return expected
	}
	ssi := system.NewSystemServiceImpl(mkdirMock, nil, nil, nil, nil)

	// Exercise SUT
	actual := ssi.CreateDirRecursive(expectedDir)

	// Verify results
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expected, actual)
	}
	if passedPerm != expectedPerm {
		t.Errorf("Perm mismatch\nExpected: %v\nActual: %v", expectedPerm, passedPerm)
	}
	if passedDir != expectedDir {
		t.Errorf("Dir mismatch\nExpected: %v\nActual: %v", expectedDir, passedDir)
	}
}

func TestGetFilesInDir_WhenCannotReadDir_ShouldReturnErr(t *testing.T) {
	// Setup fixture
	readDirMock := func(dirname string) ([]os.FileInfo, error) {
		return nil, fmt.Errorf("some error")
	}
	ssi := system.NewSystemServiceImpl(nil, nil, readDirMock, nil, nil)

	// Exercise SUT
	actual, err := ssi.GetFilesInDir("some dir")

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but none was returned")
	}
	if actual != nil {
		t.Errorf("Expected nil response, but got: %v", actual)
	}
}

func TestGetFilesInDir_WhenCanReadDir_ShouldReturnFiles(t *testing.T) {
	// Setup fixture
	readDirMock := func(dirname string) ([]os.FileInfo, error) {
		return []os.FileInfo{
			&mockFileInfo{name: "dir1", isDir: true},
			&mockFileInfo{name: "file1", isDir: false},
			&mockFileInfo{name: "dir2", isDir: true},
			&mockFileInfo{name: "file2", isDir: false},
		}, nil
	}
	ssi := system.NewSystemServiceImpl(nil, nil, readDirMock, nil, nil)
	expected := []string{"file1", "file2"}

	// Exercise SUT
	actual, err := ssi.GetFilesInDir("some dir")

	// Verify results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Response mismatch\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestPathExists_WhenCannotStatPath_ShouldReturnExists(t *testing.T) {
	// Setup fixture
	failingStat := func(name string) (os.FileInfo, error) {
		return &mockFileInfo{}, fmt.Errorf("arbitrary error")
	}
	ssi := system.NewSystemServiceImpl(nil, failingStat, nil, nil, nil)

	// Exercise SUT
	actual := ssi.PathExists("some path")

	// Verify results
	if !actual {
		t.Errorf("Expected result to be true, but was %t", actual)
	}
}

func TestPathExists_WhenCannotStatPathDueToNotExisting_ShouldReturnNotExists(t *testing.T) {
	// Setup fixture
	failingStat := func(name string) (os.FileInfo, error) {
		return &mockFileInfo{}, os.ErrNotExist
	}
	ssi := system.NewSystemServiceImpl(nil, failingStat, nil, nil, nil)

	// Exercise SUT
	actual := ssi.PathExists("some path")

	// Verify results
	if actual {
		t.Errorf("Expected result to be false, but was %t", actual)
	}
}

func TestPathExists_WhenCanStatPath_ShouldReturnExists(t *testing.T) {
	// Setup fixture
	failingStat := func(name string) (os.FileInfo, error) {
		return &mockFileInfo{}, nil
	}
	ssi := system.NewSystemServiceImpl(nil, failingStat, nil, nil, nil)

	// Exercise SUT
	actual := ssi.PathExists("some path")

	// Verify results
	if !actual {
		t.Errorf("Expected result to be true, but was %t", actual)
	}
}

type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string {
	return m.name
}

func (m *mockFileInfo) IsDir() bool {
	return m.isDir
}

func (m *mockFileInfo) Size() int64 {
	return -1
}

func (m *mockFileInfo) Mode() os.FileMode {
	return os.ModePerm
}

func (m *mockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (m *mockFileInfo) Sys() interface{} {
	return nil
}
