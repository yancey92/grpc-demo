// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package client

import (
	"fmt"
	"io"
	"os"

	"demo.test/grpc-demo/pkg/strkit"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LogSet() {

	if strkit.StrNotBlank(LogPath) {
		// log rotation
		logFile := &lumberjack.Logger{
			Filename:   LogPath,
			MaxSize:    100, // megabytes
			MaxBackups: 7,
			Compress:   false,
		}
		logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.Level(LogLevel))
	logrus.SetReportCaller(true)
	fmt.Println("the logrus setup is complete")

}
