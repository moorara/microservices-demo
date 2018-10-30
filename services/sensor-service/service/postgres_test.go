package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     string
		database string
		username string
		password string
	}{
		{
			name:     "Simple",
			host:     "localhost",
			port:     "5432",
			database: "store",
		},
		{
			name:     "WithUsername",
			host:     "localhost",
			port:     "5432",
			database: "store",
			username: "root",
		},
		{
			name:     "WithPassword",
			host:     "localhost",
			port:     "5432",
			database: "store",
			username: "root",
			password: "pass",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := NewPostgresDB(tc.host, tc.port, tc.database, tc.username, tc.password)
			assert.NotNil(t, db)
		})
	}
}
