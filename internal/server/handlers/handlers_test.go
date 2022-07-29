package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"yametrics/internal/server/models"
	"yametrics/internal/server/storage"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockMetricStorage struct {
	mock.Mock
}

func getLogger() *zap.SugaredLogger {
	l, _ := zap.NewProduction()
	return l.Sugar()
}

func (s *MockMetricStorage) Get(id string) (models.Metrics, bool) {
	rs := s.Called(id)
	return rs.Get(0).(models.Metrics), rs.Bool(1)
}

func (s *MockMetricStorage) GetAll() []models.Metrics {
	rs := s.Called()
	return rs.Get(0).([]models.Metrics)
}

func (s *MockMetricStorage) Update(m models.Metrics) {}

func TestGetV2(t *testing.T) {
	metricStorage := new(MockMetricStorage)
	v := 1.0
	storedMetric := models.Metrics{ID: "1", MType: "gauge", Value: &v}
	metricStorage.On("Get", "1").Return(storedMetric, true)
	handler := &Handler{getLogger(), metricStorage}
	tests := []struct {
		name     string
		code     int
		body     []byte
		response *models.Metrics
	}{
		{
			"200",
			200,
			func() []byte {
				json, _ := json.Marshal(&models.Metrics{ID: "1", MType: "gauge"})
				return json
			}(),
			&storedMetric,
		},
		{
			"400",
			400,
			func() []byte {
				json, _ := json.Marshal("{}")
				return json
			}(),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/value/1", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.GetV2)
			h.ServeHTTP(w, request)
			res := w.Result()
			var result models.Metrics
			res.Body.Close()
			if tt.response != nil {
				json.NewDecoder(res.Body).Decode(&result)
				assert.Equal(t, true, assert.ObjectsAreEqualValues(result, storedMetric), "wrong response")
			}
			assert.Equal(t, tt.code, res.StatusCode, "wrong status")
		})
	}
}

func TestUpdateV2(t *testing.T) {
	value := 1.0
	handler := &Handler{getLogger(), new(MockMetricStorage)}
	tests := []struct {
		name string
		code int
		body []byte
	}{
		{
			"200",
			200,
			func() []byte {
				json, _ := json.Marshal(&models.Metrics{ID: "1", MType: "gauge", Value: &value})
				return json
			}(),
		},
		{
			"400",
			400,
			func() []byte {
				json, _ := json.Marshal("{}")
				return json
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/gauge/metric_name/1", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.UpdateV2)
			h.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.code, res.StatusCode, "wrong status")
		})
	}
}

func TestUpdateV1(t *testing.T) {
	handler := &Handler{getLogger(), new(MockMetricStorage)}
	tests := []struct {
		name string
		code int
		rctx *chi.Context
	}{
		{"200", 200, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "1")
			rctx.URLParams.Add("type", "gauge")
			return rctx
		}(),
		},
		{"400 Bad Request - incorrect value", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "a")
			rctx.URLParams.Add("type", "gauge")
			return rctx
		}()},
		{"400 Bad Request - empty params", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "")
			rctx.URLParams.Add("value", "")
			rctx.URLParams.Add("type", "gauge")
			return rctx
		}()},
		{"200", 200, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "1")
			rctx.URLParams.Add("type", "counter")
			return rctx
		}(),
		},
		{"400 Bad Request - incorrect value", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "metric_name")
			rctx.URLParams.Add("value", "a")
			rctx.URLParams.Add("type", "counter")
			return rctx
		}(),
		},
		{"400 Bad Request - empty params", 400, func() (c *chi.Context) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("name", "")
			rctx.URLParams.Add("value", "")
			rctx.URLParams.Add("type", "counter")
			return rctx
		}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/", nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, tt.rctx))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.UpdateV1)
			h.ServeHTTP(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.code, res.StatusCode, "wrong status")
		})
	}
}

func TestGetV1(t *testing.T) {
	gaugeType := "gauge"
	counterType := "counter"
	existMetricName := "existMetricName"
	ubsentMetricName := "ubsentMetricName"
	var i int64 = 1
	f := 1.1
	gaugeResponse := "1.1"
	counterResponse := "1"
	model := models.Metrics{ID: "1", MType: models.GAUGE, Value: &f, Delta: &i}
	tests := []struct {
		name          string
		metricStorage storage.MetricsStorage
		code          int
		result        *string
		rctx          *chi.Context
	}{
		{
			"guage, 200 OK",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", existMetricName).Return(model, true)
				return r
			}(),
			200,
			&gaugeResponse,
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", gaugeType)
				rctx.URLParams.Add("name", existMetricName)
				return rctx
			}(),
		},
		{
			"guage, 404 Not Found",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", ubsentMetricName).Return(model, false)
				return r
			}(),
			404,
			nil,
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", gaugeType)
				rctx.URLParams.Add("name", ubsentMetricName)
				return rctx
			}(),
		},
		{
			"guage, 400 Bad Request",
			new(MockMetricStorage),
			400,
			nil,
			chi.NewRouteContext(),
		},
		{
			"counter, 200 OK",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", existMetricName).Return(model, true)
				return r
			}(),
			200,
			&counterResponse,
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", counterType)
				rctx.URLParams.Add("name", existMetricName)
				return rctx
			}(),
		},
		{
			"counter, 404 Not Found",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", ubsentMetricName).Return(model, false)
				return r
			}(),
			404,
			nil,
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
			handler := &Handler{getLogger(), tt.metricStorage}
			h := http.HandlerFunc(handler.GetV1)
			h.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode, "wrong status")

			if tt.result != nil {
				data, err := io.ReadAll(res.Body)
				if assert.NoError(t, err) {
					assert.Equal(t, *tt.result, string(data))
				}
			}
		})
	}
}
