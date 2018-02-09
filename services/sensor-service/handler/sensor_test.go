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
	"github.com/moorara/microservices-demo/services/sensor-service/service"
	"github.com/stretchr/testify/assert"
)

type mockSensorManager struct {
	CreateCalled bool
	CreateResult *service.Sensor
	CreateError  error

	AllCalled bool
	AllResult []service.Sensor
	AllError  error

	GetCalled bool
	GetResult *service.Sensor
	GetError  error

	DeleteCalled bool
	DeleteError  error
}

func (m *mockSensorManager) Create(ctx context.Context, siteID string, name, unit string, minSafe, maxSafe float64) (*service.Sensor, error) {
	m.CreateCalled = true
	return m.CreateResult, m.CreateError
}

func (m *mockSensorManager) All(ctx context.Context, siteID string) ([]service.Sensor, error) {
	m.AllCalled = true
	return m.AllResult, m.AllError
}

func (m *mockSensorManager) Get(ctx context.Context, id string) (*service.Sensor, error) {
	m.GetCalled = true
	return m.GetResult, m.GetError
}

func (m *mockSensorManager) Delete(ctx context.Context, id string) error {
	m.DeleteCalled = true
	return m.DeleteError
}

func TestNewSensorHandler(t *testing.T) {
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
			h := NewSensorHandler(db, logger)

			assert.NotNil(t, h)
		})
	}
}

func TestPostSensor(t *testing.T) {
	tests := []struct {
		name            string
		createResult    *service.Sensor
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
			`{"siteId": "1111-aaaa"`,
			400,
			``,
		},
		{
			"SensorManagerError",
			nil, errors.New("error"),
			`{"siteId": "1111-aaaa", "name": "temperature", "unit": "celsius", "minSafe": -30, "maxSafe": 30}`,
			500,
			``,
		},
		{
			"Successful",
			&service.Sensor{
				ID:      "2222-bbbb",
				SiteID:  "1111-aaaa",
				Name:    "temperature",
				Unit:    "celsius",
				MinSafe: -30.0,
				MaxSafe: 30.0,
			},
			nil,
			`{"siteId": "1111-aaaa", "name": "temperature", "unit": "celsius", "minSafe": -30, "maxSafe": 30}`,
			201,
			`{"id":"2222-bbbb","siteId":"1111-aaaa","name":"temperature","unit":"celsius","minSafe":-30,"maxSafe":30}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockSensorManager{
				CreateResult: tc.createResult,
				CreateError:  tc.createError,
			}

			h := &postgresSensorHandler{
				manager: m,
				logger:  log.NewNopLogger(),
			}

			reqBody := strings.NewReader(tc.reqBody)
			r := httptest.NewRequest("POST", "http://service/sensors", reqBody)
			w := httptest.NewRecorder()

			h.PostSensor(w, r)
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

func TestGetSensors(t *testing.T) {
	tests := []struct {
		name            string
		allResult       []service.Sensor
		allError        error
		siteID          string
		expectedStatus  int
		expectedResBody string
	}{
		{
			"NoSiteID",
			nil, nil,
			"",
			400,
			``,
		},
		{
			"SensorManagerError",
			nil, errors.New("error"),
			"1111-aaaa",
			500,
			``,
		},
		{
			"Successful",
			[]service.Sensor{
				{ID: "2222-bbbb", SiteID: "1111-aaaa", Name: "temperature", Unit: "celsius", MinSafe: -30.0, MaxSafe: 30.0},
				{ID: "4444-dddd", SiteID: "1111-aaaa", Name: "temperature", Unit: "fahrenheit", MinSafe: -22.0, MaxSafe: 86.0},
			},
			nil,
			"1111-aaaa",
			200,
			`[{"id":"2222-bbbb","siteId":"1111-aaaa","name":"temperature","unit":"celsius","minSafe":-30,"maxSafe":30},{"id":"4444-dddd","siteId":"1111-aaaa","name":"temperature","unit":"fahrenheit","minSafe":-22,"maxSafe":86}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockSensorManager{
				AllResult: tc.allResult,
				AllError:  tc.allError,
			}

			h := &postgresSensorHandler{
				manager: m,
				logger:  log.NewNopLogger(),
			}

			mr := mux.NewRouter()
			mr.Methods("GET").Path("/sensors").Queries("siteId", "{siteId}").HandlerFunc(h.GetSensors)
			ts := httptest.NewServer(mr)
			defer ts.Close()

			res, err := http.Get(ts.URL + "/sensors?siteId=" + tc.siteID)
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

func TestGetSensor(t *testing.T) {
	tests := []struct {
		name            string
		getResult       *service.Sensor
		getError        error
		sensorID        string
		expectedStatus  int
		expectedResBody string
	}{
		{
			"NoSensorID",
			nil, nil,
			"",
			404,
			``,
		},
		{
			"SensorManagerError",
			nil, errors.New("error"),
			"2222-bbbb",
			500,
			``,
		},
		{
			"NoSensorFound",
			nil, nil,
			"2222-bbbb",
			404,
			``,
		},
		{
			"Successful",
			&service.Sensor{
				ID:      "2222-bbbb",
				SiteID:  "1111-aaaa",
				Name:    "temperature",
				Unit:    "celsius",
				MinSafe: -30.0,
				MaxSafe: 30.0,
			},
			nil,
			"1111-aaaa",
			200,
			`{"id":"2222-bbbb","siteId":"1111-aaaa","name":"temperature","unit":"celsius","minSafe":-30,"maxSafe":30}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockSensorManager{
				GetResult: tc.getResult,
				GetError:  tc.getError,
			}

			h := &postgresSensorHandler{
				manager: m,
				logger:  log.NewNopLogger(),
			}

			mr := mux.NewRouter()
			mr.Methods("GET").Path("/sensors/{id}").HandlerFunc(h.GetSensor)
			ts := httptest.NewServer(mr)
			defer ts.Close()

			res, err := http.Get(ts.URL + "/sensors/" + tc.sensorID)
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

func TestDeleteSensor(t *testing.T) {
	tests := []struct {
		name            string
		deleteError     error
		sensorID        string
		expectedStatus  int
		expectedResBody string
	}{
		{
			"NoSensorID",
			nil,
			"",
			404,
			``,
		},
		{
			"SensorManagerError",
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
			m := &mockSensorManager{
				DeleteError: tc.deleteError,
			}

			h := &postgresSensorHandler{
				manager: m,
				logger:  log.NewNopLogger(),
			}

			mr := mux.NewRouter()
			mr.Methods("DELETE").Path("/sensors/{id}").HandlerFunc(h.DeleteSensor)
			ts := httptest.NewServer(mr)
			defer ts.Close()

			req, err := http.NewRequest("DELETE", ts.URL+"/sensors/"+tc.sensorID, nil)
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
