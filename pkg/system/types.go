package system

import "os"

type SystemService interface {
	CreateDirRecursive(string) error
	GetFilesInDir(string) ([]string, error)
	PathExists(path string) bool
	OverwriteFile(path string, data []byte) error
	ReadFile(path string) ([]byte, error)
}

type SystemServiceImpl struct {
	// os funcs
	mkdirAll func(path string, perm os.FileMode) error
	stat     func(name string) (os.FileInfo, error)
	// ioutil funcs
	readDir   func(dirname string) ([]os.FileInfo, error)
	writeFile func(filename string, data []byte, perm os.FileMode) error
	readFile  func(filename string) ([]byte, error)
}

var _ SystemService = &SystemServiceImpl{}

func NewSystemServiceImpl(
	mkdirAll func(path string, perm os.FileMode) error,
	stat func(name string) (os.FileInfo, error),
	readDir func(dirname string) ([]os.FileInfo, error),
	writeFile func(filename string, data []byte, perm os.FileMode) error,
	readFile func(filename string) ([]byte, error)) *SystemServiceImpl {

	return &SystemServiceImpl{
		mkdirAll:  mkdirAll,
		stat:      stat,
		readDir:   readDir,
		writeFile: writeFile,
		readFile:  readFile,
	}
}
