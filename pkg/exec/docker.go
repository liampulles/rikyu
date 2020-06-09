package exec

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

type DockerService interface {
	RunDockerContainerForOutput(image string, mounts []DockerVolumeMount, args []string) (string, string, error)
}

type DockerClientWrapper interface {
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerWait(ctx context.Context, containerID string) (int64, error)
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error
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
			Type:     mount.TypeBind,
			Source:   dvm.Host,
			Target:   dvm.Container,
			ReadOnly: false,
		}
		ms = append(ms, m)
	}
	var hostConfig *container.HostConfig
	if len(ms) != 0 {
		hostConfig = &container.HostConfig{Mounts: ms}
	}

	// Pull container
	ctx := context.Background()
	r, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return "", "", err
	}
	io.Copy(os.Stderr, r)

	// Create container
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

	// Remove container
	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		return "", "", err
	}

	return stdOutStr, stdErrStr, nil
}
