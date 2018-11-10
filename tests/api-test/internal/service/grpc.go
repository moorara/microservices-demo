package service

type (
	// GRPCPlan is for testing a service with gRPC API
	GRPCPlan struct {
		Name           string     `yaml:"name"`
		Address        string     `yaml:"address"`
		CAChainFile    string     `yaml:"ca_chain_file"`
		ClientCertFile string     `yaml:"client_cert_file"`
		ClientKeyFile  string     `yaml:"client_key_file"`
		Tests          []GRPCTest `yaml:"tests"`
	}

	// GRPCTest defines test spec for a GRPC service
	GRPCTest struct {
	}
)
