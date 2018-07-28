package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		serverAddr string
	}{
		{
			name:       "MTLSDisabled",
			serverAddr: "localhost:9999",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, conn, err := New(tc.serverAddr)
			defer conn.Close()

			assert.NoError(t, err)
			assert.NotNil(t, client)
			assert.NotNil(t, conn)
		})
	}
}

func TestNewMTLS(t *testing.T) {
	tests := []struct {
		name           string
		serverAddr     string
		serverName     string
		caChainFile    string
		clientCertFile string
		clientKeyFile  string
	}{
		{
			name:           "MTLSEnabled",
			serverAddr:     "localhost:9999",
			serverName:     "server",
			caChainFile:    "../test/ca.chain.cert",
			clientCertFile: "../test/client.cert",
			clientKeyFile:  "../test/client.key",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, conn, err := NewMTLS(tc.serverAddr, tc.serverName, tc.caChainFile, tc.clientCertFile, tc.clientKeyFile)
			defer conn.Close()

			assert.NoError(t, err)
			assert.NotNil(t, client)
			assert.NotNil(t, conn)
		})
	}
}
