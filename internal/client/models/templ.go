package models

// yaml to go struct, https://zhwt.github.io/yaml-to-go/

type Templ struct {
	PackageName string    `yaml:"packageName" json:"packageName"`
	Service     []Service `yaml:"service" json:"service"`
}

type Service struct {
	ServiceName string `yaml:"serviceName" json:"serviceName"`
	RPC         []RPC  `yaml:"rpc" json:"rpc"`
}

type RPC struct {
	RPCFuncName string   `yaml:"rpcFuncName" json:"rpcFuncName"`
	Request     Request  `yaml:"request" json:"request"`
	Response    Response `yaml:"response" json:"response"`
}

type Request struct {
	Type    string `yaml:"type" json:"type"`
	Message string `yaml:"message" json:"message"`
}
type Response struct {
	Type    string `yaml:"type" json:"type"`
	Message string `yaml:"message"`
}
