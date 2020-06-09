package system

import "os"

type SystemService interface {
	CreateDirRecursive(string) error
	GetFilesInDir(string) ([]string, error)
}

type SystemServiceImpl struct {
	// os funcs
	mkdirAll func(path string, perm os.FileMode) error
	// ioutil funcs
	readDir func(dirname string) ([]os.FileInfo, error)
}

var _ SystemService = &SystemServiceImpl{}

func NewSystemServiceImpl(
	mkdirAll func(path string, perm os.FileMode) error,
	readDir func(dirname string) ([]os.FileInfo, error)) *SystemServiceImpl {

	return &SystemServiceImpl{
		mkdirAll: mkdirAll,
		readDir:  readDir,
	}
}
