
第一步：grpc 的核心就是编写服务定义文件 *.proto
    这个demo中是 hello.proto, .proto 文件使用 Protocol Buffers (ˈproʊtəkɑːl) 语法来定义服务和消息。这个文件定义了服务的接口和通信协议，以及客户端和服务器之间交换的数据结构。

第二步：服务端编写
    1）服务定义文件写好之后，使用命令生成相应的服务端 go 代码：
        protoc --go_out=.  hello.proto
        protoc --go-grpc_out=.  hello.proto
    2）自定义一个 struct 重载自动生成的 *_grpc.pb.go 代码中的 UnimplementedSayHelloServer struct 
    3）使用重载的 struct，重写 UnimplementedSayHelloServer 的 SayHello 方法（即重写 hello.proto 中声明方法）
    4）在 main 方法中：
        创建网络监听，
        创建一个 gRPC Server，
        将我们自己的 server 对象注册进 gRPC Server，
        启动服务，gRPC Server 开始 lis.Accept，直到Stop

第三步：客户端编写
