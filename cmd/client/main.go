package main

import (
	"strconv"
	"time"

	"demo.test/grpc-demo/internal/client"
	"demo.test/grpc-demo/internal/client/api"
	"demo.test/grpc-demo/internal/client/logic"
	"demo.test/grpc-demo/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	client.InitMyFlagSet()
	client.LogSet()

	// get grpc metrics as soon as start
	if err := logMetrics(); err != nil {
		panic(err)
	}

	// http
	if client.HttpPort != 0 {
		api.SetupRouter(gin.Default()).Run(":" + strconv.Itoa(client.HttpPort))
	}

}

func logMetrics() (err error) {
	var tplog = pkg.TPluginLoger{}
	tplog.OpenLoger()
	tplog.SetPrefix("gRPC")
	tplog.SetLevel(pkg.KPluginLogLevelAll)

	metrics, err := logic.MockRPC()
	if err != nil {
		logrus.Error(err)
		return
	}

	tplog.SetHost(metrics.HostName)
	tplog.SetIp(metrics.UsedIP)
	tplog.SetNValue(pkg.KIndexTotal, metrics.ProcessEnd.Sub(metrics.ProcessBegin).Milliseconds())
	tplog.SetNValue(pkg.KIndexDNS, metrics.DomainResolutionEnd.Sub(metrics.DomainResolutionBegin).Milliseconds())
	tplog.SetNValue(pkg.KIndexConn, metrics.ConnectedEnd.Sub(metrics.ConnectedBegin).Milliseconds())
	tplog.SetNValue(pkg.KIndexSSL, metrics.TLSHandshakeEnd.Sub(metrics.TLSHandshakeBegin).Milliseconds())
	tplog.SetNValue(pkg.KIndexRequest, metrics.RequestEnd.Sub(metrics.RequestBegin).Milliseconds())
	tplog.SetNValue(pkg.KIndexFirst, metrics.ResponseBegin.Sub(metrics.RequestEnd).Milliseconds())
	tplog.SetNValue(pkg.KIndexRemain, metrics.ResponseEnd.Sub(metrics.ResponseBegin).Milliseconds())
	tplog.SetNValue(pkg.KIndexRequestLength, metrics.SendDataLength)
	tplog.SetNValue(pkg.KIndexResponseLength, metrics.ReceiveDataLength)

	responseDuration := metrics.ResponseEnd.Sub(metrics.ResponseBegin).Seconds()
	if responseDuration == 0.0 {
		responseDuration = 1.0
	}
	tplog.SetNValue(pkg.KIndexSpeed,
		int64(
			pkg.CeilFloatInt(
				(float64(metrics.ReceiveDataLength))/responseDuration,
			),
		),
	)

	time.Sleep(2 * time.Second)
	return
}
