package exec

import (
	"bytes"
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

type DockerService interface {
	RunDockerContainerForOutput(image string, mounts []DockerVolumeMount, args []string) (string, string, error)
}

type DockerClientWrapper interface {
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerWait(ctx context.Context, containerID string) (int64, error)
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
}

type DockerServiceImpl struct {
	dockerClientFunc func() (DockerClientWrapper, error)
}

var _ DockerService = &DockerServiceImpl{}

func NewDockerServiceImpl(dockerClientFunc func() (DockerClientWrapper, error)) *DockerServiceImpl {
	return &DockerServiceImpl{
		dockerClientFunc: dockerClientFunc,
	}
}

type DockerVolumeMount struct {
	Host      string
	Container string
}

func (dsi *DockerServiceImpl) RunDockerContainerForOutput(image string, mounts []DockerVolumeMount, args []string) (string, string, error) {
	// Get docker client
	cli, err := dsi.dockerClientFunc()
	if err != nil {
		return "", "", err
	}

	// Setup mounts
	var ms []mount.Mount
	for _, dvm := range mounts {
		m := mount.Mount{
			Type:   mount.TypeBind,
			Source: dvm.Host,
			Target: dvm.Container,
		}
		ms = append(ms, m)
	}
	var hostConfig *container.HostConfig
	if len(ms) != 0 {
		hostConfig = &container.HostConfig{Mounts: ms}
	}

	// Create container
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   args,
		Tty:   true,
	}, hostConfig, nil, "")
	if err != nil {
		return "", "", err
	}

	// Run container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", err
	}
	if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
		return "", "", err
	}

	// Get stdOut
	stdOut, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: false, ShowStdout: true})
	if err != nil {
		return "", "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdOut)
	stdOutStr := buf.String()

	// Get stdErr
	stdErr, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: false})
	if err != nil {
		return "", "", err
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(stdErr)
	stdErrStr := buf.String()

	return stdOutStr, stdErrStr, nil
}
