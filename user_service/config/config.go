package config

type Config struct {
	CassandraDetails struct {
		Address  string `yaml:"address"`
		KeySpace string `yaml:"key_space"`
		Port     int    `yaml:"port"`
	} `yaml:"cassandra_details"`

	GrpcDetails struct {
		Address  string `yaml:"address"`
		Network  string `yaml:"network"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"grpc_details"`

	HttpDetails struct {
		Port string `yaml:"port"`
	} `yaml:"http_details"`
}
