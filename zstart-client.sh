#!/bin/bash

source zcommon.sh

protofile_path=${project_dir}/configs/
proto_out_dir=${project_dir}/internal/client/proto

render_cmd=${project_dir}/bin/render
renderval_path=${project_dir}/configs/render_value.yaml
temp_path=${project_dir}/configs/template.go.txt


## build render
rm -rf ${render_cmd}
go build -mod=vendor -o ${render_cmd} ${project_dir}/cmd/render
chmod +x ${render_cmd}


## generate gRPC client code use .proto 
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


## render the template file my_generate.go.tmpl, generate my_generate.go
${render_cmd} --renderval_path=${renderval_path}  --temp_path=${temp_path}  > ${project_dir}/internal/client/proto/my_generate.go


## build client
rm -rf ${project_dir}/bin/client
go build -mod=vendor -o ${project_dir}/bin/client ${project_dir}/cmd/client
chmod +x ${project_dir}/bin/client


## run
 ${project_dir}/bin/client \
    --rootca_path=${project_dir}/configs/key/cacert.pem \
    --port=8080 \
    --loglevel=7 \
    --grpc_server_addr=my.grpc.com:9090 \
    --grpc_skip_verify=false
