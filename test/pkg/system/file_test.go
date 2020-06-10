package system_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/liampulles/rikyu/pkg/system"
)

func TestOverwriteFile_ShouldCallWriteFile(t *testing.T) {
	// Setup fixture
	var passedPath string
	var passedData []byte
	var passedPerm os.FileMode
	expectedErr := fmt.Errorf("some err")
	writeFile := func(filename string, data []byte, perm os.FileMode) error {
		passedPath = filename
		passedData = data
		passedPerm = perm
		return expectedErr
	}
	expectedPath := "path"
	expectedData := []byte("bytes")
	expectedPerm := os.FileMode(0644)
	ssi := system.NewSystemServiceImpl(nil, nil, nil, writeFile, nil)

	// Exercise SUT
	actualErr := ssi.OverwriteFile(expectedPath, expectedData)

	// Verify results
	if !reflect.DeepEqual(actualErr, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, actualErr)
	}
	if passedPath != expectedPath {
		t.Errorf("Path mismatch\nExpected: %s\nActual: %s", expectedPath, passedPath)
	}
	if !reflect.DeepEqual(passedData, expectedData) {
		t.Errorf("Data mismatch\nExpected: %v\nActual: %v", expectedData, passedData)
	}
	if passedPerm != expectedPerm {
		t.Errorf("Perm mismatch\nExpected: %v\nActual: %v", expectedPerm, passedPerm)
	}
}

func TestReadFile_ShouldCallReadFile(t *testing.T) {
	// Setup fixture
	var passedPath string
	expectedData := []byte("bytes")
	expectedErr := fmt.Errorf("some error")
	readFile := func(filename string) ([]byte, error) {
		passedPath = filename
		return expectedData, expectedErr
	}
	expectedPath := "path"
	ssi := system.NewSystemServiceImpl(nil, nil, nil, nil, readFile)

	// Exercise SUT
	actualData, actualErr := ssi.ReadFile(expectedPath)

	// Verify results
	if passedPath != expectedPath {
		t.Errorf("Path mismatch\nExpected: %s\nActual: %s", expectedPath, passedPath)
	}
	if !reflect.DeepEqual(actualData, expectedData) {
		t.Errorf("Data mismatch\nExpected: %v\nActual: %v", expectedData, actualData)
	}
	if !reflect.DeepEqual(actualErr, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, actualErr)
	}
}
