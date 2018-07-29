package service

import (
	"context"

	arango "github.com/arangodb/go-driver"
	arangoHttp "github.com/arangodb/go-driver/http"
	arangoUtil "github.com/moorara/microservices-demo/services/switch-service/pkg/arango"
)

const (
	errDatabaseNotFound   = "database not found"
	errCollectionNotFound = "collection not found"
)

type (
	// ArangoService selects the functions used from arango driver
	ArangoService interface {
		Query(ctx context.Context, query string, vars map[string]interface{}) (arango.Cursor, error)
		CreateDocument(ctx context.Context, doc interface{}) (arango.DocumentMeta, error)
		ReadDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error)
		UpdateDocument(ctx context.Context, key string, doc interface{}) (arango.DocumentMeta, error)
		RemoveDocument(ctx context.Context, key string) (arango.DocumentMeta, error)
	}

	arangoService struct {
		arango.Connection
		arango.Client
		arango.Database
		arango.Collection
	}
)

// NewArangoService creates a arango service
func NewArangoService(ctx context.Context, endpoints []string, user, password, databaseName, collectionName string) (ArangoService, error) {
	// Wait until Arango is ready to serve requests
	arangoHTTPService := arangoUtil.NewHTTPService(endpoints[0], 0)
	err := <-arangoHTTPService.NotifyReady(ctx)
	if err != nil {
		return nil, err
	}

	connConfig := arangoHttp.ConnectionConfig{
		Endpoints: endpoints,
	}

	connection, err := arangoHttp.NewConnection(connConfig)
	if err != nil {
		return nil, err
	}

	clientConfig := arango.ClientConfig{
		Connection:     connection,
		Authentication: arango.BasicAuthentication(user, password),
	}

	client, err := arango.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	database, err := client.Database(ctx, databaseName)
	if err != nil {
		if err.Error() != errDatabaseNotFound {
			return nil, err
		}

		// Create database if not exist
		opts := &arango.CreateDatabaseOptions{}
		database, err = client.CreateDatabase(ctx, databaseName, opts)
		if err != nil {
			return nil, err
		}
	}

	collection, err := database.Collection(ctx, collectionName)
	if err != nil {
		if err.Error() != errCollectionNotFound {
			return nil, err
		}

		// Create collection if not exist
		opts := &arango.CreateCollectionOptions{}
		collection, err = database.CreateCollection(ctx, collectionName, opts)
		if err != nil {
			return nil, err
		}
	}

	return &arangoService{
		Connection: connection,
		Client:     client,
		Database:   database,
		Collection: collection,
	}, nil
}
