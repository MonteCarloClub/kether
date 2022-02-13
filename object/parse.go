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
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func ParseYaml(yamlPath string) (*KetherObject, *KetherObjectState, error) {
	ext := filepath.Ext(yamlPath)
	if ext != ".yaml" && ext != ".yml" {
		logrus.Warnf("%v may not be a legal yaml file", yamlPath)
	}

	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		logrus.Errorf("fail to read yaml file, yamlPath: %v, err: %v", yamlPath, err)
		return nil, nil, err
	}

	ketherObjectEntity := &KetherObjectEntity{}
	err = yaml.Unmarshal(yamlBytes, &ketherObjectEntity)
	if err != nil {
		logrus.Errorf("fail to unmarshal yaml, yamlBytes: %v, err: %v", string(yamlBytes), err)
		return nil, nil, err
	}

	ketherObject := ketherObjectEntity.GetKetherObject()
	ketherObjectState := ketherObjectEntity.GetKetherObjectState()
	ketherObjectState.SetState(REGISTERED)
	return ketherObject, ketherObjectState, nil
}
