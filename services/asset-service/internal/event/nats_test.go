package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNATSConnection(t *testing.T) {
	tests := []struct {
		name           string
		servers        []string
		clientName     string
		user, password string
	}{
		{
			"NoServer",
			[]string{},
			"",
			"", "",
		},
		{
			"NoAuth",
			[]string{"localhost:4222"},
			"",
			"", "",
		},
		{
			"WithName",
			[]string{"localhost:4222"},
			"service-name",
			"", "",
		},
		{
			"WithAuth",
			[]string{"nats1:4222"},
			"service-name",
			"nats_client", "passsword",
		},
		{
			"Cluster",
			[]string{"nats1:4222", "nats2:4222", "nats3:4222"},
			"service-name",
			"nats_client", "passsword",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nc, err := NewNATSConnection(tc.servers, tc.clientName, tc.user, tc.password)

			assert.Error(t, err)
			assert.Nil(t, nc)
		})
	}
}
