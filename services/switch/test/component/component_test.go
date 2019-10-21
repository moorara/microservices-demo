package test

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/switch/internal/client"
	"github.com/moorara/microservices-demo/services/switch/internal/proto"
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
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
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
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
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
		getSwitchesResponses map[string][]proto.Switch

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
				proto.Switch{Id: "", SiteId: "1111", Name: "Light", State: "OFF", States: []string{"OFF", "ON"}},
				proto.Switch{Id: "", SiteId: "2222", Name: "Pressure", State: "Low", States: []string{"High", "Low"}},
				proto.Switch{Id: "", SiteId: "2222", Name: "Temperature", State: "Medium", States: []string{"High", "Medium", "Low"}},
			},

			[]proto.GetSwitchesRequest{
				proto.GetSwitchesRequest{SiteId: "1111"},
				proto.GetSwitchesRequest{SiteId: "2222"},
			},
			map[string][]proto.Switch{
				"1111": []proto.Switch{
					proto.Switch{Id: "", SiteId: "1111", Name: "Light", State: "OFF", States: []string{"OFF", "ON"}},
				},
				"2222": []proto.Switch{
					proto.Switch{Id: "", SiteId: "2222", Name: "Pressure", State: "Low", States: []string{"High", "Low"}},
					proto.Switch{Id: "", SiteId: "2222", Name: "Temperature", State: "Medium", States: []string{"High", "Medium", "Low"}},
				},
			},

			[]proto.SetSwitchRequest{
				proto.SetSwitchRequest{Id: "", State: "ON"},
				proto.SetSwitchRequest{Id: "", State: "High"},
				proto.SetSwitchRequest{Id: "", State: "Low"},
			},
			[]proto.SetSwitchResponse{
				proto.SetSwitchResponse{},
				proto.SetSwitchResponse{},
				proto.SetSwitchResponse{},
			},

			[]proto.GetSwitchRequest{
				proto.GetSwitchRequest{Id: ""},
				proto.GetSwitchRequest{Id: ""},
				proto.GetSwitchRequest{Id: ""},
			},
			[]proto.Switch{
				proto.Switch{Id: "", SiteId: "1111", Name: "Light", State: "ON", States: []string{"OFF", "ON"}},
				proto.Switch{Id: "", SiteId: "2222", Name: "Pressure", State: "High", States: []string{"High", "Low"}},
				proto.Switch{Id: "", SiteId: "2222", Name: "Temperature", State: "Low", States: []string{"High", "Medium", "Low"}},
			},

			[]proto.RemoveSwitchRequest{
				proto.RemoveSwitchRequest{Id: ""},
				proto.RemoveSwitchRequest{Id: ""},
				proto.RemoveSwitchRequest{Id: ""},
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
					sw, err := client.InstallSwitch(ctx, &req)
					assert.NoError(t, err)

					assert.NotEmpty(t, sw.GetId())
					tc.switchID[i] = sw.GetId()
					tc.installSwitchResponses[i].Id = sw.GetId()

					assert.Equal(t, tc.installSwitchResponses[i].Id, sw.GetId())
					assert.Equal(t, tc.installSwitchResponses[i].SiteId, sw.GetSiteId())
					assert.Equal(t, tc.installSwitchResponses[i].Name, sw.GetName())
					assert.Equal(t, tc.installSwitchResponses[i].State, sw.GetState())
					assert.Equal(t, tc.installSwitchResponses[i].States, sw.GetStates())
				}
			})

			// GET SWITCHES
			t.Run("GetSwitches", func(t *testing.T) {
				for _, req := range tc.getSwitchesRequests {
					switches := []*proto.Switch{}
					stream, err := client.GetSwitches(ctx, &req)

					assert.NoError(t, err)
					assert.NotNil(t, stream)

					for {
						sw, err := stream.Recv()
						if err == io.EOF { // No more response
							break
						}
						assert.NoError(t, err)
						switches = append(switches, sw)
					}
				}
			})

			// UPDATE SWITCHES
			t.Run("SetSwitch", func(t *testing.T) {
				for i, id := range tc.switchID {
					tc.setSwitchRequests[i].Id = id
					resp, err := client.SetSwitch(ctx, &tc.setSwitchRequests[i])

					assert.NoError(t, err)
					assert.NotNil(t, resp)
				}
			})

			// GET SWITCH
			t.Run("GetSwitch", func(t *testing.T) {
				for i, id := range tc.switchID {
					tc.getSwitchRequests[i].Id = id
					tc.getSwitchResponses[i].Id = id
					sw, err := client.GetSwitch(ctx, &tc.getSwitchRequests[i])

					assert.NoError(t, err)
					assert.Equal(t, tc.getSwitchResponses[i].Id, sw.GetId())
					assert.Equal(t, tc.getSwitchResponses[i].SiteId, sw.GetSiteId())
					assert.Equal(t, tc.getSwitchResponses[i].Name, sw.GetName())
					assert.Equal(t, tc.getSwitchResponses[i].State, sw.GetState())
					assert.Equal(t, tc.getSwitchResponses[i].States, sw.GetStates())
				}
			})

			// DELETE SWITCHES
			t.Run("RemoveSwitch", func(t *testing.T) {
				for i, id := range tc.switchID {
					tc.removeSwitchRequests[i].Id = id
					resp, err := client.RemoveSwitch(ctx, &tc.removeSwitchRequests[i])

					assert.NoError(t, err)
					assert.NotNil(t, resp)
				}
			})
		})
	}
}
