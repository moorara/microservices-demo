package service

type (
	// GRPCTest is for testing a service with gRPC API
	GRPCTest struct {
		Name           string `yaml:"name"`
		ServiceAddress string `yaml:"service_address"`
		CAChainFile    string `yaml:"ca_chain_file"`
		ClientCertFile string `yaml:"client_cert_file"`
		ClientKeyFile  string `yaml:"client_key_file"`
	}
)
