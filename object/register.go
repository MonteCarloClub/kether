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
	"github.com/sirupsen/logrus"
)

func Register(dryRun bool, yamlPath string) {
	ketherObject, ketherObjectState := ParseYaml(yamlPath)
	if ketherObject == nil {
		logrus.Infof("fail to get kether object from yaml file, yamlPath: %v", yamlPath)
		return
	}
	if ketherObjectState == nil {
		logrus.Infof("fail to get kether object state from yaml file, yamlPath: %v", yamlPath)
		return
	}
	if dryRun {
		logrus.Infof("registering kether object in dry run mode will not change any state")
	}
	// TODO 根据 ketherObject 部署服务，根据 ketherObjectState 注册服务状态
}
