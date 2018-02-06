package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	CloseCalled bool
	CloseError  error

	ExecContextCalled bool
	ExecContextResult sql.Result
	ExecContextError  error

	QueryContextCalled bool
	QueryContextRows   *sql.Rows
	QueryContextError  error

	QueryRowContextCalled bool
	QueryRowContextRow    *sql.Row
}

func (mdb *mockDB) Close() error {
	mdb.CloseCalled = true
	return mdb.CloseError
}

func (mdb *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	mdb.ExecContextCalled = true
	return mdb.ExecContextResult, mdb.ExecContextError
}

func (mdb *mockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	mdb.QueryContextCalled = true
	return mdb.QueryContextRows, mdb.QueryContextError
}

func (mdb *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	mdb.QueryRowContextCalled = true
	return mdb.QueryRowContextRow
}

func TestNewVoteManager(t *testing.T) {
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
			logger := log.NewNopLogger()
			vm := NewVoteManager(db, logger)

			assert.NotNil(t, vm)
		})
	}
}

func TestVoteManagerCreate(t *testing.T) {
	tests := []struct {
		name             string
		execContextError error
		context          context.Context
		linkID           string
		stars            int
		expectError      bool
	}{
		{
			"DatabaseFailed",
			errors.New("error"),
			context.Background(),
			"",
			0,
			true,
		},
		{
			"DatabaseSuccessful",
			nil,
			context.Background(),
			"2222-bbbb",
			5,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mdb := &mockDB{
				ExecContextResult: nil,
				ExecContextError:  tc.execContextError,
			}
			vm := &postgresVoteManager{
				db:     mdb,
				logger: log.NewNopLogger(),
			}

			vote, err := vm.Create(tc.context, tc.linkID, tc.stars)

			assert.True(t, mdb.ExecContextCalled)
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

func TestVoteManagerDelete(t *testing.T) {
	tests := []struct {
		name               string
		execContextError   error
		context            context.Context
		voteID             string
		expectError        bool
		expectedVoteID     string
		expectedVoteLinkID string
		expectedVoteStars  int
	}{
		{
			"DatabaseFailed",
			errors.New("error"),
			context.Background(),
			"",
			true,
			"", "", 0,
		},
		{
			"DatabaseSuccessful",
			nil,
			context.Background(),
			"2222-bbbb",
			false,
			"1111-aaaa", "2222-bbbb", 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mdb := &mockDB{
				ExecContextResult: nil,
				ExecContextError:  tc.execContextError,
			}
			vm := &postgresVoteManager{
				db:     mdb,
				logger: log.NewNopLogger(),
			}

			err := vm.Delete(tc.context, tc.voteID)

			assert.True(t, mdb.ExecContextCalled)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
