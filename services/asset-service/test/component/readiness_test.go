package component

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadiness(t *testing.T) {
	if Config.ComponentTest {
		t.SkipNow()
	}

	for i := 0; i < 30; i++ {
		resp, err := http.Get(Config.ServiceURL + "/readiness")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(time.Second)
	}

	assert.Fail(t, "timeout")
}
