
第一步：grpc 的核心就是编写服务定义文件 *.proto
    我们可以在一个项目中编写多个 .proto 文件（便于多人协同开发）, 或者在 .proto 文件中定义多个 service 模块（便于阅读）
    本demo中是 hello.proto, .proto 文件使用 Protocol Buffers (ˈproʊtəkɑːl) 语法来定义服务和消息。该文件定义了 服务的接口 和 通信协议，以及客户端和服务器之间交换的数据结构。

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
    1）将服务定义文件hello.proto 也复制一份到客户端
    2）生成客户端gRPC代码，和在服务端的操作命令一样
        protoc --go_out=.  hello.proto
        protoc --go-grpc_out=.  hello.proto
    3）grpc.Dial 连接服务端地址，然后我们创建一个 gRPC客户端对象
    4）使用创建好的客户端对象调用本地的方法，grpc 会将请求发送给服务端对应的方法，并返回响应结果


