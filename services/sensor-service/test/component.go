package test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	timeoutMS = 10000
)

type (
	// Component represents the component under test
	Component struct {
		ServiceURL string
		transport  *http.Transport
	}
)

// NewComponent creates a new component for testing
func NewComponent() *Component {
	serviceURL := os.Getenv("SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "http://localhost:4020"
	}

	return &Component{
		ServiceURL: serviceURL,
		transport:  &http.Transport{},
	}
}

// Call makes a http request to component under test
func (c *Component) Call(ctx context.Context, method, endpoint string, body []byte) (int, []byte, error) {
	client := &http.Client{
		Transport: c.transport,
		Timeout:   timeoutMS * time.Millisecond,
	}

	url := c.ServiceURL + endpoint
	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}

	return res.StatusCode, resBody, nil
}
