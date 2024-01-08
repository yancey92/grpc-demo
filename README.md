
#### 基本步骤
##### 第一步：grpc 的核心就是编写服务定义文件 *.proto  
    我们可以在一个项目中编写多个 .proto 文件（便于多人协同开发）, 或者在 .proto 文件中定义多个 service 模块（便于阅读）
    本demo中是 hello.proto, .proto 文件使用 Protocol Buffers (ˈproʊtəkɑːl) 语法来定义服务和消息。该文件定义了 服务的接口 和 通信协议，以及客户端和服务器之间交换的数据结构。

##### 第二步：服务端编写
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

##### 第三步：客户端编写
    1）将服务定义文件hello.proto 也复制一份到客户端
    2）生成客户端gRPC代码，和在服务端的操作命令一样
        protoc --go_out=.  hello.proto
        protoc --go-grpc_out=.  hello.proto
    3）grpc.Dial 连接服务端地址，然后我们创建一个 gRPC客户端对象
    4）使用创建好的客户端对象调用本地的方法，grpc 会将请求发送给服务端对应的方法，并返回响应结果


#### 安全传输
    gRPC 是一个典型的C/S模型，使用 SSL/TLS 方式做服务端认证 或 双向认证。

这里在 Ubuntu 中使用 openssl 来签发证书：
* 给 root CA 生成私钥：
    ```(umask 077; dir=/usr/lib/ssl/demoCA;  openssl genrsa -out $dir/private/cakey.pem 4096) ```

* 生成自签证书(root ca cert)： 
    ``` 
    dir=/usr/lib/ssl/demoCA; \
    openssl req -new -x509 -key $dir/private/cakey.pem -out $dir/cacert.pem  -days 3650  \
    -subj  "/C=CN/ST=BeiJing/L=BeiJing/O=Personal/CN=root ca server/emailAddress=test@163.com" 
    ```

* 签发服务端证书，先临时存放到 /tmp/ssl/ 下面：
    1. 生成私钥: 
        ```(umask 077; openssl genrsa -out /tmp/ssl/server.key 2048) ```
    2. 生成证书签名请求文件:
        ```
        cat > /tmp/ssl/extfile.cnf << EOF
            
        [req_ext]
        subjectAltName = @alt_names
        
        [alt_names]
        DNS.1 = my.grpc.com
        DNS.2 = my.test.com
        IP.1 = 172.16.1.241
            
        EOF
        ```
         
        ``` 
        cat /etc/pki/tls/openssl.cnf  /tmp/ssl/extfile.cnf > /tmp/ssl/openssl.config;
        
        openssl req -new \
        -days 365 \
        -key /tmp/ssl/server.key \
        -subj  "/C=CN/ST=BeiJing/L=BeiJing/O=Personal/CN=my grpc server/emailAddress=yxx@163.com" \
        -config /tmp/ssl/openssl.config \
        -reqexts req_ext \
        -out /tmp/ssl/server.csr 
        ```
    
 * server.csr、extfile.cnf 发给root CA，签发服务端证书：
    ``` 
    openssl ca -days 365 -extensions req_ext -extfile /tmp/ssl/extfile.cnf  -in /tmp/ssl/server.csr -out /usr/lib/ssl/demoCA/certs/server.crt  
    ```
 
 * 验证域名，在客户端 /etc/hosts 中配置 ：
     ``` 
     本机IP my.grpc.com 
     本机IP my.test.com 
     ```


