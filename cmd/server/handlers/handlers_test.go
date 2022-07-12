package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockGuageStorage struct{}
type MockCounterStorage struct{}

func (s *MockGuageStorage) Update(name string, value float64) {}
func (s *MockCounterStorage) Update(name string, value int64) {}

func TestHandleGuage(t *testing.T) {
	metricsStorage := &MockGuageStorage{}
	tests := []struct {
		name string
		code int
	}{
		{"return 200 OK", 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/gauge/metric_name/1", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleGuage(metricsStorage))
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.code, "wrong status")
		})
	}
}

func TestHandleCounter(t *testing.T) {
	counterStorage := &MockCounterStorage{}
	tests := []struct {
		name string
		code int
	}{
		{"return 200 OK", 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/counter/metric_name/1", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleCounter(counterStorage))
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.code, "wrong status")
		})
	}
}
