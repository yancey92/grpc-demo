#!/bin/bash

protoc_name="protoc"
protoc_gen_go_grpc_name="protoc-gen-go-grpc"
protoc_gen_go_name="protoc-gen-go"
os_name=$(uname -s)
machine_arch=$(uname -m)

if [ "$os_name" = "Linux" ]; then
    protoc_name="${protoc_name}-linux"
elif [ "$os_name" = "Darwin" ]; then
    protoc_name="${protoc_name}-osx"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}-osx"
    protoc_gen_go_name="${protoc_gen_go_name}-osx"
elif [ "$os_name" = "Windows_NT" ] || [[ "$os_name" == MINGW64_NT* ]]; then
    protoc_name="${protoc_name}-win64"
else
    echo "Unknown operating system."
    exec 1
fi

if [ "$os_name" = "Windows_NT" ] ||  [[ "$os_name" == MINGW64_NT* ]]; then
    protoc_name="${protoc_name}.exe"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}.exe"
    protoc_gen_go_name="${protoc_gen_go_name}.exe"
elif [ "$machine_arch" = "x86_64" ]; then
    protoc_name="${protoc_name}-x86_64"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}-x86_64"
    protoc_gen_go_name="${protoc_gen_go_name}-x86_64"
elif [ "$machine_arch" = "arm64" ] || [ "$machine_arch" = "aarch64" ]; then
    protoc_name="${protoc_name}-arm64"
else
    echo "Unknown architecture."
    exec 1
fi


absolute_path=$(readlink -f "$0")  # the absolute path to the file
project_dir=$(dirname $(dirname ${absolute_path}))

protofile_path=${project_dir}/configs/
proto_out_dir=${project_dir}/internal/server/proto


## generate gRPC code use .proto 
if [ $(ls ${project_dir}/configs | grep .proto | wc -l ) -gt 1 ]; then
    echo "The .proto files count is greater than 1, Please keep only one"
    exit 1
fi
protofile_name=$(find ${project_dir}/configs -maxdepth 1 -type f -name "*.proto" -exec basename {} \;)
${project_dir}/bin/${protoc_name} \
    --plugin=protoc-gen-go-grpc=${project_dir}/bin/${protoc_gen_go_grpc_name} \
    --plugin=protoc-gen-go=${project_dir}/bin/${protoc_gen_go_name} \
    --go-grpc_out=${proto_out_dir} --proto_path=${protofile_path} ${protofile_name}  
${project_dir}/bin/${protoc_name} \
    --plugin=protoc-gen-go-grpc=${project_dir}/bin/${protoc_gen_go_grpc_name} \
    --plugin=protoc-gen-go=${project_dir}/bin/${protoc_gen_go_name} \
    --go_out=${proto_out_dir} --proto_path=${protofile_path}  ${protofile_name}    


## build server
rm -rf ${project_dir}/bin/server
go build -mod=vendor -o ${project_dir}/bin/server ${project_dir}/cmd/server
chmod +x ${project_dir}/bin/server


## run
 ${project_dir}/bin/server \
    --rootca_path=${project_dir}/configs/key/cacert.pem \
    --servercert_path=${project_dir}/configs/key/server.crt \
    --serverkey_path=${project_dir}/configs/key/server.key \
    --port=9090

