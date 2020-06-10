package project

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/liampulles/rikyu/pkg/types"
)

func (psi *ProjectServiceImpl) Create(name string, dir string) (*types.Project, error) {
	projectDir := path.Join(dir, name)
	if psi.systemService.PathExists(projectDir) {
		return nil, fmt.Errorf("desired project path already exists, so cannot use (%s)", projectDir)
	}
	project := &types.Project{
		Name: name,
		Path: projectDir,
	}
	if err := psi.Save(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (psi *ProjectServiceImpl) Save(project *types.Project) error {
	projectFilePath := projectFile(project.Path)
	data, err := serialize(project)
	fmt.Println(string(data))
	if err != nil {
		return err
	}

	psi.ioLock.Lock()
	defer psi.ioLock.Unlock()

	return psi.systemService.OverwriteFile(projectFilePath, data)
}

func (psi *ProjectServiceImpl) Load(dir string) (*types.Project, error) {
	projectFilePath := projectFile(dir)

	psi.ioLock.Lock()
	data, err := psi.systemService.ReadFile(projectFilePath)
	psi.ioLock.Unlock()
	if err != nil {
		return nil, err
	}

	return deserialize(dir, data)
}

func serialize(project *types.Project) ([]byte, error) {
	return json.MarshalIndent(project, "", "  ")
}

func deserialize(path string, data []byte) (*types.Project, error) {
	var target types.Project
	if err := json.Unmarshal(data, &target); err != nil {
		return nil, err
	}
	target.Path = path
	return &target, nil
}

func projectFile(projectPath string) string {
	return path.Join(projectPath, "project.rku")
}
