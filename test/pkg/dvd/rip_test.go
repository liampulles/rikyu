package dvd_test

import (
	"fmt"
	"testing"

	"github.com/liampulles/rikyu/pkg/dvd"
	"github.com/liampulles/rikyu/pkg/exec"
)

// func TestRip(t *testing.T) {
// 	// Setup fixture
// 	ds := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
// 		return client.NewEnvClient()
// 	})
// 	ss := system.NewSystemServiceImpl(os.MkdirAll, ioutil.ReadDir)
// 	dsi := dvd.NewDVDServiceImpl(ds, ss)

// 	// Exercise SUT
// 	err := dsi.RipTitles("/media/liam/XFILES_S1_DISC1/", map[int]string{
// 		1: "/home/liam/Development/Data/xfiles/ep1",
// 		2: "/home/liam/Development/Data/xfiles/ep2",
// 	})

// 	// Verify results
// 	if err != nil {
// 		t.Errorf("Unexpected error returned: %v", err)
// 	}
// }

func TestRipTitles_WhenCreateDirFails_ShouldFail(t *testing.T) {
	// Setup fixture
	mds := &mockDockerService{"", "", nil}
	mss := &mockSystemService{createDirRecursiveErr: fmt.Errorf("some err")}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestRipTitles_GivenValidateDirEmptyFail_ShouldFail(t *testing.T) {
	// Setup fixture
	mds := &mockDockerService{"", "", nil}
	mss := &mockSystemService{getFilesInDirErr: fmt.Errorf("some error")}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestRipTitles_WhenDestinationDirIsNotEmpty_ShouldFail(t *testing.T) {
	// Setup fixture
	mds := &mockDockerService{"", "", nil}
	mss := &mockSystemService{getFilesInDirResp: []string{"some file"}}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestRipTitles_WhenDockerFails_ShouldFail(t *testing.T) {
	// Setup fixture
	mds := &mockDockerService{"", "", fmt.Errorf("docker err")}
	mss := &mockSystemService{}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestRipTitles_WhenValidatePopulatedHasError_ShouldFail(t *testing.T) {
	// Setup fixture
	mss := &mockSystemService{}
	mds := &mockDockerServiceWhichPopulatesDir{
		errToReturn:   fmt.Errorf("some err"),
		systemService: mss,
	}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestRipTitles_WhenVobNotPopulated_ShouldFail(t *testing.T) {
	// Setup fixture
	mss := &mockSystemService{}
	mds := &mockDockerServiceWhichPopulatesDir{
		filesToPlace: []string{
			"not a vob.something",
			"but an.ifo",
		},
		systemService: mss,
	}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestRipTitles_WhenIfoNotPopulated_ShouldFail(t *testing.T) {
	// Setup fixture
	mss := &mockSystemService{}
	mds := &mockDockerServiceWhichPopulatesDir{
		filesToPlace: []string{
			"is a.vob",
			"but no ifo.somethingelse",
		},
		systemService: mss,
	}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestRipTitles_WhenVobAndIfoPopulated_ShouldPass(t *testing.T) {
	// Setup fixture
	mss := &mockSystemService{}
	mds := &mockDockerServiceWhichPopulatesDir{
		filesToPlace: []string{
			"is a.vob",
			"and an.ifo",
		},
		systemService: mss,
	}
	dsi := dvd.NewDVDServiceImpl(mds, mss)

	// Exercise SUT
	err := dsi.RipTitles("some dir", map[int]string{
		1: "some dir",
	})

	// Verify results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

type mockSystemService struct {
	createDirRecursiveErr error
	getFilesInDirResp     []string
	getFilesInDirErr      error
}

func (mss *mockSystemService) CreateDirRecursive(string) error {
	return mss.createDirRecursiveErr
}

func (mss *mockSystemService) GetFilesInDir(string) ([]string, error) {
	return mss.getFilesInDirResp, mss.getFilesInDirErr
}

func (mss *mockSystemService) PathExists(path string) bool {
	return false
}

func (mss *mockSystemService) OverwriteFile(path string, data []byte) error {
	return nil
}

func (mss *mockSystemService) ReadFile(path string) ([]byte, error) {
	return nil, nil
}

type mockDockerServiceWhichPopulatesDir struct {
	systemService *mockSystemService
	filesToPlace  []string
	errToReturn   error
}

func (mds *mockDockerServiceWhichPopulatesDir) RunDockerContainerForOutput(image string, mounts []exec.DockerVolumeMount, args []string) (string, string, error) {
	mds.systemService.getFilesInDirResp = mds.filesToPlace
	mds.systemService.getFilesInDirErr = mds.errToReturn
	return "", "", nil
}
