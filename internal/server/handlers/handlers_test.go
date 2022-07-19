package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"yametrics/internal/server/storage"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockGuageStorage struct {
	mock.Mock
}
type MockCounterStorage struct {
	mock.Mock
}

func (s *MockGuageStorage) Get(v string) (float64, bool) {
	rs := s.Called(v)
	return rs.Get(0).(float64), rs.Bool(1)
}
func (s *MockGuageStorage) GetAll() map[string]float64 {
	rs := s.Called()
	return rs.Get(0).(map[string]float64)
}
func (s *MockGuageStorage) Update(name string, value float64) {}

func (s *MockCounterStorage) Get(v string) (int64, bool) {
	rs := s.Called(v)
	return rs.Get(0).(int64), rs.Bool(1)
}
func (s *MockCounterStorage) GetAll() map[string]int64 {
	rs := s.Called()
	return rs.Get(0).(map[string]int64)
}
func (s *MockCounterStorage) Update(name string, value int64) {}

func getLogger() *zap.SugaredLogger {
	l, _ := zap.NewProduction()
	return l.Sugar()
}

func TestPostGuage(t *testing.T) {
	handler := &Handler{new(MockGuageStorage), new(MockCounterStorage), getLogger()}
	tests := []struct {
		name string
		code int
		rctx *chi.Context
	}{
		{"200", 200, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "1")
			return rctx
		}(),
		},
		{"400 Bad Request - incorrect value", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "a")
			return rctx
		}()},
		{"400 Bad Request - empty params", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "")
			rctx.URLParams.Add("value", "")
			return rctx
		}()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/gauge/metric_name/1", nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, tt.rctx))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.PostGuage)
			h.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.code, res.StatusCode, "wrong status")
		})
	}
}

func TestPostCounter(t *testing.T) {
	handler := &Handler{new(MockGuageStorage), new(MockCounterStorage), getLogger()}

	tests := []struct {
		name string
		code int
		rctx *chi.Context
	}{
		{"200", 200, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "1")
			return rctx
		}(),
		},
		{"400 Bad Request - incorrect value", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "a")
			return rctx
		}(),
		},
		{"400 Bad Request - empty params", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "")
			rctx.URLParams.Add("value", "")
			return rctx
		}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/gauge/metric_name/1", nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, tt.rctx))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.PostCounter)
			h.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.code, res.StatusCode, "wrong status")
		})
	}
}

func TestGetMetric(t *testing.T) {
	gaugeType := "gauge"
	counterType := "counter"
	dummyCounterStorage := new(MockCounterStorage)
	dummyGuageStorage := new(MockGuageStorage)
	existMetricName := "existMetricName"
	ubsentMetricName := "ubsentMetricName"

	tests := []struct {
		name           string
		guageStorage   storage.GuageStorage
		counterStorage storage.CounterStorage
		code           int
		rctx           *chi.Context
	}{
		{
			"guage, 200 OK",
			func() storage.GuageStorage {
				r := new(MockGuageStorage)
				r.On("Get", existMetricName).Return(float64(1), true)
				return r
			}(),
			dummyCounterStorage,
			200,
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", gaugeType)
				rctx.URLParams.Add("name", existMetricName)
				return rctx
			}(),
		},
		{
			"guage, 404 Not Found",
			func() storage.GuageStorage {
				r := new(MockGuageStorage)
				r.On("Get", ubsentMetricName).Return(0.0, false)
				return r
			}(),
			dummyCounterStorage,
			404,
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", gaugeType)
				rctx.URLParams.Add("name", ubsentMetricName)
				return rctx
			}(),
		},
		{
			"guage, 400 Bad Request",
			dummyGuageStorage,
			dummyCounterStorage,
			400,
			chi.NewRouteContext(),
		},
		{
			"counter, 200 OK",
			dummyGuageStorage,
			func() storage.CounterStorage {
				r := new(MockCounterStorage)
				r.On("Get", existMetricName).Return(int64(1), true)
				return r
			}(),
			200,
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", counterType)
				rctx.URLParams.Add("name", existMetricName)
				return rctx
			}(),
		},
		{
			"counter, 404 Not Found",
			dummyGuageStorage,
			func() storage.CounterStorage {
				r := new(MockCounterStorage)
				r.On("Get", ubsentMetricName).Return(int64(0), false)
				return r
			}(),
			404,
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", counterType)
				rctx.URLParams.Add("name", ubsentMetricName)
				return rctx
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, tt.rctx))
			w := httptest.NewRecorder()
			handler := &Handler{tt.guageStorage, tt.counterStorage, getLogger()}
			h := http.HandlerFunc(handler.GetMetric)
			h.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.code, res.StatusCode, "wrong status")
		})
	}
}
