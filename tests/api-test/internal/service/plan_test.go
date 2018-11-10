package service

import (
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
				Name:         "Simple",
				RESTPlans:    []RESTPlan{},
				GRPCPlans:    []GRPCPlan{},
				NATSPlans:    []NATSPlan{},
				GraphQLPlans: []GraphQLPlan{},
			},
		},
		{
			name:          "Full",
			file:          "./test/plan_full.yaml",
			expectedError: "",
			expectedPlan: &TestPlan{
				Name: "Full",
				RESTPlans: []RESTPlan{
					RESTPlan{
						Name:    "site-service",
						Address: "site-service:4010",
						Base:    "/v1",
						Headers: map[string]string{
							"Is-Test": "True",
						},
						Tests: []RESTTest{
							RESTTest{
								Name:     "CheckHealth",
								Method:   "GET",
								Endpoint: "/health",
								Headers: map[string]string{
									"Test-Name": "Check-Health",
								},
								Expect: RESTExpect{
									StatusCode: 200,
								},
							},
							RESTTest{
								Name:     "GetSites",
								Method:   "GET",
								Endpoint: "/sites",
								Headers: map[string]string{
									"Test-Name": "Get-Sites",
								},
								Expect: RESTExpect{
									StatusCode: 200,
								},
							},
						},
					},
				},
				GRPCPlans: []GRPCPlan{
					GRPCPlan{
						Name:    "switch-service",
						Address: "switch-service:4030",
					},
				},
				NATSPlans: []NATSPlan{
					NATSPlan{
						Name:         "asset-service",
						NATSServers:  []string{"nats-0:4222", "nats-1:4222", "nats-2:4222"},
						NATSUser:     "api-test",
						NATSPassword: "password",
					},
				},
				GraphQLPlans: []GraphQLPlan{
					GraphQLPlan{
						Name:    "graphql-server",
						Address: "graphql-server:5000",
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
			} else {
				assert.Contains(t, err.Error(), tc.expectedError)
			}

			assert.Equal(t, tc.expectedPlan, plan)
		})
	}
}
