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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
)

func CreateDockerContainer(ctx context.Context, containerConfig *container.Config, hostConfig *container.HostConfig) (string, error) {
	containerCreateCreatedBody, err := DockerApiClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if len(containerCreateCreatedBody.Warnings) > 0 {
		logrus.Warnf("warning(s): %v", containerCreateCreatedBody.Warnings)
	}
	if err != nil {
		logrus.Errorf("fail to create container, err: %v", err)
		return "", err
	}
	logrus.Infof("container created, id: %v", containerCreateCreatedBody.ID)
	return containerCreateCreatedBody.ID, nil
}

func RunDockerContainer(ctx context.Context, id string) error {
	err := DockerApiClient.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if err != nil {
		logrus.Errorf("fail to start container, id: %v, err: %v", id, err)
		return err
	}
	logrus.Infof("container started, id: %v", id)

	containerWaitOKBodyChan, errChan := DockerApiClient.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	select {
	case <-containerWaitOKBodyChan:
		logrus.Infof("container not running, id: %v", id)
	case err := <-errChan:
		logrus.Errorf("error encountered while container running, id: %v, err: %v", id, err)
		return err
	}
	return nil
}
