package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		name                         string
		endpoints                    []string
		user, password               string
		databaseName, collectionName string
		expectError                  bool
	}{
		{
			"ConnectionError",
			[]string{"http://localhost:9999"},
			"root", "pass",
			"resources", "things",
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			arango, err := NewArangoService(ctx, tc.endpoints, tc.user, tc.password, tc.databaseName, tc.collectionName)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, arango)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, arango)
			}
		})
	}
}
