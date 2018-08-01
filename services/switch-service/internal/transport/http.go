package transport

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	// HTTPServer is the interface for http.Server
	HTTPServer interface {
		ListenAndServe() error
		Serve(net.Listener) error
		Close() error
		Shutdown(context.Context) error
	}
)

// NewHTTPServer creates a new http server
func NewHTTPServer(addr string, liveHandler, readyHandler, metricsHandler http.HandlerFunc) HTTPServer {
	httpRouter := mux.NewRouter()
	httpRouter.NotFoundHandler = http.NotFoundHandler()
	httpRouter.Methods("GET").Path("/live").HandlerFunc(liveHandler)
	httpRouter.Methods("GET").Path("/ready").HandlerFunc(readyHandler)
	httpRouter.Methods("GET").Path("/metrics").HandlerFunc(metricsHandler)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: httpRouter,
	}

	return httpServer
}
