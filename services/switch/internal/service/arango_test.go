package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockResp struct {
	StatusCode int
	Body       string
}

func TestArangoService(t *testing.T) {
	tests := []struct {
		name                 string
		user, password       string
		database, collection string
		contextTimeout       time.Duration
		mocks                map[string]mockResp
		expectError          bool
	}{
		{
			"Unauthorized",
			"user", "pass",
			"animals", "mammals",
			50 * time.Millisecond,
			map[string]mockResp{
				"/_db/animals/_api/database/current":   mockResp{http.StatusUnauthorized, `{"code":401, "error":true, "errorMessage":"not authorized"}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusUnauthorized, `{"code":401, "error":true, "errorMessage":"not authorized"}`},
			},
			true,
		},
		{
			"Unauthorized",
			"user", "pass",
			"animals", "mammals",
			50 * time.Millisecond,
			map[string]mockResp{
				"/_db/animals/_api/database/current":   mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusUnauthorized, `{"code":401, "error":true, "errorMessage":"not authorized"}`},
			},
			true,
		},
		{
			"DatabaseAndCollectionExist",
			"user", "pass",
			"animals", "mammals",
			50 * time.Millisecond,
			map[string]mockResp{
				"/_db/animals/_api/database/current":   mockResp{http.StatusOK, `{"code":200, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusOK, `{"code":200, "error":false}`},
			},
			false,
		},
		{
			"CreateDatabaseError",
			"user", "pass",
			"animals", "mammals",
			50 * time.Millisecond,
			map[string]mockResp{
				"/_db/animals/_api/database/current": mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"database not found"}`},
				"/_db/_system/_api/database":         mockResp{http.StatusBadRequest, `{"code":400, "error":false}`},
			},
			true,
		},
		{
			"CreateCollectionError",
			"user", "pass",
			"animals", "mammals",
			50 * time.Millisecond,
			map[string]mockResp{
				"/_db/animals/_api/database/current":   mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"database not found"}`},
				"/_db/_system/_api/database":           mockResp{http.StatusCreated, `{"code":201, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"collection not found"}`},
				"/_db/animals/_api/collection":         mockResp{http.StatusBadRequest, `{"code":400, "error":false}`},
			},
			true,
		},
		{
			"CreateDatabaseAndCollection",
			"user", "pass",
			"animals", "mammals",
			50 * time.Millisecond,
			map[string]mockResp{
				"/_db/animals/_api/database/current":   mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"database not found"}`},
				"/_db/_system/_api/database":           mockResp{http.StatusCreated, `{"code":201, "error":false}`},
				"/_db/animals/_api/collection/mammals": mockResp{http.StatusNotFound, `{"code":404, "error":true, "errorMessage":"collection not found"}`},
				"/_db/animals/_api/collection":         mockResp{http.StatusOK, `{"code":200, "error":false}`},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock arango back-end
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.mocks[path].StatusCode)
				w.Write([]byte(tc.mocks[path].Body))
			}))
			defer ts.Close()

			ctx, cancel := context.WithTimeout(context.Background(), tc.contextTimeout)
			defer cancel()

			service, err := NewArangoService([]string{ts.URL}, tc.user, tc.password)
			assert.NoError(t, err)

			err = service.Connect(ctx, tc.database, tc.collection)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
