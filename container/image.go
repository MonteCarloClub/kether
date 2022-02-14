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
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/MonteCarloClub/kether/log"
	"github.com/docker/docker/api/types"
)

func PullDockerImage(ctx context.Context, imageName string) error {
	var err error
	if imageName == "" {
		err = fmt.Errorf("empty image name")
		log.Error("empty image name", "err", err)
		return err
	}

	readCloser, err := DockerApiClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	var dst io.Writer
	if log.IfTraceOrDebug() {
		dst = os.Stdout
	} else {
		dst = ioutil.Discard
	}
	io.Copy(dst, readCloser)
	if err != nil {
		log.Error("fail to pull docker image", "refStr", imageName, "err", err)
		return err
	}
	log.Info("docker image pulled", "refStr", imageName)
	return nil
}
