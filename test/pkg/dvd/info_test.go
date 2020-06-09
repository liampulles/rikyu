package dvd_test

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"testing"

	"github.com/liampulles/rikyu/pkg/dvd"
	"github.com/liampulles/rikyu/pkg/exec"
	"github.com/liampulles/rikyu/pkg/types"
)

func TestReadDVDInfoForDirectory_WhenDockerServiceReturnsError_ShouldReturnError(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("some error")
	mds := &mockDockerService{err: expectedErr}
	DVDServiceImpl := dvd.NewDVDServiceImpl(mds, nil)

	// Exercise SUT
	actual, err := DVDServiceImpl.ReadDVDInfoForDirectory("some dir")

	// Verify results
	if actual != nil {
		t.Errorf("Expected response to be nil, but was %v", actual)
	}
	if err == nil {
		t.Errorf("Expected an err, but was nil")
	} else if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Expected error did not match actual\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestReadDVDInfoForDirectory_WhenLsdvdReturnsUnexpectedStr_ShouldReturnError(t *testing.T) {
	// Setup fixture
	mds := &mockDockerService{stdOut: "clearly not xml"}
	DVDServiceImpl := dvd.NewDVDServiceImpl(mds, nil)

	// Exercise SUT
	actual, err := DVDServiceImpl.ReadDVDInfoForDirectory("some dir")

	// Verify results
	if actual != nil {
		t.Errorf("Expected response to be nil, but was %v", actual)
	}
	if err == nil {
		t.Errorf("Expected an err, but was nil")
	}
}

func TestReadDVDInfoForDirectory_WhenLsdvdReturnsValidXML_ShouldTransformIntoValidResult(t *testing.T) {
	// Setup fixture
	content, _ := ioutil.ReadFile(path.Join("testdata", "example.lsdvd.output"))
	mds := &mockDockerService{stdOut: string(content)}
	DVDServiceImpl := dvd.NewDVDServiceImpl(mds, nil)
	// Setup expected
	expected := &types.DVD{DiscTitle: "XFILES_S1_DISC1", TitleSets: []types.TitleSet{types.TitleSet{Titles: []types.Title{types.Title{TitleNumber: 1, Angles: 1, FPS: 25, Width: 720, Height: 576, Format: "PAL", AutomaticLetterbox: false, Length: 131759999999, AudioTracks: []types.Audio{types.Audio{Language: types.Language{Code: "cs", Name: "Czech"}, Content: "Undefined", Format: "ac3", Frequency: 48000, DynamicRangeControl: true, Channels: 2}}, SubtitleTracks: []types.Subtitle{types.Subtitle{Language: types.Language{Code: "bg", Name: "Bulgarian"}, Content: "Undefined"}}, Chapters: []types.Chapter{types.Chapter{Length: 131759999999, StartCell: 1}}, Cells: []types.Cell{types.Cell{Length: 131759999999}}, Palette: types.Palette{"108080", "eb8080", "808080", "808080", "808080", "808080", "808080", "eb8080", "808080", "b08080", "808080", "808080", "808080", "808080", "808080", "808080"}}}}}}

	// Exercise SUT
	actual, err := DVDServiceImpl.ReadDVDInfoForDirectory("some dir")

	// Verify results
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Response did not match expected\n\nExpected: %#v\n\nActual: %#v", expected, actual)
	}
}

// func TestRealDVD(t *testing.T) {
// 	ds := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
// 		return client.NewEnvClient()
// 	})
// 	DVDServiceImpl := types.NewDVDServiceImpl(ds)

// 	dvd, err := DVDServiceImpl.ReadDVDInfoForDirectory("/media/liam/XFILES_S1_DISC1/")

// 	fmt.Println(dvd)
// 	fmt.Println(err)
// }

type mockDockerService struct {
	stdOut string
	stdErr string
	err    error
}

func (mds *mockDockerService) RunDockerContainerForOutput(image string, mounts []exec.DockerVolumeMount, args []string) (string, string, error) {
	return mds.stdOut, mds.stdErr, mds.err
}
