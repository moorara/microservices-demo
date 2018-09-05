package db

import (
	"context"

	arango "github.com/arangodb/go-driver"
	arangoHttp "github.com/arangodb/go-driver/http"
)

type (
	// ArangoService selects the functions used from arango driver
	ArangoService interface {
		Connect(ctx context.Context, database, collection string) error
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

// NewArangoService creates an Arango service
func NewArangoService(endpoints []string, user, password string) (ArangoService, error) {
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

	return &arangoService{
		Connection: connection,
		Client:     client,
	}, nil
}

func (s *arangoService) Connect(ctx context.Context, databaseName, collectionName string) error {
	database, err := s.Client.Database(ctx, databaseName)
	if err != nil {
		if !arango.IsNotFound(err) {
			return err
		}

		// Create database if not exist
		opts := &arango.CreateDatabaseOptions{}
		database, err = s.Client.CreateDatabase(ctx, databaseName, opts)
		if err != nil {
			return err
		}
	}

	collection, err := database.Collection(ctx, collectionName)
	if err != nil {
		if !arango.IsNotFound(err) {
			return err
		}

		// Create collection if not exist
		opts := &arango.CreateCollectionOptions{}
		collection, err = database.CreateCollection(ctx, collectionName, opts)
		if err != nil {
			return err
		}
	}

	s.Database = database
	s.Collection = collection

	return nil
}
