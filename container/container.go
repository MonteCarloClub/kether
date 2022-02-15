/*
Copyright (c) 2022 Zhang Zhanpeng <zhangregister@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package container

import (
	"context"

	"github.com/MonteCarloClub/kether/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func CreateDockerContainer(ctx context.Context, containerConfig *container.Config, hostConfig *container.HostConfig) (string, error) {
	containerCreateCreatedBody, err := DockerApiClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if len(containerCreateCreatedBody.Warnings) > 0 {
		log.Warn("ContainerCreate returned warning(s)", "warning", containerCreateCreatedBody.Warnings)
	}
	if err != nil {
		log.Error("fail to create container", "err", err)
		return "", err
	}
	log.Info("container created", "id", containerCreateCreatedBody.ID)
	return containerCreateCreatedBody.ID, nil
}

func RunDockerContainer(ctx context.Context, id string) error {
	err := DockerApiClient.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if err != nil {
		log.Error("fail to start container", "id", id, "err", err)
		return err
	}
	log.Info("container started", "id", id)

	statusCh, errChan := DockerApiClient.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	select {
	case <-statusCh:
		log.Info("container not running", "id", id)
	case err := <-errChan:
		log.Error("error encountered while container running", "id", id, "err", err)
		return err
	}

	// TODO 输出容器日志，参考：https://docs.docker.com/engine/api/sdk/examples/
	return nil
}
