package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		expectedError string
		expectedPlan  *TestPlan
	}{
		{
			name:          "NotExist",
			file:          "./test/noplan.yaml",
			expectedError: "no such file or directory",
			expectedPlan:  nil,
		},
		{
			name:          "EmptyFile",
			file:          "./test/plan_empty.yaml",
			expectedError: "EOF",
			expectedPlan:  nil,
		},
		{
			name:          "InvalidFile",
			file:          "./test/plan_error.yaml",
			expectedError: "cannot unmarshal",
			expectedPlan:  nil,
		},
		{
			name:          "Simple",
			file:          "./test/plan_simple.yaml",
			expectedError: "",
			expectedPlan: &TestPlan{
				Name: "Simple",
			},
		},
		{
			name:          "Full",
			file:          "./test/plan_full.yaml",
			expectedError: "",
			expectedPlan: &TestPlan{
				Name: "Simple",
				RESTTests: []RESTTest{
					RESTTest{
						Name:           "site-service",
						ServiceAddress: "site-service:4010",
					},
					RESTTest{
						Name:           "sensor-service",
						ServiceAddress: "sensor-service:4020",
					},
				},
				GRPCTests: []GRPCTest{
					GRPCTest{
						Name:           "switch-service",
						ServiceAddress: "switch-service:4030",
					},
				},
				NATSTests: []NATSTest{
					NATSTest{
						Name:         "asset-service",
						NATSServers:  []string{"nats-0:4222", "nats-1:4222", "nats-2:4222"},
						NATSUser:     "api-test",
						NATSPassword: "password",
					},
				},
				GraphQLTests: []GraphQLTest{
					GraphQLTest{
						Name:           "graphql-server",
						ServiceAddress: "graphql-server:5000",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := Read(tc.file)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, plan)
				fmt.Printf("%+v\n", plan)
			} else {
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, plan)
			}
		})
	}
}
