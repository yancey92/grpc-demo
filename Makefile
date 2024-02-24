
project_dir= .
render_cmd = ${project_dir}/bin/render
client_cmd = ${project_dir}/bin/client
server_cmd = ${project_dir}/bin/server
protofile_path=${project_dir}/configs/
protofile_name=$(shell find ${project_dir}/configs -maxdepth 1 -type f -name "*.proto" -exec basename {} \;)

protoc_name = protoc
protoc_gen_go_grpc_name = protoc-gen-go-grpc
protoc_gen_go_name = protoc-gen-go


#-------------------------------------------------------------------------------------
# Supported operating systems and architectures
os_supported = Linux Darwin Windows_NT MINGW64_NT
os_name = $(shell uname -s)

arch_supported = x86_64 amd64 aarch64 arm64
machine_arch = $(shell uname -m)

# os
ifneq ($(findstring $(os_name),$(os_supported)),)
	ifeq ($(os_name),Linux)
		protoc_name:=${protoc_name}-linux
		protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}-linux
		protoc_gen_go_name:=${protoc_gen_go_name}-linux
	endif
	ifeq ($(os_name),Darwin)
		protoc_name:=${protoc_name}-osx
		protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}-osx
		protoc_gen_go_name:=${protoc_gen_go_name}-osx
	endif
	ifeq ($(os_name),Windows_NT)
		protoc_name:=${protoc_name}-win
		protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}-win
		protoc_gen_go_name:=${protoc_gen_go_name}-win
	endif
	ifneq ($(findstring MINGW64_NT,$(os_name)),)
		protoc_name:=${protoc_name}-win
		protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}-win
		protoc_gen_go_name:=${protoc_gen_go_name}-win
	endif
else 
	$(error Unknown operating system.)
endif

# arch
ifneq ($(findstring $(machine_arch),$(arch_supported)),)
	ifneq ($(findstring $(machine_arch),aarch64 arm64),)
		protoc_name:=${protoc_name}-arm64
		protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}-arm64
		protoc_gen_go_name:=${protoc_gen_go_name}-arm64
	endif
	ifneq ($(findstring $(machine_arch),x86_64 amd64),)
		protoc_name:=${protoc_name}-amd64
		protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}-amd64
		protoc_gen_go_name:=${protoc_gen_go_name}-amd64
	endif
else 
	$(error Unknown architectures.)
endif

# windows file suffix
ifeq ($(os_name),Windows_NT)
	protoc_name:=${protoc_name}.exe
	protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}.exe
	protoc_gen_go_name:=${protoc_gen_go_name}.exe
endif
ifneq ($(findstring MINGW64_NT,$(os_name)),)
	protoc_name:=${protoc_name}.exe
	protoc_gen_go_grpc_name:=${protoc_gen_go_grpc_name}.exe
	protoc_gen_go_name:=${protoc_gen_go_name}.exe
endif

#-------------------------------------------------------------------------------------






#----------------------------------------build/clean---------------------------------------------
.PHONY: build
build: render client server

# render
.PHONY: render
render:
	go build -mod=vendor -o ${render_cmd} ${project_dir}/cmd/render
	chmod +x ${render_cmd}

# client
.PHONY: client
client: generate_grpc_client
	go build -mod=vendor -o ${project_dir}/bin/client ${project_dir}/cmd/client
	chmod +x ${project_dir}/bin/client

# server
.PHONY: server
server: generate_grpc_server
	go build -mod=vendor -o ${project_dir}/bin/server ${project_dir}/cmd/server
	chmod +x ${project_dir}/bin/server

.PHONY: clean
clean:
	rm -rf ${render_cmd}
	rm -rf ${client_cmd}
	rm -rf ${server_cmd}

# Specific target variable
generate_grpc_client: proto_out_dir=${project_dir}/internal/client/proto
generate_grpc_client: renderval_path=${project_dir}/configs/render_value.yaml
generate_grpc_client: temp_path=${project_dir}/configs/template.go.txt

.PHONY: generate_grpc_client
generate_grpc_client:
	${project_dir}/bin/${protoc_name} \
		--plugin=protoc-gen-go-grpc=${project_dir}/bin/${protoc_gen_go_grpc_name} \
		--plugin=protoc-gen-go=${project_dir}/bin/${protoc_gen_go_name} \
		--go-grpc_out=${proto_out_dir} --proto_path=${protofile_path} ${protofile_name}
	${project_dir}/bin/${protoc_name} \
		--plugin=protoc-gen-go-grpc=${project_dir}/bin/${protoc_gen_go_grpc_name} \
		--plugin=protoc-gen-go=${project_dir}/bin/${protoc_gen_go_name} \
		--go_out=${proto_out_dir} --proto_path=${protofile_path}  ${protofile_name}    
# render the template file my_generate.go.tmpl, generate my_generate.go
	${render_cmd} --renderval_path=${renderval_path}  --temp_path=${temp_path}  > ${project_dir}/internal/client/proto/my_generate.go


generate_grpc_server: proto_out_dir=${project_dir}/internal/server/proto
.PHONY: generate_grpc_server
generate_grpc_server: 
	${project_dir}/bin/${protoc_name} \
		--plugin=protoc-gen-go-grpc=${project_dir}/bin/${protoc_gen_go_grpc_name} \
		--plugin=protoc-gen-go=${project_dir}/bin/${protoc_gen_go_name} \
		--go-grpc_out=${proto_out_dir} --proto_path=${protofile_path} ${protofile_name}  
	${project_dir}/bin/${protoc_name} \
		--plugin=protoc-gen-go-grpc=${project_dir}/bin/${protoc_gen_go_grpc_name} \
		--plugin=protoc-gen-go=${project_dir}/bin/${protoc_gen_go_name} \
		--go_out=${proto_out_dir} --proto_path=${protofile_path}  ${protofile_name}    


#----------------------------------------start programs---------------------------------------------
.PHONY: start_server
start_server: server
	${project_dir}/bin/server \
	--rootca_path=${project_dir}/configs/key/cacert.pem \
	--servercert_path=${project_dir}/configs/key/server.crt \
	--serverkey_path=${project_dir}/configs/key/server.key \
	--port=9090

.PHONY: start_client
start_client: render client
	${project_dir}/bin/client \
    --rootca_path=${project_dir}/configs/key/cacert.pem \
    --port=8080 \
    --loglevel=7 \
    --grpc_server_addr=my.grpc.com:9090 \
    --grpc_skip_verify=false


#-------------------------------------------------------------------------------------

