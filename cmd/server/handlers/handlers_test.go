package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

type MockGuageStorage struct{}
type MockCounterStorage struct{}

func (s *MockGuageStorage) Get(string) (float64, bool)        { return 0, false }
func (s *MockGuageStorage) GetAll() map[string]float64        { return map[string]float64{} }
func (s *MockGuageStorage) Update(name string, value float64) {}

func (s *MockCounterStorage) Get(string) (int64, bool)        { return 0, false }
func (s *MockCounterStorage) GetAll() map[string]int64        { return map[string]int64{} }
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
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "1")

			request := httptest.NewRequest(http.MethodPost, "/update/gauge/metric_name/1", nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostGuage(metricsStorage))
			h.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.code, res.StatusCode, "wrong status")
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
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "1")

			request := httptest.NewRequest(http.MethodPost, "/update/counter/metric_name/1", nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostCounter(counterStorage))
			h.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.code, "wrong status")
		})
	}
}
