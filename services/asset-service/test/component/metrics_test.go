package component

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	if Config.ComponentTest {
		t.SkipNow()
	}

	url := Config.ServiceURL + "/metrics"
	res, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	body := string(data)

	assert.Contains(t, body, "# TYPE go_goroutines gauge")
	assert.Contains(t, body, "# TYPE go_gc_duration_seconds summary")
	assert.Contains(t, body, "# TYPE go_memstats_heap_objects gauge")
	assert.Contains(t, body, "# TYPE go_memstats_alloc_bytes gauge")
	assert.Contains(t, body, "# TYPE go_memstats_alloc_bytes_total counter")
	assert.Contains(t, body, "# TYPE asset_service_process_open_fds gauge")
	assert.Contains(t, body, "# TYPE asset_service_process_resident_memory_bytes gauge")
	assert.Contains(t, body, "# TYPE asset_service_process_virtual_memory_bytes gauge")
	assert.Contains(t, body, "# TYPE asset_service_jaeger_traces counter")
	assert.Contains(t, body, "# TYPE asset_service_jaeger_started_spans counter")
	assert.Contains(t, body, "# TYPE asset_service_jaeger_finished_spans counter")
}
