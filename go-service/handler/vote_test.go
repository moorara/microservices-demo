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
	"github.com/moorara/microservices-demo/go-service/service"
	"github.com/stretchr/testify/assert"
)

type mockVoteManager struct {
	CreateCalled bool
	CreateVote   *service.Vote
	CreateError  error

	GetAllCalled bool
	GetAllVotes  []service.Vote
	GetAllError  error

	GetCalled bool
	GetVote   *service.Vote
	GetError  error

	DeleteCalled bool
	DeleteError  error
}

func (vm *mockVoteManager) Create(ctx context.Context, linkID string, stars int) (*service.Vote, error) {
	vm.CreateCalled = true
	return vm.CreateVote, vm.CreateError
}

func (vm *mockVoteManager) GetAll(ctx context.Context, linkID string) ([]service.Vote, error) {
	vm.GetAllCalled = true
	return vm.GetAllVotes, vm.GetAllError
}

func (vm *mockVoteManager) Get(ctx context.Context, id string) (*service.Vote, error) {
	vm.GetCalled = true
	return vm.GetVote, vm.GetError
}

func (vm *mockVoteManager) Delete(ctx context.Context, id string) error {
	vm.DeleteCalled = true
	return vm.DeleteError
}

func TestNewVoteHandler(t *testing.T) {
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
			db := service.NewPostgresDB(tc.postgresURL)
			logger := log.NewNopLogger()
			vh := NewVoteHandler(db, logger)

			assert.NotNil(t, vh)
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
			``,
		},
		{
			"InvalidJSON",
			nil, nil,
			`{"linkId": "1111-aaaa"`,
			400,
			``,
		},
		{
			"VoteManagerError",
			nil, errors.New("error"),
			`{"linkId": "1111-aaaa", "stars": 5}`,
			500,
			``,
		},
		{
			"Successful",
			&service.Vote{
				ID:     "2222-bbbb",
				LinkID: "1111-aaaa",
				Stars:  5,
			},
			nil,
			`{"linkId": "1111-aaaa", "stars": 5}`,
			201,
			`{"id":"2222-bbbb","linkId":"1111-aaaa","stars":5}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vm := &mockVoteManager{
				CreateVote:  tc.createVote,
				CreateError: tc.createError,
			}

			vh := &postgresVoteHandler{
				vm:     vm,
				logger: log.NewNopLogger(),
			}

			reqBody := strings.NewReader(tc.reqBody)
			r := httptest.NewRequest("POST", "http://service/votes", reqBody)
			w := httptest.NewRecorder()

			vh.PostVote(w, r)
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

func TestGetVotes(t *testing.T) {
	tests := []struct {
		name            string
		getAllVotes     []service.Vote
		getAllError     error
		linkID          string
		expectedStatus  int
		expectedResBody string
	}{
		{
			"NoLinkID",
			nil, nil,
			"",
			404,
			``,
		},
		{
			"VoteManagerError",
			nil, errors.New("error"),
			"1111-aaaa",
			500,
			``,
		},
		{
			"Successful",
			[]service.Vote{
				{ID: "2222-bbbb", LinkID: "1111-aaaa", Stars: 3},
				{ID: "4444-dddd", LinkID: "1111-aaaa", Stars: 4},
			},
			nil,
			"1111-aaaa",
			200,
			`[{"id":"2222-bbbb","linkId":"1111-aaaa","stars":3},{"id":"4444-dddd","linkId":"1111-aaaa","stars":4}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vm := &mockVoteManager{
				GetAllVotes: tc.getAllVotes,
				GetAllError: tc.getAllError,
			}

			vh := &postgresVoteHandler{
				vm:     vm,
				logger: log.NewNopLogger(),
			}

			mr := mux.NewRouter()
			mr.Path("/votes").Queries("linkId", "{linkId:[0-9a-f-]+}").Methods("GET").HandlerFunc(vh.GetVotes)
			ts := httptest.NewServer(mr)
			defer ts.Close()

			res, err := http.Get(ts.URL + "/votes?linkId=" + tc.linkID)
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
			``,
		},
		{
			"VoteManagerError",
			nil, errors.New("error"),
			"2222-bbbb",
			500,
			``,
		},
		{
			"NoVoteFound",
			nil, nil,
			"2222-bbbb",
			404,
			``,
		},
		{
			"Successful",
			&service.Vote{
				ID:     "2222-bbbb",
				LinkID: "1111-aaaa",
				Stars:  5,
			},
			nil,
			"1111-aaaa",
			200,
			`{"id":"2222-bbbb","linkId":"1111-aaaa","stars":5}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vm := &mockVoteManager{
				GetVote:  tc.getVote,
				GetError: tc.getError,
			}

			vh := &postgresVoteHandler{
				vm:     vm,
				logger: log.NewNopLogger(),
			}

			mr := mux.NewRouter()
			mr.Path("/votes/{id:[0-9a-f-]+}").Methods("GET").HandlerFunc(vh.GetVote)
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

func TestDeleteVote(t *testing.T) {
	tests := []struct {
		name            string
		deleteError     error
		voteID          string
		expectedStatus  int
		expectedResBody string
	}{
		{
			"NoVoteID",
			nil,
			"",
			404,
			``,
		},
		{
			"VoteManagerError",
			errors.New("error"),
			"2222-bbbb",
			500,
			``,
		},
		{
			"Successful",
			nil,
			"2222-bbbb",
			204,
			``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vm := &mockVoteManager{
				DeleteError: tc.deleteError,
			}

			vh := &postgresVoteHandler{
				vm:     vm,
				logger: log.NewNopLogger(),
			}

			mr := mux.NewRouter()
			mr.Path("/votes/{id:[0-9a-f-]+}").Methods("DELETE").HandlerFunc(vh.DeleteVote)
			ts := httptest.NewServer(mr)
			defer ts.Close()

			req, err := http.NewRequest("DELETE", ts.URL+"/votes/"+tc.voteID, nil)
			assert.NoError(t, err)
			client := &http.Client{}
			res, err := client.Do(req)
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
