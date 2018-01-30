package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisPersisterSave(t *testing.T) {
	tests := []struct {
		name     string
		redisURL string
		key      string
		data     []byte
		ttl      time.Duration
	}{
		{
			"WithoutUserPass",
			"redis://redis:6379",
			"1234",
			[]byte("data"),
			10 * time.Minute,
		},
		{
			"WithUserPass",
			"redis://user:pass@redis:6389",
			"abcd",
			[]byte("content"),
			1 * time.Hour,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rp := NewRedisPersister(tc.redisURL)
			assert.NotNil(t, rp)

			err := rp.Save(tc.key, tc.data, tc.ttl)
			assert.Error(t, err)
		})
	}
}

func TestRedisPersisterLoad(t *testing.T) {
	tests := []struct {
		name     string
		redisURL string
		key      string
	}{
		{
			"WithoutUserPass",
			"redis://redis:6379",
			"1234",
		},
		{
			"WithUserPass",
			"redis://user:pass@redis:6389",
			"abcd",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rp := NewRedisPersister(tc.redisURL)
			assert.NotNil(t, rp)

			data, err := rp.Load(tc.key)
			assert.Error(t, err)
			assert.Nil(t, data)
		})
	}
}
