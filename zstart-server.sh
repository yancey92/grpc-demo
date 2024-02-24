#!/bin/bash

source zcommon.sh

absolute_path=$(readlink -f "$0")  # the absolute path to the file
project_dir=$(dirname ${absolute_path})

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

