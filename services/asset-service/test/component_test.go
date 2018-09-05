package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func componentTest() bool {
	value := os.Getenv("COMPONENT_TEST")
	return value == "true" || value == "TRUE"
}

func getServiceURL() string {
	serviceURL := os.Getenv("SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "http://localhost:4040"
	}

	return serviceURL
}

func makeQuery(client *http.Client, query string) (map[string]interface{}, error) {
	serviceURL := getServiceURL()
	endpoint := serviceURL + "/graphql"

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(map[string]string{
		"query": query,
	})

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func TestLiveness(t *testing.T) {
	if !componentTest() {
		t.SkipNow()
	}

	serviceURL := getServiceURL()

	for i := 0; i < 30; i++ {
		resp, err := http.Get(serviceURL + "/liveness")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(time.Second)
	}

	assert.Fail(t, "timeout")
}

func TestReadiness(t *testing.T) {
	if !componentTest() {
		t.SkipNow()
	}

	serviceURL := getServiceURL()

	for i := 0; i < 30; i++ {
		resp, err := http.Get(serviceURL + "/readiness")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(time.Second)
	}

	assert.Fail(t, "timeout")
}

func TestMetrics(t *testing.T) {
	if !componentTest() {
		t.SkipNow()
	}

	serviceURL := getServiceURL()
	url := serviceURL + "/metrics"
	res, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	body := string(data)

	assert.Contains(t, body, "# TYPE go_goroutines gauge")
	assert.Contains(t, body, "# TYPE go_gc_duration_seconds summary")
	assert.Contains(t, body, "# TYPE go_memstats_heap_objects gauge")
	assert.Contains(t, body, "# TYPE go_memstats_alloc_bytes gauge")
	assert.Contains(t, body, "# TYPE go_memstats_alloc_bytes_total counter")
	assert.Contains(t, body, "# TYPE asset_service_process_open_fds gauge")
	assert.Contains(t, body, "# TYPE asset_service_process_resident_memory_bytes gauge")
	assert.Contains(t, body, "# TYPE asset_service_process_virtual_memory_bytes gauge")
	assert.Contains(t, body, "# TYPE asset_service_jaeger_traces counter")
	assert.Contains(t, body, "# TYPE asset_service_jaeger_started_spans counter")
	assert.Contains(t, body, "# TYPE asset_service_jaeger_finished_spans counter")
}

func TestAPI(t *testing.T) {
	if !componentTest() {
		t.SkipNow()
	}

	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   2 * time.Second,
	}

	tests := []struct {
		name             string
		query            string
		expectedResponse map[string]interface{}
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			response, _ := makeQuery(client, tc.query)

			responseData, err := json.Marshal(response)
			assert.NoError(t, err)

			expectedData, err := json.Marshal(tc.expectedResponse)
			assert.NoError(t, err)

			assert.Equal(t, string(expectedData), string(responseData))
		})
	}
}
