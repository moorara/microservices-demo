package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/asset-service/internal/queue"
	"github.com/stretchr/testify/assert"

	nats "github.com/nats-io/go-nats"
)

const (
	NatsClientName = "tester"
)

type Event struct {
	ID        string
	Payload   string
	Timestamp time.Time
}

func TestNATSConnection(t *testing.T) {
	tests := []struct {
		name         string
		subject      string
		queue        string
		workersCount int
		events       []Event
	}{
		{
			"NoWorker",
			"events",
			"",
			0,
			[]Event{
				Event{
					ID:        "1111-1111",
					Payload:   "Hakuna Matata",
					Timestamp: time.Now(),
				},
				Event{
					ID:        "2222-2222",
					Payload:   "Chill out",
					Timestamp: time.Now(),
				},
			},
		},
		{
			"MultipleWorkers",
			"events",
			"workers",
			4,
			[]Event{
				Event{
					ID:        "aaaa-aaaa",
					Payload:   "Je t'aime",
					Timestamp: time.Now(),
				},
				Event{
					ID:        "bbbb-bbbb",
					Payload:   "Te amo",
					Timestamp: time.Now(),
				},
			},
		},
	}

	conn, err := queue.NewNATSConnection(Config.NatsServers, NatsClientName, Config.NatsUser, Config.NatsPassword)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			count := len(tc.events)
			subscriberDone := make(chan bool, count)
			workerDone := make(chan bool, count)

			// Subscribe to the subject
			go func() {
				_, err := conn.Subscribe(tc.subject, func(msg *nats.Msg) {
					var event Event
					err := json.Unmarshal(msg.Data, &event)
					assert.NoError(t, err)
					subscriberDone <- true
				})
				assert.NoError(t, err)
			}()

			// Subscribe to a working group in the subject
			for i := 0; i < tc.workersCount; i++ {
				go func() {
					_, err := conn.QueueSubscribe(tc.subject, tc.queue, func(msg *nats.Msg) {
						var event Event
						err := json.Unmarshal(msg.Data, &event)
						assert.NoError(t, err)
						workerDone <- true
					})
					assert.NoError(t, err)
				}()
			}

			// Publish to the subject
			go func() {
				for _, event := range tc.events {
					data, err := json.Marshal(&event)
					assert.NoError(t, err)
					err = conn.Publish(tc.subject, data)
					assert.NoError(t, err)
				}
			}()

			err := conn.Flush()
			assert.NoError(t, err)

			for i := 0; i < count; i++ {
				<-subscriberDone
				if tc.workersCount > 0 {
					<-workerDone
				}
			}
		})
	}
}
