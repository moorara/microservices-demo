package component

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/asset-service/internal/queue"
	"github.com/stretchr/testify/assert"
)

const (
	natsClientName = "tester"
)

func TestAPI(t *testing.T) {
	if !Config.ComponentTest {
		t.SkipNow()
	}

	tests := []struct {
		name             string
		subject          string
		request          string
		expectedResponse string
	}{
		{
			name:             "Test",
			subject:          "asset_service",
			request:          "Hello, World!",
			expectedResponse: "Your message is received!",
		},
	}

	nats, err := queue.NewNATSConnection(Config.NatsServers, natsClientName, Config.NatsUser, Config.NatsPassword)
	assert.NoError(t, err)
	assert.NotNil(t, nats)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			msg, err := nats.RequestWithContext(ctx, tc.subject, []byte(tc.request))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, string(msg.Data))

			fmt.Println(string(msg.Data))
		})
	}
}
