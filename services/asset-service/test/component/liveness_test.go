package component

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLiveness(t *testing.T) {
	if Config.ComponentTest {
		t.SkipNow()
	}

	for i := 0; i < 30; i++ {
		resp, err := http.Get(Config.ServiceURL + "/liveness")
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
