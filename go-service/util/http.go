package util

import (
	"fmt"
	"net/http"
)

type (
	// ResponseWriter extends the standard http.ResponseWriter
	ResponseWriter interface {
		http.ResponseWriter
		StatusCode() int
		StatusClass() string
	}

	responseWriter struct {
		http.ResponseWriter
		statusCode  int
		statusClass string
	}
)

// NewResponseWriter creates a new response writer
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		ResponseWriter: rw,
	}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)

	if rw.statusCode == 0 {
		rw.statusCode = statusCode
		rw.statusClass = fmt.Sprintf("%dxx", statusCode/100)
	}
}

func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *responseWriter) StatusClass() string {
	return rw.statusClass
}
