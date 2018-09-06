package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/moorara/microservices-demo/services/asset-service/internal/db"
	"github.com/stretchr/testify/assert"

	arango "github.com/arangodb/go-driver"
	arangoUtil "github.com/moorara/microservices-demo/services/asset-service/pkg/arango"
)

func TestArangoHTTPService(t *testing.T) {
	if !Config.IntegrationTest {
		t.SkipNow()
	}

	databaseName := fmt.Sprintf("animals_%d", time.Now().Unix())
	docCollectionName := fmt.Sprintf("mammals_%d", time.Now().Unix())
	edgeCollectionName := fmt.Sprintf("relations_%d", time.Now().Unix())
	documentKey := "4444"

	tests := []struct {
		name                   string
		method, endpoint, body string
		expectedStatus         int
	}{
		{
			name:           "GetDatabaseNames",
			method:         "GET",
			endpoint:       "/_api/database",
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "CreateDatabase",
			method:         "POST",
			endpoint:       "/_api/database",
			body:           fmt.Sprintf(`{"name":"%s"}`, databaseName),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "GetDatabaseInfo",
			method:         "GET",
			endpoint:       fmt.Sprintf("/_db/%s/_api/database/current", databaseName),
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "CreateDocumentCollection",
			method:         "POST",
			endpoint:       fmt.Sprintf("/_db/%s/_api/collection", databaseName),
			body:           fmt.Sprintf(`{"name":"%s"}`, docCollectionName),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetDocumentCollection",
			method:         "GET",
			endpoint:       fmt.Sprintf("/_db/%s/_api/collection/%s", databaseName, docCollectionName),
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "CreateEdgeCollection",
			method:         "POST",
			endpoint:       fmt.Sprintf("/_db/%s/_api/collection", databaseName),
			body:           fmt.Sprintf(`{"name":"%s", "type":3}`, edgeCollectionName),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetEdgeCollection",
			method:         "GET",
			endpoint:       fmt.Sprintf("/_db/%s/_api/collection/%s", databaseName, edgeCollectionName),
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "CreateDocument",
			method:         "POST",
			endpoint:       fmt.Sprintf("/_db/%s/_api/document/%s", databaseName, docCollectionName),
			body:           fmt.Sprintf(`{"_key":"%s", "name":"chipmunk", "class":"mammalia"}`, documentKey),
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "GetDocument",
			method:         "GET",
			endpoint:       fmt.Sprintf("/_db/%s/_api/document/%s/%s", databaseName, docCollectionName, documentKey),
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DeleteDocument",
			method:         "DELETE",
			endpoint:       fmt.Sprintf("/_db/%s/_api/document/%s/%s", databaseName, docCollectionName, documentKey),
			body:           ``,
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "DeleteDocumentCollection",
			method:         "DELETE",
			endpoint:       fmt.Sprintf("/_db/%s/_api/collection/%s", databaseName, docCollectionName),
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DeleteEdgeCollection",
			method:         "DELETE",
			endpoint:       fmt.Sprintf("/_db/%s/_api/collection/%s", databaseName, edgeCollectionName),
			body:           ``,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DeleteDatabase",
			method:         "DELETE",
			endpoint:       fmt.Sprintf("/_api/database/%s", databaseName),
			body:           ``,
			expectedStatus: http.StatusOK,
		},
	}

	arangoHTTPService := arangoUtil.NewHTTPService(Config.ArangoHTTPAddr)
	assert.NotNil(t, arangoHTTPService)

	t.Run("NotifyReady", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := <-arangoHTTPService.NotifyReady(ctx)
		assert.NoError(t, err)
	})

	t.Run("Login", func(t *testing.T) {
		ctx := context.Background()
		err := arangoHTTPService.Login(ctx, Config.ArangoUser, Config.ArangoPassword)
		assert.NoError(t, err)
	})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			data, status, err := arangoHTTPService.Call(ctx, tc.method, tc.endpoint, tc.body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, status)
			assert.NotNil(t, data)
		})
	}
}

func TestArangoService(t *testing.T) {
	if !Config.IntegrationTest {
		t.SkipNow()
	}

	tests := []struct {
		name           string
		databaseName   string
		collectionName string
		document       map[string]interface{}
		update         map[string]interface{}
		query          string
		vars           map[string]interface{}
	}{
		{
			name:           "MammalsCollection",
			databaseName:   "animals",
			collectionName: "mammals",
			document: map[string]interface{}{
				"name":  "moose",
				"class": "mammalia",
			},
			update: map[string]interface{}{
				"family": "cervidae",
			},
			query: `INSERT { name: @name, class: @class } INTO mammals RETURN NEW`,
			vars: map[string]interface{}{
				"name":  "sloth",
				"class": "mammalia",
			},
		},
		{
			name:           "BirdsCollection",
			databaseName:   "animals",
			collectionName: "birds",
			document: map[string]interface{}{
				"name":  "parrot",
				"class": "aves",
			},
			update: map[string]interface{}{
				"phylum": "chordata",
			},
			query: `INSERT { name: @name, class: @class } INTO birds RETURN NEW`,
			vars: map[string]interface{}{
				"name":  "flamingo",
				"class": "aves",
			},
		},
		{
			name:           "ReptilesCollection",
			databaseName:   "animals",
			collectionName: "reptiles",
			document: map[string]interface{}{
				"name":  "chameleon",
				"class": "reptilia",
			},
			update: map[string]interface{}{
				"order": "squamata",
			},
			query: `INSERT { name: @name, class: @class } INTO reptiles RETURN NEW`,
			vars: map[string]interface{}{
				"name":  "salamander",
				"class": "amphibia",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var arangoService db.ArangoService

			var doc map[string]interface{}
			var docMeta arango.DocumentMeta

			ctx := context.Background()

			t.Run("CreateArangoService", func(t *testing.T) {
				arangoService, err = db.NewArangoService(Config.ArangoEndpoints, Config.ArangoUser, Config.ArangoPassword)
				assert.NoError(t, err)
				assert.NotNil(t, arangoService)
			})

			t.Run("Connect", func(t *testing.T) {
				err = arangoService.Connect(ctx, tc.databaseName, tc.collectionName)
				assert.NoError(t, err)
			})

			t.Run("CreateDocument", func(t *testing.T) {
				docMeta, err = arangoService.CreateDocument(ctx, tc.document)
				assert.NoError(t, err)
				assert.NotEmpty(t, docMeta)
			})

			t.Run("ReadDocument", func(t *testing.T) {
				_, err = arangoService.ReadDocument(ctx, docMeta.Key, &doc)
				assert.NoError(t, err)
			})

			t.Run("UpdateDocument", func(t *testing.T) {
				docMeta, err = arangoService.UpdateDocument(ctx, docMeta.Key, tc.update)
				assert.NoError(t, err)
				assert.NotEmpty(t, docMeta)
			})

			t.Run("RemoveDocument", func(t *testing.T) {
				docMeta, err = arangoService.RemoveDocument(ctx, docMeta.Key)
				assert.NoError(t, err)
				assert.NotEmpty(t, docMeta)
			})

			t.Run("Query", func(t *testing.T) {
				cursor, err := arangoService.Query(ctx, tc.query, tc.vars)
				assert.NoError(t, err)
				defer cursor.Close()
			})
		})
	}
}
