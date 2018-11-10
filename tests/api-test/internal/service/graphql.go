package service

type (
	// GraphQLPlan is for testing a service with GraphQL API
	GraphQLPlan struct {
		Name    string        `yaml:"name"`
		Address string        `yaml:"address"`
		Tests   []GraphQLTest `yaml:"tests"`
	}

	// GraphQLTest defines test spec for a GraphQL service
	GraphQLTest struct {
	}
)
