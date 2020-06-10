package system

import (
	"os"
)

func (ssi *SystemServiceImpl) CreateDirRecursive(path string) error {
	return ssi.mkdirAll(path, os.ModePerm)
}

func (ssi *SystemServiceImpl) GetFilesInDir(dir string) ([]string, error) {
	fileInfos, err := ssi.readDir(dir)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		result = append(result, fileInfo.Name())
	}
	return result, nil
}

func (ssi *SystemServiceImpl) PathExists(path string) bool {
	_, err := ssi.stat(path)
	return !os.IsNotExist(err)
}
