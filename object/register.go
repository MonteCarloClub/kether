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

	"github.com/MonteCarloClub/kether/flag"
	"github.com/MonteCarloClub/kether/log"
)

func Register(ctx context.Context, yamlPath string) (*KetherObject, *KetherObjectState, error) {
	var err error
	ketherObject, ketherObjectState, err := ParseYaml(yamlPath)
	if ketherObject == nil {
		log.Error("fail to get kether object from yaml file", "yamlPath", yamlPath)
		err = fmt.Errorf("empty ketherObject")
	} else if ketherObjectState == nil {
		log.Error("fail to get kether object state from yaml file", "yamlPath", yamlPath)
		err = fmt.Errorf("empty ketherObjectState")
	} else if err != nil {
		log.Error("fail to parse yaml file", "err", err)
	}
	if ctx.Value(flag.ContextKey).(flag.ContextValType).DryRun {
		log.Info("registering kether object in dry run mode will not change any state")
	}
	return ketherObject, ketherObjectState, err
}
