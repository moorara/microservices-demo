package test

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/switch-service/internal/client"
	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"github.com/stretchr/testify/assert"
)

func componentTest() bool {
	value := os.Getenv("COMPONENT_TEST")
	return value == "true" || value == "TRUE"
}

func serviceConfig() (string, string) {
	grpcAddr := os.Getenv("SERVICE_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:4030"
	}

	httpAddr := os.Getenv("SERVICE_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "http://localhost:4031"
	}

	return grpcAddr, httpAddr
}

func TestLive(t *testing.T) {
	if !componentTest() {
		t.SkipNow()
	}

	_, httpAddr := serviceConfig()

	for i := 0; i < 30; i++ {
		resp, err := http.Get(httpAddr + "/live")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		time.Sleep(time.Second)
	}

	assert.Fail(t, "timeout")
}

func TestReady(t *testing.T) {
	if !componentTest() {
		t.SkipNow()
	}

	_, httpAddr := serviceConfig()

	for i := 0; i < 30; i++ {
		resp, err := http.Get(httpAddr + "/ready")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		time.Sleep(time.Second)
	}

	assert.Fail(t, "timeout")
}

func TestContract(t *testing.T) {
	if !componentTest() {
		t.SkipNow()
	}

	grpcAddr, _ := serviceConfig()
	client, conn, err := client.New(grpcAddr)
	assert.NoError(t, err)
	defer conn.Close()

	tests := []struct {
		name     string
		switchID []string

		installSwitchRequests  []proto.InstallSwitchRequest
		installSwitchResponses []proto.Switch

		getSwitchesRequests  []proto.GetSwitchesRequest
		getSwitchesResponses [][]proto.Switch

		setSwitchRequests  []proto.SetSwitchRequest
		setSwitchResponses []proto.SetSwitchResponse

		getSwitchRequests  []proto.GetSwitchRequest
		getSwitchResponses []proto.Switch

		removeSwitchRequests  []proto.RemoveSwitchRequest
		removeSwitchResponses []proto.RemoveSwitchResponse
	}{
		{
			"Successful",
			make([]string, 3),

			[]proto.InstallSwitchRequest{
				proto.InstallSwitchRequest{SiteId: "1111", Name: "Light", State: "OFF", States: []string{"OFF", "ON"}},
				proto.InstallSwitchRequest{SiteId: "2222", Name: "Pressure", State: "Low", States: []string{"High", "Low"}},
				proto.InstallSwitchRequest{SiteId: "2222", Name: "Temperature", State: "Medium", States: []string{"High", "Medium", "Low"}},
			},
			[]proto.Switch{
				proto.Switch{Id: "TBD", SiteId: "1111", Name: "Light", State: "OFF", States: []string{"OFF", "ON"}},
				proto.Switch{Id: "TBD", SiteId: "2222", Name: "Pressure", State: "Low", States: []string{"High", "Low"}},
				proto.Switch{Id: "TBD", SiteId: "2222", Name: "Temperature", State: "Medium", States: []string{"High", "Medium", "Low"}},
			},

			[]proto.GetSwitchesRequest{
				proto.GetSwitchesRequest{SiteId: "1111"},
				proto.GetSwitchesRequest{SiteId: "2222"},
			},
			[][]proto.Switch{
				[]proto.Switch{
					proto.Switch{Id: "TBD", SiteId: "1111", Name: "Light", State: "OFF", States: []string{"OFF", "ON"}},
				},
				[]proto.Switch{
					proto.Switch{Id: "TBD", SiteId: "2222", Name: "Pressure", State: "Low", States: []string{"High", "Low"}},
					proto.Switch{Id: "TBD", SiteId: "3333", Name: "Temperature", State: "Medium", States: []string{"High", "Medium", "Low"}},
				},
			},

			[]proto.SetSwitchRequest{
				proto.SetSwitchRequest{Id: "TBD", State: "ON"},
				proto.SetSwitchRequest{Id: "TBD", State: "High"},
				proto.SetSwitchRequest{Id: "TBD", State: "Low"},
			},
			[]proto.SetSwitchResponse{
				proto.SetSwitchResponse{},
				proto.SetSwitchResponse{},
				proto.SetSwitchResponse{},
			},

			[]proto.GetSwitchRequest{
				proto.GetSwitchRequest{Id: "TBD"},
				proto.GetSwitchRequest{Id: "TBD"},
				proto.GetSwitchRequest{Id: "TBD"},
			},
			[]proto.Switch{
				proto.Switch{Id: "TBD", SiteId: "1111", Name: "Light", State: "ON", States: []string{"OFF", "ON"}},
				proto.Switch{Id: "TBD", SiteId: "2222", Name: "Pressure", State: "High", States: []string{"High", "Low"}},
				proto.Switch{Id: "TBD", SiteId: "2222", Name: "Temperature", State: "Low", States: []string{"High", "Medium", "Low"}},
			},

			[]proto.RemoveSwitchRequest{
				proto.RemoveSwitchRequest{Id: "TBD"},
				proto.RemoveSwitchRequest{Id: "TBD"},
				proto.RemoveSwitchRequest{Id: "TBD"},
			},
			[]proto.RemoveSwitchResponse{
				proto.RemoveSwitchResponse{},
				proto.RemoveSwitchResponse{},
				proto.RemoveSwitchResponse{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// CREATE SWITCHES
			t.Run("InstallSwitch", func(t *testing.T) {
				for i, req := range tc.installSwitchRequests {
					res, err := client.InstallSwitch(ctx, &req)
					tc.switchID[i] = res.GetId()

					assert.NoError(t, err)
					assert.NotEmpty(t, res.GetId())
					assert.Equal(t, tc.installSwitchResponses[i].SiteId, res.GetSiteId())
					assert.Equal(t, tc.installSwitchResponses[i].Name, res.GetName())
					assert.Equal(t, tc.installSwitchResponses[i].State, res.GetState())
					assert.Equal(t, tc.installSwitchResponses[i].States, res.GetStates())
				}
			})

			// GET SWITCHES
			t.Run("GetSwitches", func(t *testing.T) {
				for _, req := range tc.getSwitchesRequests {
					stream, err := client.GetSwitches(ctx, &req)

					assert.NoError(t, err)
					assert.NotNil(t, stream)

					for {
						res, err := stream.Recv()
						if err == io.EOF { // No more response
							break
						}

						assert.NoError(t, err)
						assert.NotNil(t, res)
						/* assert.Contains(t, tc.getSwitchesResponses[i], res) */
					}
				}
			})

			// UPDATE SWITCHES
			t.Run("SetSwitch", func(t *testing.T) {
				for i, id := range tc.switchID {
					tc.setSwitchRequests[i].Id = id
					res, err := client.SetSwitch(ctx, &tc.setSwitchRequests[i])

					assert.NoError(t, err)
					assert.NotNil(t, res)
				}
			})

			// GET SWITCHES
			t.Run("GetSwitch", func(t *testing.T) {
				for i, id := range tc.switchID {
					tc.getSwitchRequests[i].Id = id
					tc.getSwitchResponses[i].Id = id
					res, err := client.GetSwitch(ctx, &tc.getSwitchRequests[i])

					assert.NoError(t, err)
					assert.Equal(t, tc.getSwitchResponses[i].Id, res.GetId())
					/* assert.Equal(t, tc.getSwitchResponses[i].SiteId, res.GetSiteId())
					assert.Equal(t, tc.getSwitchResponses[i].Name, res.GetName())
					assert.Equal(t, tc.getSwitchResponses[i].State, res.GetState())
					assert.Equal(t, tc.getSwitchResponses[i].States, res.GetStates()) */
				}
			})

			// DELETE SWITCHES
			t.Run("RemoveSwitch", func(t *testing.T) {
				for i, id := range tc.switchID {
					tc.removeSwitchRequests[i].Id = id
					res, err := client.RemoveSwitch(ctx, &tc.removeSwitchRequests[i])

					assert.NoError(t, err)
					assert.NotNil(t, res)
				}
			})
		})
	}
}
