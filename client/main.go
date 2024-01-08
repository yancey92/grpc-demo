package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	pb "demo.test/grpc-demo/client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("./key/cacert.pem")
	if err != nil {
		log.Fatalf("could not read ca certificate, err, %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append certs")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
		// InsecureSkipVerify: true,
	})

	// 连接到grpc server端
	conn, err := grpc.Dial("my.grpc.com:9090", grpc.WithTransportCredentials(creds))
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
