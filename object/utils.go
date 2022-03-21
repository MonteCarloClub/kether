/*
Copyright (c) 2022 Zhang Zhanpeng <zhangregister@outlook.com>, Cai Dongliang <18307130121@fudan.edu.cn>, Zhong Chongpeng <1940064747@qq.com>

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
	"strings"

	kethercontainer "github.com/MonteCarloClub/kether/container"
	"github.com/MonteCarloClub/kether/log"
	"github.com/MonteCarloClub/kether/machine"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type ResourceDescriptionEntity struct {
	DockerImageRepository string `yaml:"repository"`
	DockerImageTag        string `yaml:"tag"`
}

type RunDescriptionEntity struct {
	Detach      bool     `yaml:"detach"`
	PublishList []string `yaml:"publish_list"`
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
			Detach:      ketherObjectEntity.Requirement.Detach,
			PublishList: ketherObjectEntity.Requirement.PublishList,
		},
	}
}

func (ketherObjectEntity *KetherObjectEntity) GetKetherObjectState() *KetherObjectState {
	return &KetherObjectState{
		Name: ketherObjectEntity.Name,
	}
}

func getImageName(repository string, tag string) string {
	// assert: repository != ""
	if tag == "" {
		return repository
	}
	return fmt.Sprintf("%v:%v", repository, tag)
}

func (ketherObject *KetherObject) GetImageName() string {
	candidateRepository := make([]string, 0)
	if ketherObject.Priority.DockerImageRepository != "" {
		candidateRepository = append(candidateRepository, ketherObject.Priority.DockerImageRepository)
	}
	if ketherObject.Predicate.DockerImageRepository != "" {
		candidateRepository = append(candidateRepository, ketherObject.Predicate.DockerImageRepository)
	}
	candidateTag := []string{ketherObject.Priority.DockerImageTag, ketherObject.Predicate.DockerImageTag}

	for _, repository := range candidateRepository {
		for _, tag := range candidateTag {
			candidateImageName := getImageName(repository, tag)
			if kethercontainer.CheckIfDockerImageAvailable(candidateImageName) {
				return candidateImageName
			}
		}
	}
	log.Warn("no available image name specified", "candidateRepository", candidateRepository, "candidateTag", candidateTag)
	return ""
}

func (ketherObject *KetherObject) GetContainerConfig() (*container.Config, *container.HostConfig) {
	publishList := ketherObject.Requirement.PublishList
	if len(publishList) == 0 {
		log.Info("empty publish list")
		return nil, nil
	}

	hostPortMap := make(map[string]string, len(publishList))           // 主机端口 -> 容器端口
	containerPortMap := make(map[string]nat.PortSet, len(publishList)) // 容器端口 -> 主机端口集
	for _, portPair := range publishList {
		portSlice := strings.Split(portPair, ":")
		if len(portSlice) != 2 {
			log.Warn("invalid port map", "portPair", portPair)
			continue
		}
		// 同一主机端口不能被不同容器端口映射
		if _, ok := hostPortMap[portSlice[0]]; ok {
			if hostPortMap[portSlice[0]] != portSlice[1] {
				log.Error("host port conflict", "host port", portSlice[0], "container ports", fmt.Sprintf("%v, %v...", hostPortMap[portSlice[0]], portSlice[1]), "err", fmt.Errorf("host port conflict"))
				return nil, nil
			}
		} else {
			hostPortMap[portSlice[0]] = portSlice[1]
		}
	}
	for hostPort, containerPort := range hostPortMap {
		if _, ok := containerPortMap[containerPort]; !ok {
			containerPortMap[containerPort] = make(nat.PortSet)
		}
		containerPortMap[containerPort][nat.Port(hostPort)] = struct{}{}
	}

	exposedPorts := make(nat.PortSet) // 容器端口列表
	portBindings := make(nat.PortMap) // 容器端口列表 -> 主机端口集列表
	for containerPort, hostPortSet := range containerPortMap {
		if containerPort == "" {
			log.Warn("no exposed port, ignored")
			continue
		}
		exposedPorts[nat.Port(containerPort)] = struct{}{}

		portBindingsValue := make([]nat.PortBinding, 0)
		var ifhostPortAvailable bool
		if len(hostPortSet) == 0 {
			altHostPort := machine.GetAvailableHostPort()
			if altHostPort == "" {
				log.Warn("fail to get alternate host port")
				continue
			}
			hostPortSet[nat.Port(altHostPort)] = struct{}{}
			ifhostPortAvailable = true
			log.Info("no available host port specified", "alternate host port", altHostPort)
		}
		for hostPort := range hostPortSet {
			var portBindingsValueElem nat.PortBinding
			if ifhostPortAvailable || machine.CheckIfHostPortAvailable(string(hostPort)) {
				portBindingsValueElem.HostPort = string(hostPort)
			} else {
				altHostPort := nat.Port(machine.GetAvailableHostPort())
				if altHostPort == "" {
					log.Warn("fail to get alternate host port", "unavailable host port", hostPort)
					continue
				}
				portBindingsValueElem.HostPort = string(altHostPort)
				log.Info("specified host port unavailable", "specified host port", hostPort, "alternate host port", altHostPort)
			}
			portBindingsValue = append(portBindingsValue, portBindingsValueElem)
		}
		portBindings[nat.Port(containerPort)] = portBindingsValue
	}

	containerConfig := &container.Config{
		Image:        ketherObject.GetImageName(),
		ExposedPorts: exposedPorts,
	}
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
	}

	return containerConfig, hostConfig
}

func (ketherObjectState *KetherObjectState) SetState(state KetherObjectStateType) {
	// TODO 根据 ketherObjectState 注册服务状态
	ketherObjectState.State = state
}
