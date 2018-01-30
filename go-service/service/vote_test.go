package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	delay   = 200 * time.Millisecond
	timeout = 100 * time.Millisecond
)

type mockPersister struct {
	saveCalled bool
	saveError  error

	loadCalled bool
	loadData   []byte
	loadError  error
}

func (mp *mockPersister) Save(key string, data []byte, ttl time.Duration) error {
	time.Sleep(delay)
	mp.saveCalled = true
	return mp.saveError
}

func (mp *mockPersister) Load(key string) ([]byte, error) {
	time.Sleep(delay)
	mp.loadCalled = true
	return mp.loadData, mp.loadError
}

func TestNewVoteManager(t *testing.T) {
	tests := []struct {
		name     string
		redisURL string
	}{
		{
			"WithoutUserPass",
			"redis://redis:6379",
		},
		{
			"WithUserPass",
			"redis://user:pass@redis:6389",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rp := NewRedisPersister(tc.redisURL)
			sm := NewVoteManager(rp)
			assert.NotNil(t, sm)
		})
	}
}

func TestVoteManagerCreate(t *testing.T) {
	tests := []struct {
		name        string
		saveError   error
		context     func() context.Context
		linkID      string
		stars       int
		expectError bool
	}{
		{
			"PersisterError",
			errors.New("error"),
			func() context.Context {
				return context.Background()
			},
			"",
			0,
			true,
		},
		{
			"ContextTimeout",
			nil,
			func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), timeout)
				return ctx
			},
			"",
			0,
			true,
		},
		{
			"Successful",
			nil,
			func() context.Context {
				return context.Background()
			},
			"2222-bbbb",
			5,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sm := &redisVoteManager{
				persister: &mockPersister{
					saveError: tc.saveError,
				},
				voteTTL: 1 * time.Minute,
			}

			vote, err := sm.Create(tc.context(), tc.linkID, tc.stars)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.linkID, vote.LinkID)
				assert.Equal(t, tc.stars, vote.Stars)
			}
		})
	}
}

func TestVoteManagerGet(t *testing.T) {
	tests := []struct {
		name               string
		loadData           []byte
		loadError          error
		context            func() context.Context
		voteID             string
		expectError        bool
		expectedVoteID     string
		expectedVoteLinkID string
		expectedVoteStars  int
	}{
		{
			"PersisterError",
			nil, errors.New("error"),
			func() context.Context {
				return context.Background()
			},
			"",
			true,
			"", "", 0,
		},
		{
			"ContextTimeout",
			nil, nil,
			func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), timeout)
				return ctx
			},
			"",
			true,
			"", "", 0,
		},
		{
			"Successful",
			[]byte(`{"id": "1111-aaaa", "linkId": "2222-bbbb", "stars": 5}`), nil,
			func() context.Context {
				return context.Background()
			},
			"2222-bbbb",
			false,
			"1111-aaaa", "2222-bbbb", 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sm := &redisVoteManager{
				persister: &mockPersister{
					loadData:  tc.loadData,
					loadError: tc.loadError,
				},
				voteTTL: 1 * time.Minute,
			}

			vote, err := sm.Get(tc.context(), tc.voteID)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedVoteID, vote.ID)
				assert.Equal(t, tc.expectedVoteLinkID, vote.LinkID)
				assert.Equal(t, tc.expectedVoteStars, vote.Stars)
			}
		})
	}
}
