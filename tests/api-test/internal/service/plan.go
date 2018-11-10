package service

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type (
	// TestPlan defines all tests to be run
	TestPlan struct {
		Name         string        `yaml:"name"`
		RESTPlans    []RESTPlan    `yaml:"rest"`
		GRPCPlans    []GRPCPlan    `yaml:"grpc"`
		NATSPlans    []NATSPlan    `yaml:"nats"`
		GraphQLPlans []GraphQLPlan `yaml:"graphql"`
	}
)

// Read returns a test plan for an input file
func Read(file string) (*TestPlan, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	plan := new(TestPlan)
	err = yaml.NewDecoder(f).Decode(plan)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
