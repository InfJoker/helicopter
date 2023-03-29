package config

type Config struct {
	GrpcServer struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
	} `yaml:"grpc_server"`
	HttpServer struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
	} `yaml:"http_server"`
	LseqDb struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
	} `yaml:"lseqdb"`
}
