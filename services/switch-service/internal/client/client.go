package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"github.com/moorara/microservices-demo/services/switch-service/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// New creates a new client without MTLS
func New(serverAddr string) (proto.SwitchServiceClient, *grpc.ClientConn, error) {
	options := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(serverAddr, options...)
	if err != nil {
		return nil, nil, err
	}

	client := proto.NewSwitchServiceClient(conn)

	return client, conn, nil
}

// NewMTLS creates a new client with MTLS enabled
func NewMTLS(serverAddr, serverName, caChainFile, clientCertFile, clientKeyFile string) (proto.SwitchServiceClient, *grpc.ClientConn, error) {
	ca, err := ioutil.ReadFile(caChainFile)
	if err != nil {
		return nil, nil, err
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(ca); !ok {
		return nil, nil, errors.New("Failed to append certificate authority")
	}

	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ServerName:   serverName,
	}

	creds := credentials.NewTLS(tlsConfig)
	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	conn, err := grpc.Dial(serverAddr, options...)
	if err != nil {
		return nil, nil, err
	}

	client := proto.NewSwitchServiceClient(conn)

	return client, conn, nil
}
