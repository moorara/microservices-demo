package test

import (
	"os"
)

func integrationTest() bool {
	value := os.Getenv("INTEGRATION_TEST")
	return value == "true" || value == "TRUE"
}
