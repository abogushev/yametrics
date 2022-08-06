package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"yametrics/internal/protocol"
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

func (s *MockMetricStorage) Get(id string, mtype string) (models.Metrics, bool) {
	rs := s.Called(id, mtype)
	return rs.Get(0).(models.Metrics), rs.Bool(1)
}

func (s *MockMetricStorage) GetAll() []models.Metrics {
	rs := s.Called()
	return rs.Get(0).([]models.Metrics)
}

func (s *MockMetricStorage) Update(m models.Metrics) {}

func TestGetV2(t *testing.T) {
	metricStorage := new(MockMetricStorage)
	storedMetric := models.Metrics{ID: "1", MType: "gauge", Value: 1}
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
				r := protocol.Metrics{ID: "1", MType: "gauge"}
				json, _ := json.Marshal(r)
				metricStorage.On("Get", "1", "gauge").Return(storedMetric, true)
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
				json, _ := json.Marshal(&models.Metrics{ID: "1", MType: "gauge", Value: 1})
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
	existMetricName := "existMetricName"
	ubsentMetricName := "ubsentMetricName"
	model := models.Metrics{ID: "1", MType: models.GAUGE, Value: 1, Delta: 1}
	tests := []struct {
		name          string
		metricStorage storage.MetricsStorage
		code          int
		result        string
		rctx          *chi.Context
	}{
		{
			"guage, 200 OK",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", existMetricName, models.GAUGE).Return(model, true)
				return r
			}(),
			200,
			"1",
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", protocol.GAUGE)
				rctx.URLParams.Add("name", existMetricName)
				return rctx
			}(),
		},
		{
			"guage, 404 Not Found",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", ubsentMetricName, models.GAUGE).Return(model, false)
				return r
			}(),
			404,
			"",
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", protocol.GAUGE)
				rctx.URLParams.Add("name", ubsentMetricName)
				return rctx
			}(),
		},
		{
			"guage, 400 Bad Request",
			new(MockMetricStorage),
			400,
			"",
			chi.NewRouteContext(),
		},
		{
			"counter, 200 OK",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", existMetricName, models.COUNTER).Return(model, true)
				return r
			}(),
			200,
			"1",
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", protocol.COUNTER)
				rctx.URLParams.Add("name", existMetricName)
				return rctx
			}(),
		},
		{
			"counter, 404 Not Found",
			func() storage.MetricsStorage {
				r := new(MockMetricStorage)
				r.On("Get", ubsentMetricName, models.COUNTER).Return(model, false)
				return r
			}(),
			404,
			"",
			func() *chi.Context {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", protocol.COUNTER)
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

			if tt.result != "" {
				data, err := io.ReadAll(res.Body)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.result, string(data))
				}
			}
		})
	}
}
