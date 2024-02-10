package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "demo.test/grpc-demo/internal/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// 重载 UnimplementedSayHelloServer
type server struct {
	pb.UnimplementedSayHelloServer
}

// 重写 UnimplementedSayHelloServer 中的方法
func (s *server) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	if req.RequestName != "" {
		fmt.Printf("reveive a message: %#v\n", req)
		return &pb.HelloResponse{
			ResponseMsg: "Hello " + req.RequestName,
		}, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "required parameters are missing")
}

/*--------------------------------------------------------------------------------------------*/

func main() {
	var (
		caCertPath string
		serverCertPath string
		serverKeyPath string
		port int
	)
	myFlagSet := flag.NewFlagSet("my_flagset", flag.ExitOnError)
	myFlagSet.IntVar(&port, "port", port, "server run port")
	myFlagSet.StringVar(&caCertPath, "rootca_path", caCertPath, "the server's root certificate")
	myFlagSet.StringVar(&serverCertPath, "servercert_path", serverCertPath, "server key path")
	myFlagSet.StringVar(&serverKeyPath, "serverkey_path", serverKeyPath, "server key path")
	myFlagSet.Parse(os.Args[1:])

	certificate, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		log.Fatalf("failed to load key pair, err, %v", err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("could not read ca certificate, err, %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append certs")
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate}, // 设置证书链
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    certPool,
	})
	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.MaxRecvMsgSize(1024 * 1024 * 100),
	}

	listen, err := net.Listen("tcp", ":9090") // 开启网络监听
	if err != nil {
		log.Fatalf("failed to listen, err: %v", err)
	} else {
		log.Println("listen addr is ", listen.Addr().String())
	}

	grpcServer := grpc.NewServer(opts...)            // 创建 gRPC Server
	pb.RegisterSayHelloServer(grpcServer, &server{}) // 将我们自己的 server 对象注册进 gRPC Server
	err = grpcServer.Serve(listen)                   // 启动服务
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}
