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
		Connect(ctx context.Context, endpoints []string, user, password, database, collection string) error
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
func NewArangoService() ArangoService {
	return &arangoService{}
}

func (s *arangoService) Connect(ctx context.Context, endpoints []string, user, password, database, collection string) error {
	// Wait until Arango is ready to serve requests
	arangoHTTPService := arangoUtil.NewHTTPService(endpoints[0], 0)
	err := <-arangoHTTPService.NotifyReady(ctx)
	if err != nil {
		return err
	}

	connConfig := arangoHttp.ConnectionConfig{
		Endpoints: endpoints,
	}

	s.Connection, err = arangoHttp.NewConnection(connConfig)
	if err != nil {
		return err
	}

	clientConfig := arango.ClientConfig{
		Connection:     s.Connection,
		Authentication: arango.BasicAuthentication(user, password),
	}

	s.Client, err = arango.NewClient(clientConfig)
	if err != nil {
		return err
	}

	s.Database, err = s.Client.Database(ctx, database)
	if err != nil {
		if err.Error() != errDatabaseNotFound {
			return err
		}

		// Create database if not exist
		opts := &arango.CreateDatabaseOptions{}
		s.Database, err = s.Client.CreateDatabase(ctx, database, opts)
		if err != nil {
			return err
		}
	}

	s.Collection, err = s.Database.Collection(ctx, collection)
	if err != nil {
		if err.Error() != errCollectionNotFound {
			return err
		}

		// Create collection if not exist
		opts := &arango.CreateCollectionOptions{}
		s.Collection, err = s.Database.CreateCollection(ctx, collection, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
