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
package object

import (
	"context"
	"fmt"

	"github.com/MonteCarloClub/kether/container"
	"github.com/MonteCarloClub/kether/flag"
	"github.com/MonteCarloClub/kether/log"
)

func Deploy(ctx context.Context, ketherObject *KetherObject, ketherObjectState *KetherObjectState) error {
	var err error
	imageName := ketherObject.GetImageName()
	containerConfig, hostConfig := ketherObject.GetContainerAndHostConfig()
	networkingConfig := ketherObject.GetNetworkingConfig()
	containerName := ketherObject.GetContainerName()

	if ctx.Value(flag.ContextKey).(flag.ContextValType).DryRun {
		log.Info("image name gotten", "imageName", imageName)
		if containerConfig != nil {
			log.Info("container config gotten", "containerConfig", containerConfig)
		}
		if hostConfig != nil {
			log.Info("host config gotten", "hostConfig", hostConfig)
		}
		if networkingConfig != nil {
			log.Info("networking config gotten", "networkingConfig", networkingConfig)
		}
		if containerName != "" {
			log.Info("container name gotten", "containerName", containerName)
		}
		log.Info("deploying kether object in dry run mode will not change any state")
		return nil
	}

	if !ketherObject.Requirement.LocalImage {
		log.Info("docker image will be pulled from remote repository", "imageName", imageName)
		err = container.PullDockerImage(ctx, imageName)
		if err != nil {
			log.Error("fail to pull docker image", "imageName", imageName, "err", err)
			ketherObjectState.SetState(ctx, FAIL_TO_DEPLOY)
			return err
		}
		log.Info("docker image pulled")
	}

	id, err := container.CreateDockerContainer(ctx, containerConfig, hostConfig, networkingConfig, ketherObject.Name)
	if id == "" {
		err = fmt.Errorf("empty container id")
		log.Error("fail to create docker container, empty id", "err", err)
		ketherObjectState.SetState(ctx, FAIL_TO_DEPLOY)
		return err
	}
	if err != nil {
		log.Error("fail to create docker container", "id", id, "err", err)
		ketherObjectState.SetState(ctx, FAIL_TO_DEPLOY)
		return err
	}
	log.Info("container created")

	if ketherObject.Requirement.Detach {
		err = container.RunDockerContainerInBackground(ctx, id)
	} else {
		err = container.RunDockerContainer(ctx, id)
	}
	if err != nil {
		log.Error("fail to run docker container in {foreground|background}", "err", err)
		ketherObjectState.SetState(ctx, FAIL_TO_DEPLOY)
		return err
	}
	log.Info("container run in {foreground|background}")
	err = ketherObjectState.SetState(ctx, DEPLOYED)
	if err != nil {
		log.Error("fail to set state of kether object", "name", ketherObjectState.Name, "state", ketherObjectState.State, "err", err)
		return err
	}
	log.Info("state of kether object set", "name", ketherObjectState.Name, "state", ketherObjectState.State)
	return nil
}
