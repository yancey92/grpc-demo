// Copyright © 2024 yangxinxin_mail@163.com. All right reserved.
// Customize or override a struct or method
// This go file should be generated from the template, Do not modify this file if you are not sure
// Generate time is 2024-02-10 03:03:42
package service
 
import (
	context "context"

	grpc "google.golang.org/grpc"
)

type MyCustomRequest struct {
	Request HelloRequest
}

type MyCustomResponse struct {
	Response *HelloResponse
}

type MyCustomClient interface {
	CustomFunc(ctx context.Context, in *MyCustomRequest, opts ...grpc.CallOption) (*MyCustomResponse, error)
}

type myclient struct {
	cli sayHelloClient // The first letter of service name is lowercase
}

func NewMyCustomClient(cc grpc.ClientConnInterface) MyCustomClient {
	c := myclient{}
	c.cli.cc = cc
	return &c
}

func (c *myclient) CustomFunc(ctx context.Context, in *MyCustomRequest, opts ...grpc.CallOption) (*MyCustomResponse, error) {
	result, err := c.cli.Hello(ctx, &in.Request, opts...)
	if err != nil {
		return nil, err
	}
	return &MyCustomResponse{Response: result}, err
}
