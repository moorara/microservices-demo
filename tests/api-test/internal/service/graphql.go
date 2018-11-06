package service

type (
	// GraphQLTest is for testing a service with GraphQL API
	GraphQLTest struct {
		Name           string `yaml:"name"`
		ServiceAddress string `yaml:"service_address"`
	}
)
