package logic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http/httptrace"
	"os"
	"strings"
	"time"

	"demo.test/grpc-demo/internal/client"
	"demo.test/grpc-demo/internal/client/models"
	pb "demo.test/grpc-demo/internal/client/proto"
	"demo.test/grpc-demo/pkg/strkit"
	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func getCertPool(rootCAPath string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	if _, err := os.Stat(rootCAPath); err != nil {
		logrus.Error(err)
		return nil, err
	}
	if ca, err := os.ReadFile(rootCAPath); err != nil {
		logrus.Errorf("could not read ca certificate, err, %v", err)
		return nil, err
	} else {
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, fmt.Errorf("failed to append certs")
		}
		return certPool, nil
	}
}

func MockRPC() (metrics *models.MetricsGRPC, err error) {
	var (
		certPool *x509.CertPool
		creds    credentials.TransportCredentials
		conn     *grpc.ClientConn
		timeout  = 20 * time.Second
	)
	if metrics == nil {
		metrics = &models.MetricsGRPC{}
	}
	if strkit.StrNotBlank(client.RootCAPath) {
		if certPool, err = getCertPool(client.RootCAPath); err != nil {
			return
		}
		creds = credentials.NewTLS(&tls.Config{
			RootCAs: certPool, InsecureSkipVerify: client.InsecureSkipVerify,
		})
	} else {
		creds = insecure.NewCredentials()
	}

	// 在这里可以自定义网络连接的创建逻辑
	customDialer := func(ctx context.Context, address string) (net.Conn, error) {

		clientTrace := &httptrace.ClientTrace{
			DNSStart: func(info httptrace.DNSStartInfo) {
				t := time.Now()
				metrics.DomainResolutionBegin = t
				metrics.HostName = info.Host
			},
			DNSDone: func(info httptrace.DNSDoneInfo) {
				t := time.Now()
				metrics.DomainResolutionEnd = t
				for _, ip := range info.Addrs {
					metrics.ResolvedIP = append(metrics.ResolvedIP, ip.IP.String())
				}
			},
			ConnectStart: func(network, addr string) {
				t := time.Now()
				metrics.ConnectedBegin = t
				metrics.UsedIP = strings.Split(addr, ":")[0]
			},
			ConnectDone: func(network, addr string, err error) {
				t := time.Now()
				metrics.ConnectedEnd = t
			},
			// httptrace 中不是所有的方法都能用的(具体请看 httptrace 的源码)，上面这4个方法也是借助 internal/nettrace.Trace 实现的，
			// 因为grpc没有对这部分做实现，并且你会发现 nettrace.Trace 是在 internal 下面的，所以外面的mod将无法直接使用。
		}
		ctx = httptrace.WithClientTrace(ctx, clientTrace)
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	// invoker interceptor
	clientInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		metrics.SendDataLength = int64(len(fmt.Sprintf("%v", req)))
		err := invoker(ctx, method, req, reply, cc, opts...)
		metrics.ReceiveDataLength = int64(len(fmt.Sprintf("%v", reply)))
		return err
	}

	// handshakeHandler := &trace.HandshakeHandler{
	// 	HandshakeStart: func() {
	// 		t := time.Now()
	// 		metrics.HandshakeBegin = t
	// 	},
	// 	HandshakeDone: func(err error) {
	// 		t := time.Now()
	// 		metrics.HandshakeEnd = t
	// 	},
	// 	TLSHandshakeStart: func(serverName string) {
	// 		t := time.Now()
	// 		metrics.TLSHandshakeBegin = t
	// 	},
	// 	TLSHandshakeDone: func(serverName string, err error) {
	// 		t := time.Now()
	// 		metrics.TLSHandshakeEnd = t
	// 	},
	// }

	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024 * 100)),
		grpc.WithUnaryInterceptor(clientInterceptor),
		grpc.WithTransportCredentials(creds),
		// grpc.WithTransportCredentialsHandler(handshakeHandler),
		grpc.WithContextDialer(customDialer), // 根据设置的超时时间，会指数回退方式尝试调用customDialer
	}
	ctxDial, cancelDial := context.WithTimeout(context.Background(), timeout)
	defer cancelDial()

	// Connect to the grpc server
	// See information on how grpc.Dial() and DialContext work
	metrics.ProcessBegin = time.Now()
	defer func() {
		metrics.ProcessEnd = time.Now()
	}()
	if conn, err = grpc.DialContext(ctxDial, client.GRPCServerAddr, dialOptions...); err != nil {
		logrus.Errorf("grpc connection err: %v", err)
		return
	}
	defer conn.Close()

	ctxCall, cancelCall := context.WithTimeout(context.Background(), timeout)
	defer cancelCall()
	// requestHandler := &trace.RequestHandler{
	// 	RequestStart: func() {
	// 		t := time.Now()
	// 		metrics.RequestBegin = t
	// 	},
	// 	RequestDone: func(err error) {
	// 		t := time.Now()
	// 		metrics.RequestEnd = t
	// 	},
	// 	ResponseStart: func() {
	// 		t := time.Now()
	// 		metrics.ResponseBegin = t
	// 	},
	// 	ResponseDone: func(err error) {
	// 		t := time.Now()
	// 		metrics.ResponseEnd = t
	// 	},
	// }
	// ctxCall = context.WithValue(ctxCall, trace.RequestHandlerContextKey{}, requestHandler)
	client := pb.NewMyCustomClient(conn)
	data := pb.MyCustomRequest{}
	faker.FakeData(&data)
	_, err = client.CustomFunc(ctxCall,
		// &pb.MyCustomRequest{Request: pb.HelloRequest{RequestName: "Alice", Age: 20}},
		&data,
	) // 通过本地调用来调用远程的方法（是不是很方便!）
	if err != nil {
		logrus.Errorf("grpc client call failed, err: %v", err)
		return
	} else {
		logrus.Info("grpc client call successful.")
	}

	return
}
