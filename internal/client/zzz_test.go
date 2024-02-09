package client_test

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"testing"

	"demo.test/grpc-demo/internal/client"
	"github.com/bxcodec/faker"
	"google.golang.org/protobuf/runtime/protoimpl"
)

func init() {
	fmt.Printf("init, flag parsed: %v\n", flag.Parsed())
}

func TestInitMyFlagSet(t *testing.T) {
	fmt.Printf("flag parsed: %v\n", flag.Parsed())

	// 为了让 Initialization 中 myFlagSet.Parse(os.Args[1:]) 可以过去。仅模拟测试使用
	os.Args = []string{""}
	client.InitMyFlagSet()
	fmt.Printf("get my flag set, httpport: %v\n", client.HttpPort)
}

func TestFaker(t *testing.T) {
	type CustomReq struct {
		state         protoimpl.MessageState
		sizeCache     protoimpl.SizeCache
		unknownFields protoimpl.UnknownFields

		RequestName string `protobuf:"bytes,1,opt,name=requestName,proto3" json:"requestName,omitempty"`
		Age         int64  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
	}

	data := CustomReq{}
	faker.FakeData(&data) // 给这个结构体赋随机值
	fmt.Printf("%+v", &data)
}

// The operation will be executed only once.
func TestOnceDo(t *testing.T) {
	once := &sync.Once{}
	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("hello world")
		})
	}
}
