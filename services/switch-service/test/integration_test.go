package test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/switch-service/internal/service"
	"github.com/moorara/microservices-demo/services/switch-service/pkg/arango"
	"github.com/stretchr/testify/assert"
)

const (
	waitTime = 20 * time.Second
)

func integrationTest() bool {
	value := os.Getenv("INTEGRATION_TEST")
	return value == "true" || value == "TRUE"
}

func arangoConfig() (string, string, string) {
	address := os.Getenv("ARANGO_ADDR")
	if address == "" {
		address = "http://localhost:8529"
	}

	user := os.Getenv("ARANGO_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("ARANGO_PASSWORD")
	if password == "" {
		password = "pass"
	}

	return address, user, password
}

func TestArangoHTTPService(t *testing.T) {
	if !integrationTest() {
		t.SkipNow()
	}

	address, user, password := arangoConfig()

	tests := []struct {
		name                   string
		timeout                time.Duration
		method, endpoint, body string
		expectedStatus         int
	}{
		{
			name:           "CreateDatabase",
			timeout:        time.Second,
			method:         "POST",
			endpoint:       "/_api/database",
			body:           fmt.Sprintf(`{"name":"example-%d"}`, time.Now().Unix()),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "GetDatabases",
			timeout:        time.Second,
			method:         "GET",
			endpoint:       "/_api/database",
			body:           ``,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arangoHTTPService := arango.NewHTTPService(address, tc.timeout)

			t.Run("NotifyReady", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), waitTime)
				defer cancel()
				err := <-arangoHTTPService.NotifyReady(ctx)
				assert.NoError(t, err)
			})

			t.Run("Login", func(t *testing.T) {
				err := arangoHTTPService.Login(user, password)
				assert.NoError(t, err)
			})

			t.Run("Call", func(t *testing.T) {
				ctx := context.Background()
				data, status, err := arangoHTTPService.Call(ctx, tc.method, tc.endpoint, tc.body)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, status)
				assert.NotNil(t, data)
			})
		})
	}
}

func TestArangoService(t *testing.T) {
	if !integrationTest() {
		t.SkipNow()
	}

	address, user, password := arangoConfig()

	tests := []struct {
		name           string
		databaseName   string
		collectionName string
	}{
		{
			name:           "MammalsCollection",
			databaseName:   "animals",
			collectionName: "mammals",
		},
		{
			name:           "BirdsCollection",
			databaseName:   "animals",
			collectionName: "birds",
		},
		{
			name:           "ReptilesCollection",
			databaseName:   "animals",
			collectionName: "reptiles",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var arangoService service.ArangoService

			t.Run("CreateArangoService", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()

				arangoService = service.NewArangoService()
				err = arangoService.Connect(ctx, []string{address}, user, password, tc.databaseName, tc.collectionName)

				assert.NoError(t, err)
				assert.NotNil(t, arangoService)
			})
		})
	}
}
