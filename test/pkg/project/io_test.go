package project_test

import (
	"fmt"
	"path"
	"reflect"
	"testing"

	"github.com/liampulles/rikyu/pkg/types"

	"github.com/liampulles/rikyu/pkg/project"
)

func TestCreate_WhenProjectFileAlreadyExists_ShouldFail(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp: true,
	})

	// Exercise SUT
	actual, err := psi.Create("some name", "some dir")

	// Verify results
	if actual != nil {
		t.Errorf("Expected nil response, but was %v", actual)
	}
	if err == nil {
		t.Errorf("Expected an err, but got nil")
	}
}

func TestCreate_WhenCannotSave_ShouldFail(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp:   false,
		overWriteFileErr: fmt.Errorf("some error"),
	})

	// Exercise SUT
	actual, err := psi.Create("some name", "some dir")

	// Verify results
	if actual != nil {
		t.Errorf("Expected nil response, but was %v", actual)
	}
	if err == nil {
		t.Errorf("Expected an err, but got nil")
	}
}

func TestCreate_WhenCanSave_ShouldReturnNewProject(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp:   false,
		overWriteFileErr: nil,
	})
	expected := &types.Project{
		Name: "some name",
		Path: path.Join("some dir", "some name"),
	}

	// Exercise SUT
	actual, err := psi.Create("some name", "some dir")

	// Verify results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Result mismatch\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestSave_WhenCannotOverwriteFile_ShouldFail(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp:   false,
		overWriteFileErr: fmt.Errorf("some error"),
	})

	// Exercise SUT
	err := psi.Save(&types.Project{
		Name: "some name",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected an err, but got nil")
	}
}

func TestSave_WhenCanOverwriteFile_ShouldPass(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp:   false,
		overWriteFileErr: nil,
	})

	// Exercise SUT
	err := psi.Save(&types.Project{
		Name: "some name",
	})

	// Verify results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestLoad_WhenCannotReadFile_ShouldFail(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp: false,
		readFileErr:    fmt.Errorf("some err"),
	})

	// Exercise SUT
	actual, err := psi.Load("some path")

	// Verify results
	if err == nil {
		t.Errorf("Expected an err, but got nil")
	}
	if actual != nil {
		t.Errorf("Expected response to be nil, but was %v", actual)
	}
}

func TestLoad_WhenFileContainsInvalidJSON_ShouldFail(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp: false,
		readFileResp:   []byte("#1 not json !?"),
	})

	// Exercise SUT
	actual, err := psi.Load("some path")

	// Verify results
	if err == nil {
		t.Errorf("Expected an err, but got nil")
	}
	if actual != nil {
		t.Errorf("Expected response to be nil, but was %v", actual)
	}
}

func TestLoad_WhenFileContainsValidJSON_ShouldReturnParsedProject(t *testing.T) {
	// Setup fixture
	psi := project.NewProjectServiceImpl(&mockSystemService{
		pathExistsResp: false,
		readFileResp:   []byte("{\"name\": \"some name\"}"),
	})
	expected := &types.Project{
		Name: "some name",
		Path: "some path",
	}

	// Exercise SUT
	actual, err := psi.Load("some path")

	// Verify results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected response\nExpected: %v\nActual: %v", expected, actual)
	}
}

type mockSystemService struct {
	pathExistsResp   bool
	overWriteFileErr error
	readFileResp     []byte
	readFileErr      error
}

func (mss *mockSystemService) CreateDirRecursive(string) error {
	return nil
}

func (mss *mockSystemService) GetFilesInDir(string) ([]string, error) {
	return nil, nil
}

func (mss *mockSystemService) PathExists(path string) bool {
	return mss.pathExistsResp
}

func (mss *mockSystemService) OverwriteFile(path string, data []byte) error {
	return mss.overWriteFileErr
}

func (mss *mockSystemService) ReadFile(path string) ([]byte, error) {
	return mss.readFileResp, mss.readFileErr
}
