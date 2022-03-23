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

	"github.com/MonteCarloClub/kether/log"
	"gopkg.in/yaml.v2"
)

func ParseYaml(yamlPath string) (*KetherObject, *KetherObjectState, error) {
	ext := filepath.Ext(yamlPath)
	if ext != ".yaml" && ext != ".yml" {
		log.Warn("illegal yaml file extension", "yamlPath", yamlPath)
	}

	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Error("fail to read yaml file", "yamlPath", yamlPath, "err", err)
		return nil, nil, err
	}

	ketherObjectEntity := &KetherObjectEntity{}
	err = yaml.Unmarshal(yamlBytes, &ketherObjectEntity)
	if err != nil {
		log.Error("fail to unmarshal yaml", yamlBytes, string(yamlBytes), "err", err)
		return nil, nil, err
	}

	ketherObject := ketherObjectEntity.GetKetherObject()
	ketherObjectState := ketherObjectEntity.GetKetherObjectState()
	return ketherObject, ketherObjectState, nil
}
