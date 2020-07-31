package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/asset/internal/queue"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

const (
	natsClientName = "tester"
)

type Event struct {
	ID        string
	Payload   string
	Timestamp time.Time
}

func TestNATSConnection(t *testing.T) {
	if !Config.IntegrationTest {
		t.SkipNow()
	}

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

	conn, err := queue.NewNATSConnection(Config.NatsServers, natsClientName, Config.NatsUser, Config.NatsPassword)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eventCount := len(tc.events)
			workerCount := tc.workersCount

			clientSubscribed := make(chan bool, 1)
			clientProcessed := make(chan bool, eventCount)

			workerSubscribed := make(chan bool, workerCount)
			workerProcessed := make(chan bool, eventCount)

			// Subscribe to the subject
			go func() {
				_, err := conn.Subscribe(tc.subject, func(msg *nats.Msg) {
					var event Event
					err := json.Unmarshal(msg.Data, &event)
					assert.NoError(t, err)
					clientProcessed <- true
				})
				assert.NoError(t, err)
				clientSubscribed <- true
			}()

			// Subscribe to a working group in the subject
			for i := 0; i < tc.workersCount; i++ {
				go func() {
					_, err := conn.QueueSubscribe(tc.subject, tc.queue, func(msg *nats.Msg) {
						var event Event
						err := json.Unmarshal(msg.Data, &event)
						assert.NoError(t, err)
						workerProcessed <- true
					})
					assert.NoError(t, err)
					workerSubscribed <- true
				}()
			}

			// Publish to the subject
			go func() {
				<-clientSubscribed
				for i := 0; i < workerCount; i++ {
					<-workerSubscribed
				}

				for _, event := range tc.events {
					data, err := json.Marshal(&event)
					assert.NoError(t, err)
					err = conn.Publish(tc.subject, data)
					assert.NoError(t, err)
				}
			}()

			// Wait for all subscribers to finish
			for i := 0; i < eventCount; i++ {
				<-clientProcessed
				if tc.workersCount > 0 {
					<-workerProcessed
				}
			}
		})
	}
}
