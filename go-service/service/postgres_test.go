package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB(t *testing.T) {
	tests := []struct {
		name        string
		postgresURL string
	}{
		{
			"WithoutUserPass",
			"postgres://localhost",
		},
		{
			"WithUserPass",
			"postgres://root:pass@localhost",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := NewPostgresDB(tc.postgresURL)
			assert.NotNil(t, db)
		})
	}
}
