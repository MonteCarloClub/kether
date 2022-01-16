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
package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

func InitLogger() {
	Logger = logrus.New()
	Logger.Out = os.Stdout
	Logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	}
	Logger.Info("logger inited")
}

func IfTraceOrDebug() bool {
	level := Logger.GetLevel()
	return level == logrus.DebugLevel || level == logrus.TraceLevel
}

func getFields(fieldArgs ...interface{}) logrus.Fields {
	fields := map[string]interface{}{}
	for i := 0; i < len(fieldArgs); i += 2 {
		key := fieldArgs[i].(string)
		var value interface{}
		if i+1 < len(fieldArgs) {
			value = fieldArgs[i+1]
		} else {
			value = nil
			Logger.Warn("parameter length of {Info|Warning|Error} should be odd")
		}
		fields[key] = value
	}
	return fields
}

func Info(msg string, fieldArgs ...interface{}) {
	Logger.WithFields(getFields(fieldArgs...)).Info(msg)
}

func Warn(msg string, fieldArgs ...interface{}) {
	Logger.WithFields(getFields(fieldArgs...)).Warn(msg)
}

func Error(msg string, fieldArgs ...interface{}) {
	Logger.WithFields(getFields(fieldArgs...)).Error(msg)
}
