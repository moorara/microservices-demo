package handler

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/toys/microservices/go-service/service"
	"github.com/stretchr/testify/assert"
)

type mockVoteManager struct {
	createCalled bool
	createVote   *service.Vote
	createError  error

	getCalled bool
	getVote   *service.Vote
	getError  error
}

func (sm *mockVoteManager) Create(ctx context.Context, linkID string, stars int) (*service.Vote, error) {
	sm.createCalled = true
	return sm.createVote, sm.createError
}

func (sm *mockVoteManager) Get(ctx context.Context, id string) (*service.Vote, error) {
	sm.getCalled = true
	return sm.getVote, sm.getError
}

func TestNewVoteHandler(t *testing.T) {
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
			rp := service.NewRedisPersister(tc.redisURL)
			logger := log.NewNopLogger()
			sh := NewVoteHandler(rp, logger)

			assert.NotNil(t, sh)
		})
	}
}

func TestPostVote(t *testing.T) {
	tests := []struct {
		name            string
		createVote      *service.Vote
		createError     error
		reqBody         string
		expectedStatus  int
		expectedResBody string
	}{
		{
			"InvalidRequest",
			nil, nil,
			`{}`,
			400,
			"",
		},
		{
			"InvalidJSON",
			nil, nil,
			`{"linkId": "2222-bbbb"`,
			400,
			"",
		},
		{
			"VoteManagerError",
			nil, errors.New("error"),
			`{"linkId": "2222-bbbb", "stars": 5}`,
			500,
			"",
		},
		{
			"Successful",
			&service.Vote{
				ID:     "1111-aaaa",
				LinkID: "2222-bbbb",
				Stars:  5,
			},
			nil,
			`{"linkId": "2222-bbbb", "stars": 5}`,
			201,
			`{"id":"1111-aaaa","linkId":"2222-bbbb","stars":5}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sm := &mockVoteManager{
				createVote:  tc.createVote,
				createError: tc.createError,
			}

			sh := &redisVoteHandler{
				sm:     sm,
				logger: log.NewNopLogger(),
			}

			reqBody := strings.NewReader(tc.reqBody)
			r := httptest.NewRequest("POST", "http://service/votes", reqBody)
			w := httptest.NewRecorder()

			sh.PostVote(w, r)
			res := w.Result()
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			if tc.expectedStatus == http.StatusCreated {
				assert.Contains(t, string(body), tc.expectedResBody)
			}
		})
	}
}

func TestGetVote(t *testing.T) {
	tests := []struct {
		name            string
		getVote         *service.Vote
		getError        error
		voteID          string
		expectedStatus  int
		expectedResBody string
	}{
		{
			"NoVoteID",
			nil, nil,
			"",
			404,
			`{}`,
		},
		{
			"VoteManagerError",
			nil, errors.New("error"),
			"22bb",
			404,
			`{}`,
		},
		{
			"Successful",
			&service.Vote{
				ID:     "1111-aaaa",
				LinkID: "2222-bbbb",
				Stars:  5,
			},
			nil,
			"44dd",
			200,
			`{"id":"1111-aaaa","linkId":"2222-bbbb","stars":5}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sm := &mockVoteManager{
				getVote:  tc.getVote,
				getError: tc.getError,
			}

			sh := &redisVoteHandler{
				sm:     sm,
				logger: log.NewNopLogger(),
			}

			mr := mux.NewRouter()
			mr.HandleFunc("/votes/{id:[0-9a-f]+}", sh.GetVote)
			ts := httptest.NewServer(mr)
			defer ts.Close()

			res, err := http.Get(ts.URL + "/votes/" + tc.voteID)
			assert.NoError(t, err)
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			if tc.expectedStatus == http.StatusOK {
				assert.Contains(t, string(body), tc.expectedResBody)
			}
		})
	}
}
