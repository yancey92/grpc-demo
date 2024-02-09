// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package client

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

var (
	HttpPort      int
	PprofHttpPort int

	LogLevel     int
	LogPath      string

	GRPCServerAddr     string
	InsecureSkipVerify bool
	RootCAPath         string
	ClientCertPath     string
	ClientKeyPath      string

	onceFlag = &sync.Once{}
)

func InitMyFlagSet() {

	onceFlag.Do(func() {
		myFlagSet := flag.NewFlagSet("my_flagset", flag.ExitOnError)
		// logrus.Infoln(os.Args)

		myFlagSet.IntVar(&HttpPort, "port", HttpPort, "http run port")
		myFlagSet.StringVar(&LogPath, "logpath", LogPath, "the path where the logs are stored")
		myFlagSet.IntVar(&LogLevel, "loglevel", LogLevel, "log level")
		myFlagSet.StringVar(&GRPCServerAddr, "grpc_server_addr", GRPCServerAddr, "grpc server address, e.g: my.grpc.com:9090")
		myFlagSet.BoolVar(&InsecureSkipVerify, "grpc_skip_verify", InsecureSkipVerify, "a grpc client skip verifies the server's certificate")
		myFlagSet.StringVar(&RootCAPath, "rootca_path", RootCAPath, "the server's root certificate")
		// ClientCertPath = "/client.cert"
		// ClientKeyPath = "/client.key"

		myFlagSet.Parse(os.Args[1:])
		// It can be parsed multiple times, it will overwrite the previous one
		// myFlagSet.Parse(os.Args[1:]) 

		fmt.Println("my flag set parsed successful")

	})

}




