package test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type JSON map[string]interface{}

func isComponentTest() bool {
	value := os.Getenv("COMPONENT_TEST")
	return value == "true" || value == "TRUE"
}

func getServiceURL() string {
	serviceURL := os.Getenv("SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "http://localhost:4010"
	}
	return serviceURL
}

func TestUnit(t *testing.T) {
	if isComponentTest() {
		t.SkipNow()
	}

	tests := []struct {
		name             string
		method, endpoint string
		reqBody          string
		statusCode       int
		resBody          string
	}{
		{
			"EmptyBody",
			"GET", "/",
			``,
			400,
			``,
		},
		{
			"SimpleGET",
			"GET", "/",
			`{}`,
			200,
			`{}`,
		},
		{
			"SimplePOST",
			"POST", "/v1/votes",
			`{"linkId":"1111-aaaa","stars":5.0}`,
			201,
			`{"id": "2222-bbbb","linkId":"1111-aaaa","stars":5.0}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rt := mux.NewRouter()
			rt.Path(tc.endpoint).Methods(tc.method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tc.resBody))
			})
			ts := httptest.NewServer(rt)
			defer ts.Close()

			os.Setenv("SERVICE_URL", ts.URL)
			cmp := NewComponent()
			statusCode, resBody, err := cmp.Call(context.Background(), tc.method, tc.endpoint, []byte(tc.reqBody))
			assert.NoError(t, err)

			assert.Equal(t, tc.statusCode, statusCode)
			assert.Equal(t, tc.resBody, string(resBody))
			os.Unsetenv("SERVICE_URL")
		})
	}
}

func TestComponentHealth(t *testing.T) {
	if !isComponentTest() {
		t.SkipNow()
	}

	endpoint := getServiceURL() + "/health"
	res, err := http.Get(endpoint)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestComponentMetrics(t *testing.T) {
	if !isComponentTest() {
		t.SkipNow()
	}

	endpoint := getServiceURL() + "/metrics"
	res, err := http.Get(endpoint)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	body := string(data)

	assert.Contains(t, body, "# TYPE go_gc_duration_seconds summary")
	assert.Contains(t, body, "# TYPE go_goroutines gauge")
	assert.Contains(t, body, "# TYPE go_memstats_alloc_bytes gauge")
	assert.Contains(t, body, "# TYPE go_memstats_frees_total counter")
	assert.Contains(t, body, "# TYPE go_memstats_heap_alloc_bytes gauge")
	assert.Contains(t, body, "# TYPE go_memstats_heap_objects gauge")
	assert.Contains(t, body, "# TYPE go_service_process_cpu_seconds_total counter")
	assert.Contains(t, body, "# TYPE go_service_process_open_fds gauge")
	assert.Contains(t, body, "# TYPE go_service_process_resident_memory_bytes gauge")
	assert.Contains(t, body, "# TYPE go_service_process_virtual_memory_bytes gauge")
}

func TestComponentSuccess(t *testing.T) {
	if !isComponentTest() {
		t.SkipNow()
	}

	tests := []struct {
		name                     string
		linkID                   string
		postVotes                []JSON
		expectedPostStatusCode   int
		expectedGetStatusCode    int
		expectedDeleteStatusCode int
	}{
		{
			"NoVote",
			"0000-0000",
			[]JSON{},
			0, 200, 0,
		},
		{
			"WithVotes",
			"1111-aaaa",
			[]JSON{
				{"linkId": "1111-aaaa"},
				{"linkId": "1111-aaaa", "stars": 4.0},
			},
			201, 200, 204,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ids := make([]string, 0)
			cmp := NewComponent()

			// CREATE VOTES
			for _, vote := range tc.postVotes {
				reqBody, err := json.Marshal(vote)
				assert.NoError(t, err)

				statusCode, resBody, err := cmp.Call(context.Background(), "POST", "/v1/votes", reqBody)
				assert.NoError(t, err)

				res := make(JSON)
				err = json.Unmarshal(resBody, &res)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedPostStatusCode, statusCode)
				assert.NotEmpty(t, res["id"])
				assert.Equal(t, vote["linkId"], res["linkId"])
				if stars, ok := vote["stars"].(float64); ok && stars > 0 {
					assert.Equal(t, stars, res["stars"])
				}
			}

			// GET VOTES
			{
				endpoint := "/v1/votes?linkId=" + tc.linkID
				statusCode, resBody, err := cmp.Call(context.Background(), "GET", endpoint, nil)
				assert.NoError(t, err)

				res := make([]JSON, 0)
				err = json.Unmarshal(resBody, &res)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedGetStatusCode, statusCode)
				for _, vote := range res {
					ids = append(ids, vote["id"].(string))
				}
			}

			// GET VOTE
			for _, id := range ids {
				endpoint := "/v1/votes/" + id
				statusCode, _, err := cmp.Call(context.Background(), "GET", endpoint, nil)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedGetStatusCode, statusCode)
			}

			// DELETE VOTE
			for _, id := range ids {
				endpoint := "/v1/votes/" + id
				statusCode, _, err := cmp.Call(context.Background(), "DELETE", endpoint, nil)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedDeleteStatusCode, statusCode)
			}
		})
	}
}

func TestComponentError(t *testing.T) {
	if !isComponentTest() {
		t.SkipNow()
	}

	tests := []struct {
		name               string
		method             string
		endpoint           string
		reqBody            JSON
		expectedStatusCode int
		expectedResBody    interface{}
	}{
		{
			"GetNotExistVote",
			"GET", "/v1/votes/0000-0000", nil,
			404, nil,
		},
		{
			"DeleteNotExistVote",
			"DELETE", "/v1/votes/0000-0000", nil,
			204, nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmp := NewComponent()

			reqBody, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			statusCode, resBody, err := cmp.Call(context.Background(), tc.method, tc.endpoint, reqBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatusCode, statusCode)
			if tc.expectedResBody != nil {
				res := make([]JSON, 0)
				err = json.Unmarshal(resBody, &res)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedResBody, res)
			}
		})
	}
}
