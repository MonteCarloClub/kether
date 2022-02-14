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
	"github.com/MonteCarloClub/kether/log"
)

func Deploy(ketherObject *KetherObject, ketherObjectState *KetherObjectState) error {
	// TODO ctx 应该是单子，定义在 main
	ctx := context.Background()

	imageName := ketherObject.GetImageName()
	err := container.PullDockerImage(ctx, imageName)
	if err != nil {
		ketherObjectState.SetState(FAIL_TO_DEPLOY)
		log.Error("fail to pull docker image", "imageName", imageName, "err", err)
		return err
	}
	log.Info("docker image pulled")

	containerConfig, hostConfig := ketherObject.GetContainerConfig()
	id, err := container.CreateDockerContainer(ctx, containerConfig, hostConfig)
	if id == "" {
		ketherObjectState.SetState(FAIL_TO_DEPLOY)
		err = fmt.Errorf("empty container id")
		log.Error("fail to create docker container, empty id", "err", err)
		return err
	}
	if err != nil {
		ketherObjectState.SetState(FAIL_TO_DEPLOY)
		log.Error("fail to create docker container", "id", id, "err", err)
		return err
	}
	log.Info("container created")

	err = container.RunDockerContainer(ctx, id)
	if err != nil {
		ketherObjectState.SetState(FAIL_TO_DEPLOY)
		log.Error("fail to run docker container", "err", err)
		return err
	}
	ketherObjectState.SetState(DEPLOYED)
	log.Info("container run")
	return nil
}
