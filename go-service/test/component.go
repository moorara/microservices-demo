package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	timeoutMS = 10000
)

type (
	// JSON represent the json type
	JSON map[string]interface{}

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
		serviceURL = "http://localhost:4010"
	}

	return &Component{
		ServiceURL: serviceURL,
		transport:  &http.Transport{},
	}
}

// Call makes a http request to component under test
func (c *Component) Call(ctx context.Context, method, endpoint string, body JSON) (int, JSON, error) {
	client := &http.Client{
		Transport: c.transport,
		Timeout:   timeoutMS * time.Millisecond,
	}

	bodyData, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	url := c.ServiceURL + endpoint
	bodyReader := bytes.NewReader(bodyData)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, nil, err
	}

	req = req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}

	resJSON := JSON{}
	err = json.Unmarshal(resData, &resJSON)
	if err != nil {
		return res.StatusCode, nil, err
	}

	return res.StatusCode, resJSON, nil
}
