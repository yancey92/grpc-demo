package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	pb "demo.test/grpc-demo/server/proto"
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
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	if req.RequestName != "" {
		return &pb.HelloResponse{
			ResponseMsg: "Hello " + req.RequestName,
		}, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "required parameters are missing")
}

/*--------------------------------------------------------------------------------------------*/

func main() {

	certificate, err := tls.LoadX509KeyPair("./key/server.crt", "./key/server.key")
	if err != nil {
		log.Fatalf("failed to load key pair, err, %v", err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("./key/cacert.pem")
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
	}

	listen, err := net.Listen("tcp", ":9090") // 开启网络监听
	if err != nil {
		log.Fatalf("failed to listen, err: %v", err)
	}

	grpcServer := grpc.NewServer(opts...) // 创建 gRPC Server

	pb.RegisterSayHelloServer(grpcServer, &server{}) // 将我们自己的 server 对象注册进 gRPC Server

	err = grpcServer.Serve(listen) // 启动服务
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}
