package exec_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/liampulles/rikyu/pkg/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

func TestRunDockerContainerForOutput_WhenDockerClientFuncFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return nil, expectedErr
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenImagePullFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullErr: expectedErr,
		}, nil
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenContainerCreateFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullResp:      "image-pull",
			containerCreateErr: expectedErr,
		}, nil
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenContainerStartFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullResp:         "image-pull",
			containerCreateRespId: "id",
			containerStartErr:     expectedErr,
		}, nil
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenContainerWaitFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullResp:         "image-pull",
			containerCreateRespId: "id",
			containerWaitErr:      expectedErr,
		}, nil
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenContainerLogsForStdOutFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullResp:          "image-pull",
			containerCreateRespId:  "id",
			containerLogsStdOutErr: expectedErr,
		}, nil
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenContainerLogsForStdErrFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullResp:           "image-pull",
			containerCreateRespId:   "id",
			containerLogsStdOutResp: "stdout",
			containerLogsStdErrErr:  expectedErr,
		}, nil
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenContainerRemoveFails_ShouldFail(t *testing.T) {
	// Setup fixture
	expectedErr := fmt.Errorf("fail client")
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullResp:           "image-pull",
			containerCreateRespId:   "id",
			containerLogsStdOutResp: "stdout",
			containerLogsStdErrResp: "stderr",
			containerRemoveErr:      expectedErr,
		}, nil
	})

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", nil, nil)

	// Verify results
	if actualOut != "" {
		t.Errorf("Expected actualOut to be empty, but was %s", actualOut)
	}
	if actualErr != "" {
		t.Errorf("Expected actualErr to be empty, but was %s", actualErr)
	}
	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Error mismatch\nExpected: %v\nActual: %v", expectedErr, err)
	}
}

func TestRunDockerContainerForOutput_WhenDockerPerformsAsExpected_ShouldReturnResult(t *testing.T) {
	// Setup fixture
	expectedStdOut := "stdout"
	expectedStdErr := "stderr"
	dsi := exec.NewDockerServiceImpl(func() (exec.DockerClientWrapper, error) {
		return &mockClient{
			imagePullResp:           "image-pull",
			containerCreateRespId:   "id",
			containerLogsStdOutResp: expectedStdOut,
			containerLogsStdErrResp: expectedStdErr,
		}, nil
	})
	mounts := []exec.DockerVolumeMount{
		{
			Host:      "hostpath",
			Container: "containerpath",
		},
	}

	// Exercise SUT
	actualOut, actualErr, err := dsi.RunDockerContainerForOutput("some image", mounts, nil)

	// Verify results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualOut != expectedStdOut {
		t.Errorf("Unexpected Stdout\nExpected: %s\nActual: %s", expectedStdOut, actualOut)
	}
	if actualErr != expectedStdErr {
		t.Errorf("Unexpected Stderr\nExpected: %s\nActual: %s", expectedStdErr, actualErr)
	}
}

type mockClient struct {
	imagePullResp           string
	imagePullErr            error
	containerCreateRespId   string
	containerCreateErr      error
	containerStartErr       error
	containerWaitResp       int64
	containerWaitErr        error
	containerLogsStdOutResp string
	containerLogsStdOutErr  error
	containerLogsStdErrResp string
	containerLogsStdErrErr  error
	containerRemoveErr      error
}

func (mc *mockClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(mc.imagePullResp)), mc.imagePullErr
}

func (mc *mockClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{
		ID: mc.containerCreateRespId,
	}, mc.containerCreateErr
}

func (mc *mockClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return mc.containerStartErr
}

func (mc *mockClient) ContainerWait(ctx context.Context, containerID string) (int64, error) {
	return mc.containerWaitResp, mc.containerWaitErr
}

func (mc *mockClient) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	if options.ShowStdout && !options.ShowStderr {
		return ioutil.NopCloser(strings.NewReader(mc.containerLogsStdOutResp)), mc.containerLogsStdOutErr
	}
	if !options.ShowStdout && options.ShowStderr {
		return ioutil.NopCloser(strings.NewReader(mc.containerLogsStdErrResp)), mc.containerLogsStdErrErr
	}
	return nil, nil
}

func (mc *mockClient) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	return mc.containerRemoveErr
}
