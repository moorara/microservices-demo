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
		serviceURL = "http://localhost:4020"
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
			"POST", "/v1/sensors",
			`{"siteId":"1111-aaaa","name":"temperature","unit":"celcius","minSafe":-30.0,"maxSafe":30.0}`,
			201,
			`{"id": "2222-bbbb","siteId":"1111-aaaa","name":"temperature","unit":"celcius","minSafe":-30.0,"maxSafe":30.0}`,
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
	assert.Contains(t, body, "# TYPE sensor_service_process_cpu_seconds_total counter")
	assert.Contains(t, body, "# TYPE sensor_service_process_open_fds gauge")
	assert.Contains(t, body, "# TYPE sensor_service_process_resident_memory_bytes gauge")
	assert.Contains(t, body, "# TYPE sensor_service_process_virtual_memory_bytes gauge")
}

func TestComponentSuccess(t *testing.T) {
	if !isComponentTest() {
		t.SkipNow()
	}

	tests := []struct {
		name                     string
		siteID                   string
		postSensors              []JSON
		putSensors               []JSON
		expectedPostStatusCode   int
		expectedAllStatusCode    int
		expectedPutStatusCode    int
		expectedGetStatusCode    int
		expectedDeleteStatusCode int
	}{
		{
			"NoSensor",
			"0000-0000",
			[]JSON{},
			[]JSON{},
			0, 200, 204, 200, 0,
		},
		{
			"WithSensors",
			"1111-aaaa",
			[]JSON{
				{"siteId": "1111-aaaa", "name": "temperature", "unit": "celsius", "minSafe": -30.0, "maxSafe": 30.0},
				{"siteId": "1111-aaaa", "name": "pressure", "unit": "atmosphere", "minSafe": 0.5, "maxSafe": 1.0},
			},
			[]JSON{
				{"siteId": "1111-aaaa", "name": "temperature", "unit": "fahrenheit", "minSafe": -22.0, "maxSafe": 86.0},
				{"siteId": "1111-aaaa", "name": "pressure", "unit": "atmosphere", "minSafe": 50000.0, "maxSafe": 100000.0},
			},
			201, 200, 204, 200, 204,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ids := make([]string, 0)
			cmp := NewComponent()

			// CREATE SENSORS
			t.Run("CreateSensors", func(t *testing.T) {
				for _, sensor := range tc.postSensors {
					reqBody, err := json.Marshal(sensor)
					assert.NoError(t, err)

					endpoint := "/v1/sensors"
					statusCode, resBody, err := cmp.Call(context.Background(), "POST", endpoint, reqBody)
					assert.NoError(t, err)

					res := make(JSON)
					err = json.Unmarshal(resBody, &res)
					assert.NoError(t, err)

					assert.Equal(t, tc.expectedPostStatusCode, statusCode)
					assert.NotEmpty(t, res["id"])
					assert.Equal(t, sensor["siteId"], res["siteId"])
					assert.Equal(t, sensor["name"], res["name"])
					assert.Equal(t, sensor["unit"], res["unit"])
					assert.Equal(t, sensor["minSafe"], res["minSafe"])
					assert.Equal(t, sensor["maxSafe"], res["maxSafe"])

					ids = append(ids, res["id"].(string))
				}
			})

			// GET SENSORS
			t.Run("GetSensors", func(t *testing.T) {
				endpoint := "/v1/sensors?siteId=" + tc.siteID
				statusCode, resBody, err := cmp.Call(context.Background(), "GET", endpoint, nil)
				assert.NoError(t, err)

				sensors := make([]JSON, 0)
				err = json.Unmarshal(resBody, &sensors)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedAllStatusCode, statusCode)
				for i, sensor := range sensors {
					assert.Equal(t, ids[i], sensor["id"])
					assert.Equal(t, tc.postSensors[i]["siteId"], sensor["siteId"])
					assert.Equal(t, tc.postSensors[i]["name"], sensor["name"])
					assert.Equal(t, tc.postSensors[i]["unit"], sensor["unit"])
					assert.Equal(t, tc.postSensors[i]["minSafe"], sensor["minSafe"])
					assert.Equal(t, tc.postSensors[i]["maxSafe"], sensor["maxSafe"])
				}
			})

			// UPDATE SENSOR
			t.Run("UpdateSensors", func(t *testing.T) {
				for i, sensor := range tc.putSensors {
					reqBody, err := json.Marshal(sensor)
					assert.NoError(t, err)

					endpoint := "/v1/sensors/" + ids[i]
					statusCode, _, err := cmp.Call(context.Background(), "PUT", endpoint, reqBody)
					assert.NoError(t, err)

					assert.Equal(t, tc.expectedPutStatusCode, statusCode)
				}
			})

			// GET SENSOR
			t.Run("GetSensor", func(t *testing.T) {
				for i, id := range ids {
					endpoint := "/v1/sensors/" + id
					statusCode, resBody, err := cmp.Call(context.Background(), "GET", endpoint, nil)
					assert.NoError(t, err)

					sensor := make(JSON)
					err = json.Unmarshal(resBody, &sensor)
					assert.NoError(t, err)

					assert.Equal(t, tc.expectedGetStatusCode, statusCode)
					assert.Equal(t, id, sensor["id"])
					assert.Equal(t, tc.putSensors[i]["siteId"], sensor["siteId"])
					assert.Equal(t, tc.putSensors[i]["name"], sensor["name"])
					assert.Equal(t, tc.putSensors[i]["unit"], sensor["unit"])
					assert.Equal(t, tc.putSensors[i]["minSafe"], sensor["minSafe"])
					assert.Equal(t, tc.putSensors[i]["maxSafe"], sensor["maxSafe"])
				}
			})

			// DELETE SENSOR
			t.Run("DeleteSensor", func(t *testing.T) {
				for _, id := range ids {
					endpoint := "/v1/sensors/" + id
					statusCode, _, err := cmp.Call(context.Background(), "DELETE", endpoint, nil)
					assert.NoError(t, err)

					assert.Equal(t, tc.expectedDeleteStatusCode, statusCode)
				}
			})
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
			"GetNotExistSensor",
			"GET", "/v1/sensors/0000-0000", nil,
			404, nil,
		},
		{
			"UpdateSensorWithoutSiteId",
			"PUT", "/v1/sensors/0000-0000", JSON{},
			400, nil,
		},
		{
			"UpdateSensorWithoutName",
			"PUT", "/v1/sensors/0000-0000", JSON{"siteId": "1111-aaaa"},
			400, nil,
		},
		{
			"UpdateSensorWithoutUnit",
			"PUT", "/v1/sensors/0000-0000", JSON{"siteId": "1111-aaaa", "name": "temperature"},
			400, nil,
		},
		{
			"UpdateSensorWithInvalidMinSafeAndMaxSafe",
			"PUT", "/v1/sensors/0000-0000", JSON{"siteId": "1111-aaaa", "name": "temperature", "unit": "celsius", "minSafe": 1.0, "maxSafe": 0.0},
			400, nil,
		},
		{
			"UpdateNotExistSensor",
			"PUT", "/v1/sensors/0000-0000", JSON{"siteId": "1111-aaaa", "name": "temperature", "unit": "celsius", "minSafe": -10.0, "maxSafe": 10.0},
			404, nil,
		},
		{
			"DeleteNotExistSensor",
			"DELETE", "/v1/sensors/0000-0000", nil,
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
