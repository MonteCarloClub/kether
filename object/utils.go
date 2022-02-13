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
	"fmt"

	"github.com/MonteCarloClub/kether/machine"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

type ResourceDescriptionEntity struct {
	DockerImageRepository string `yaml:"repository"`
	DockerImageTag        string `yaml:"tag"`
}

type RunDescriptionEntity struct {
	HostPort      string `yaml:"host_port"`
	ContainerPort string `yaml:"container_port"`
}

type KetherObjectEntity struct {
	Name        string                    `yaml:"name"`
	Kind        string                    `yaml:"kind"`
	Predicate   ResourceDescriptionEntity `yaml:"predicate"`
	Priority    ResourceDescriptionEntity `yaml:"priority"`
	Requirement RunDescriptionEntity      `yaml:"requirement"`
}

// ResourceDescription 描述 Kether 对象的资源需求
type ResourceDescription ResourceDescriptionEntity

// RunDescription 描述运行 Kether 对象的需求，对应 `docker run` 的选项
type RunDescription RunDescriptionEntity

// KetherObject 是 Kether 对象
type KetherObject struct {
	Name                string
	Predicate, Priority *ResourceDescription
	Requirement         *RunDescription
}

// KetherObjectStateType Kether 对象状态类型
type KetherObjectStateType int8

const (
	FAIL_TO_DEPLOY KetherObjectStateType = -2
	UNREGISTERED   KetherObjectStateType = 0
	REGISTERED     KetherObjectStateType = 1
	DEPLOYED       KetherObjectStateType = 2
	// TODO 新增后缀状态，包含状态转换中、成功和失败，建议成功和失败的状态值互为相反数
)

// KetherObjectState 是 Kether 对象状态
type KetherObjectState struct {
	Name  string
	State KetherObjectStateType
}

func (ketherObjectEntity *KetherObjectEntity) GetKetherObject() *KetherObject {
	return &KetherObject{
		Name: ketherObjectEntity.Name,
		Predicate: &ResourceDescription{
			DockerImageRepository: ketherObjectEntity.Predicate.DockerImageRepository,
			DockerImageTag:        ketherObjectEntity.Predicate.DockerImageTag,
		},
		Priority: &ResourceDescription{
			DockerImageRepository: ketherObjectEntity.Priority.DockerImageRepository,
			DockerImageTag:        ketherObjectEntity.Priority.DockerImageTag,
		},
		Requirement: &RunDescription{
			HostPort:      ketherObjectEntity.Requirement.HostPort,
			ContainerPort: ketherObjectEntity.Requirement.ContainerPort,
		},
	}
}

func (ketherObjectEntity *KetherObjectEntity) GetKetherObjectState() *KetherObjectState {
	return &KetherObjectState{
		Name: ketherObjectEntity.Name,
	}
}

func (ketherObject *KetherObject) GetImageName() string {
	repository, tag := ketherObject.Predicate.DockerImageRepository, ketherObject.Predicate.DockerImageTag
	if repository == "" {
		repository = ketherObject.Priority.DockerImageRepository
	}
	if tag == "" {
		tag = ketherObject.Priority.DockerImageTag
	}
	return fmt.Sprintf("%v:%v", repository, tag)
}

func (ketherObject *KetherObject) GetContainerConfig() (*container.Config, *container.HostConfig) {
	hostPort, containerPort := ketherObject.Requirement.HostPort, ketherObject.Requirement.ContainerPort
	if containerPort == "" {
		logrus.Infof("no exposed port")
		return nil, nil
	}
	if !machine.CheckIfHostPortAvailable(hostPort) {
		hostPort = machine.GetAvailableHostPort()
		logrus.Infof("no available host port specified, mapped to %v", hostPort)
	}

	containerConfig := &container.Config{
		Image: ketherObject.GetImageName(),
		ExposedPorts: nat.PortSet{
			nat.Port(containerPort): struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(containerPort): []nat.PortBinding{
				{
					HostPort: hostPort,
				},
			},
		},
	}
	return containerConfig, hostConfig
}

func (ketherObjectState *KetherObjectState) SetState(state KetherObjectStateType) {
	ketherObjectState.State = state
}
