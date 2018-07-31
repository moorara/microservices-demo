package transport

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	// GRPCServer is the interface for grpc.Server
	GRPCServer interface {
		Serve(net.Listener) error
		ServeHTTP(http.ResponseWriter, *http.Request)
		Stop()
		GracefulStop()
	}
)

// NewGRPCServer creates a new grpc server
func NewGRPCServer(caFile, certFile, keyFile string, switchService proto.SwitchServiceServer) (GRPCServer, error) {
	opts := []grpc.ServerOption{}

	// Configure MTLS
	if caFile != "" && certFile != "" && keyFile != "" {
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("Failed to append certificate authority")
		}

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    pool,
		}

		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterSwitchServiceServer(grpcServer, switchService)

	return grpcServer, nil
}
