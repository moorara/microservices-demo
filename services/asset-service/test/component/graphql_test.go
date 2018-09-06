package component

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func makeQuery(client *http.Client, query string) (map[string]interface{}, error) {
	endpoint := Config.ServiceURL + "/graphql"

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(map[string]string{
		"query": query,
	})

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func TestAPI(t *testing.T) {
	if !Config.ComponentTest {
		t.SkipNow()
	}

	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   2 * time.Second,
	}

	tests := []struct {
		name             string
		query            string
		expectedResponse map[string]interface{}
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			response, _ := makeQuery(client, tc.query)

			responseData, err := json.Marshal(response)
			assert.NoError(t, err)

			expectedData, err := json.Marshal(tc.expectedResponse)
			assert.NoError(t, err)

			assert.Equal(t, string(expectedData), string(responseData))
		})
	}
}
