package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	tests := []struct {
		name          string
		useFile       bool
		varName       string
		value         string
		defaultValue  string
		expectedValue string
	}{
		{
			"Default",
			false,
			"SECRET",
			"",
			"password",
			"password",
		},
		{
			"EnvironmentVariable",
			false,
			"SECRET",
			"secret",
			"",
			"secret",
		},
		{
			"File",
			true,
			"SECRET",
			"secret",
			"",
			"secret",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.useFile {
				file, err := ioutil.TempFile("", "config-")
				assert.NoError(t, err)
				defer os.Remove(file.Name())

				err = ioutil.WriteFile(file.Name(), []byte(tc.value), 0644)
				assert.NoError(t, err)

				err = os.Setenv(tc.varName+"_FILE", file.Name())
				assert.NoError(t, err)
			} else {
				err := os.Setenv(tc.varName, tc.value)
				assert.NoError(t, err)
			}

			value := getValue(tc.varName, tc.defaultValue)
			assert.Equal(t, tc.expectedValue, value)

			os.Unsetenv(tc.varName)
			os.Unsetenv(tc.varName + "_FILE")
		})
	}
}

func TestGetConfig(t *testing.T) {
	config := GetConfig()

	assert.Equal(t, defaultLogLevel, config.LogLevel)
	assert.Equal(t, defaultServiceName, config.ServiceName)
	assert.Equal(t, defaultServicePort, config.ServicePort)
	assert.Equal(t, defaultRedisURL, config.RedisURL)
}
