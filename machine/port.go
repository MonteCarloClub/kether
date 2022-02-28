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
package machine

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"github.com/MonteCarloClub/kether/log"
)

func CheckIfHostPortAvailable(hostPort string) bool {
	hostPortInt, err := strconv.Atoi(hostPort)
	if err != nil || hostPortInt < 1024 || hostPortInt > 49151 {
		return false
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", hostPort))
	if err != nil {
		return false
	}
	err = ln.Close()
	return err == nil
}

func GetAvailableHostPort() string {
	for i := 0; i < 1000; i++ {
		randHostPort := strconv.Itoa(8000 + rand.Intn(1000))
		if CheckIfHostPortAvailable(randHostPort) {
			return randHostPort
		}
	}
	log.Warn("fail to get available host port in [8000, 9000)")
	return ""
}
