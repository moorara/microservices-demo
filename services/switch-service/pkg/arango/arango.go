package arango

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	waitTime = time.Second
)

type (
	// JSON is an alias for map
	JSON map[string]interface{}

	// HTTPService is for making REST calls to Aragno
	HTTPService interface {
		NotifyReady(ctx context.Context) chan error
		Login(user, password string) error
		Call(ctx context.Context, method, endpoint, body string) (JSON, int, error)
	}

	httpService struct {
		client  *http.Client
		address string
		jwt     string
	}
)

// NewHTTPService creates a new arrango http
func NewHTTPService(address string, timeout time.Duration) HTTPService {
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   timeout,
	}

	return &httpService{
		client:  client,
		address: address,
	}
}

func (s *httpService) NotifyReady(ctx context.Context) chan error {
	ch := make(chan error)

	go func() {
		for {
			url := s.address + "/_admin/server/availability"
			req, _ := http.NewRequest("GET", url, nil)

			resp, err := s.client.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					ch <- nil
					return
				}
			}

			select {
			case <-time.After(waitTime):
			case <-ctx.Done():
				ch <- ctx.Err()
				return
			}
		}
	}()

	return ch
}

func (s *httpService) Login(user, password string) error {
	method := "POST"
	url := s.address + "/_open/auth"
	body := strings.NewReader(
		fmt.Sprintf(`{"username":"%s", "password":"%s"}`, user, password),
	)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d authentication failed", resp.StatusCode)
	}

	data := struct{ Jwt string }{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}

	s.jwt = data.Jwt

	return nil
}

func (s *httpService) Call(ctx context.Context, method, endpoint, body string) (JSON, int, error) {
	url := s.address + endpoint
	bodyReader := strings.NewReader(body)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+s.jwt)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var data JSON
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, 0, err
	}

	return data, resp.StatusCode, nil
}
