package main

import (
	"context"
	"fmt"
	"log"

	pb "demo.test/grpc-demo/client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// 连接到grpc server端，此处禁用安全传输，没有加密和验证
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc connection err: %v", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewSayHelloClient(conn)

	// 通过本地调用来调用远程的方法（是不是很方便!）
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "Zhang三", Age: 20})
	if err != nil {
		log.Printf("grpc client call failed, err: %v", err)
	} else {
		fmt.Println(resp.GetResponseMsg())
	}

}
