package component

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/asset/internal/queue"
	"github.com/stretchr/testify/assert"
)

const (
	natsClientName = "tester"
)

func assertResponse(t *testing.T, expected, actual map[string]interface{}) {
	for key, expectedVal := range expected {
		expectedJSON, ok := expectedVal.(map[string]interface{})
		if ok {
			actualJSON, ok := actual[key].(map[string]interface{})
			assert.True(t, ok)
			assertResponse(t, expectedJSON, actualJSON)
		} else {
			actualVal := actual[key]
			assert.Equal(t, expectedVal, actualVal)
		}
	}
}

func TestAPI(t *testing.T) {
	if !Config.ComponentTest {
		t.SkipNow()
	}

	tests := []struct {
		name             string
		subject          string
		request          map[string]interface{}
		expectedResponse map[string]interface{}
	}{
		{
			"CreateAlarm",
			"asset_service",
			map[string]interface{}{
				"kind": "createAlarm",
				"input": map[string]interface{}{
					"siteId":   "1111-1111",
					"serialNo": "1001",
					"material": "co",
				},
			},
			map[string]interface{}{
				"kind": "createAlarm",
				"alarm": map[string]interface{}{
					"siteId":   "1111-1111",
					"serialNo": "1001",
					"material": "co",
				},
			},
		},
		{
			"AllAlarm",
			"asset_service",
			map[string]interface{}{
				"kind":   "allAlarm",
				"siteId": "0000-0000",
			},
			map[string]interface{}{
				"kind":   "allAlarm",
				"alarms": []interface{}{},
			},
		},
		{
			"GetAlarm",
			"asset_service",
			map[string]interface{}{
				"kind": "getAlarm",
				"id":   "aaaa-aaaa",
			},
			map[string]interface{}{
				"kind":  "getAlarm",
				"alarm": nil,
			},
		},
		{
			"UpdateAlarm",
			"asset_service",
			map[string]interface{}{
				"kind": "updateAlarm",
				"id":   "aaaa-aaaa",
				"input": map[string]interface{}{
					"siteId":   "1111-1111",
					"serialNo": "1001",
					"material": "smoke",
				},
			},
			map[string]interface{}{
				"kind":    "updateAlarm",
				"updated": false,
			},
		},
		{
			"DeleteAlarm",
			"asset_service",
			map[string]interface{}{
				"kind": "deleteAlarm",
				"id":   "aaaa-aaaa",
			},
			map[string]interface{}{
				"kind":    "deleteAlarm",
				"deleted": false,
			},
		},
		{
			"CreateCamera",
			"asset_service",
			map[string]interface{}{
				"kind": "createCamera",
				"input": map[string]interface{}{
					"siteId":     "1111-1111",
					"serialNo":   "2001",
					"resolution": 921600,
				},
			},
			map[string]interface{}{
				"kind": "createCamera",
				"camera": map[string]interface{}{
					"siteId":     "1111-1111",
					"serialNo":   "2001",
					"resolution": float64(921600),
				},
			},
		},
		{
			"AllCamera",
			"asset_service",
			map[string]interface{}{
				"kind":   "allCamera",
				"siteId": "0000-0000",
			},
			map[string]interface{}{
				"kind":    "allCamera",
				"cameras": []interface{}{},
			},
		},
		{
			"GetCamera",
			"asset_service",
			map[string]interface{}{
				"kind": "getCamera",
				"id":   "bbbb-bbbb",
			},
			map[string]interface{}{
				"kind":   "getCamera",
				"camera": nil,
			},
		},
		{
			"UpdateCamera",
			"asset_service",
			map[string]interface{}{
				"kind": "updateCamera",
				"id":   "bbbb-bbbb",
				"input": map[string]interface{}{
					"siteId":     "1111-1111",
					"serialNo":   "2001",
					"resolution": 2073600,
				},
			},
			map[string]interface{}{
				"kind":    "updateCamera",
				"updated": false,
			},
		},
		{
			"DeleteCamera",
			"asset_service",
			map[string]interface{}{
				"kind": "deleteCamera",
				"id":   "bbbb-bbbb",
			},
			map[string]interface{}{
				"kind":    "deleteCamera",
				"deleted": false,
			},
		},
	}

	nats, err := queue.NewNATSConnection(Config.NatsServers, natsClientName, Config.NatsUser, Config.NatsPassword)
	assert.NoError(t, err)
	assert.NotNil(t, nats)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			request, err := json.Marshal(tc.request)
			assert.NoError(t, err)
			assert.NotNil(t, request)

			msg, err := nats.RequestWithContext(ctx, tc.subject, request)
			assert.NoError(t, err)

			var response map[string]interface{}
			err = json.Unmarshal(msg.Data, &response)
			assert.NoError(t, err)

			assertResponse(t, tc.expectedResponse, response)
		})
	}
}
