package main

import (
	"context"
	"fmt"
	"net"

	pb "demo.test/grpc-demo/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func main() {

	listen, _ := net.Listen("tcp", "localhost:9090") // 开启网络监听

	grpcServer := grpc.NewServer() // 创建 gRPC Server

	pb.RegisterSayHelloServer(grpcServer, &server{}) // 将我们自己的 server 对象注册进 gRPC Server

	err := grpcServer.Serve(listen) // 启动服务
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}
