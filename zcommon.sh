#!/bin/bash

protoc_name="protoc"
protoc_gen_go_grpc_name="protoc-gen-go-grpc"
protoc_gen_go_name="protoc-gen-go"
os_name=$(uname -s)
machine_arch=$(uname -m)

if [ "$os_name" = "Linux" ]; then
    protoc_name="${protoc_name}-linux"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}-linux"
    protoc_gen_go_name="${protoc_gen_go_name}-linux"
elif [ "$os_name" = "Darwin" ]; then
    protoc_name="${protoc_name}-osx"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}-osx"
    protoc_gen_go_name="${protoc_gen_go_name}-osx"
elif [ "$os_name" = "Windows_NT" ] || [[ "$os_name" == MINGW64_NT* ]]; then
    protoc_name="${protoc_name}-win"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}-win"
    protoc_gen_go_name="${protoc_gen_go_name}-win"
else
    echo "Unknown operating system."
    exec 1
fi

if [ "$machine_arch" = "x86_64" ] || [ "$machine_arch" = "amd64" ]; then
    protoc_name="${protoc_name}-amd64"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}-amd64"
    protoc_gen_go_name="${protoc_gen_go_name}-amd64"
elif [ "$machine_arch" = "arm64" ] || [ "$machine_arch" = "aarch64" ]; then
    protoc_name="${protoc_name}-arm64"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}-arm64"
    protoc_gen_go_name="${protoc_gen_go_name}-arm64"
else
    echo "Unknown architecture."
    exec 1
fi
if [ "$os_name" = "Windows_NT" ] ||  [[ "$os_name" == MINGW64_NT* ]]; then
    protoc_name="${protoc_name}.exe"
    protoc_gen_go_grpc_name="${protoc_gen_go_grpc_name}.exe"
    protoc_gen_go_name="${protoc_gen_go_name}.exe"
fi

absolute_path=$(readlink -f "$0")  # the absolute path to the file
project_dir=$(dirname ${absolute_path})

