package project

import (
	"sync"

	"github.com/liampulles/rikyu/pkg/system"
	"github.com/liampulles/rikyu/pkg/types"
)

type ProjectService interface {
	Create(name string, path string) (*types.Project, error)
	Save(project *types.Project) error
	Load(dir string) (*types.Project, error)
}

type ProjectServiceImpl struct {
	systemService system.SystemService
	ioLock        sync.Mutex
}

var _ ProjectService = &ProjectServiceImpl{}

func NewProjectServiceImpl(systemService system.SystemService) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		systemService: systemService,
	}
}
