package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"yametrics/internal/agent/models/api"
	"yametrics/internal/agent/models/storage"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAgent_sendMetricsV1(t *testing.T) {
	urls := []string{
		"/update/gauge/Alloc/0",
		"/update/gauge/BuckHashSys/0",
		"/update/gauge/Frees/0",
		"/update/gauge/GCCPUFraction/0",
		"/update/gauge/GCSys/0",
		"/update/gauge/HeapAlloc/0",
		"/update/gauge/HeapIdle/0",
		"/update/gauge/HeapInuse/0",
		"/update/gauge/HeapObjects/0",
		"/update/gauge/HeapReleased/0",
		"/update/gauge/HeapSys/0",
		"/update/gauge/LastGC/0",
		"/update/gauge/Lookups/0",
		"/update/gauge/MCacheInuse/0",
		"/update/gauge/MCacheSys/0",
		"/update/gauge/MSpanInuse/0",
		"/update/gauge/MSpanSys/0",
		"/update/gauge/Mallocs/0",
		"/update/gauge/NextGC/0",
		"/update/gauge/NumForcedGC/0",
		"/update/gauge/NumGC/0",
		"/update/gauge/OtherSys/0",
		"/update/gauge/PauseTotalNs/0",
		"/update/gauge/StackInuse/0",
		"/update/gauge/StackSys/0",
		"/update/gauge/Sys/0",
		"/update/gauge/TotalAlloc/0",
		"/update/gauge/RandomValue/0",
		"/update/counter/PollCount/0",
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.Method, "POST")

		assert.Contains(t, urls, req.URL.String())
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	logger, _ := zap.NewProduction()
	fmt.Println(server.URL)
	agent := &Agent{
		url:     server.URL,
		client:  *server.Client(),
		logger:  logger.Sugar(),
		metrics: &storage.Metrics{MemStats: &runtime.MemStats{}, PollCount: 0, RandomValue: 0.0},
	}
	agent.sendMetricsV1()
}

func TestAgent_sendMetricsV2(t *testing.T) {
	results := &storage.Metrics{MemStats: &runtime.MemStats{}, PollCount: 0, RandomValue: 0.0}
	apiModels := results.ToAPI()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.Method, "POST")
		var metric api.Metrics
		err := json.NewDecoder(req.Body).Decode(&metric)
		if assert.NoError(t, err) {
			assert.Contains(t, apiModels, metric)
		}
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	logger, _ := zap.NewProduction()
	agent := &Agent{
		url:     server.URL,
		client:  *server.Client(),
		logger:  logger.Sugar(),
		metrics: results,
	}
	agent.sendMetricsV2()
}
