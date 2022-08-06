package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"yametrics/internal/protocol"
	"yametrics/internal/server/models"
	"yametrics/internal/server/storage"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Handler struct {
	logger         *zap.SugaredLogger
	metricsStorage storage.MetricsStorage
}

func NewHandler(logger *zap.SugaredLogger, metricsStorage storage.MetricsStorage) Handler {
	return Handler{logger: logger, metricsStorage: metricsStorage}
}

func (h *Handler) UpdateV2(w http.ResponseWriter, r *http.Request) {
	var metric protocol.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.metricsStorage.Update(models.Metrics{ID: metric.ID, MType: metric.MType, Delta: *metric.Delta, Value: *metric.Value})
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetV2(w http.ResponseWriter, r *http.Request) {
	var metric protocol.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric, found := h.metricsStorage.Get(metric.ID, metric.MType); found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(metric)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *Handler) UpdateV1(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")
	mtype := chi.URLParam(r, "type")

	var reqError error
	var metric models.Metrics

	if name != "" {
		switch mtype {
		case protocol.GAUGE:
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				metric = models.Metrics{ID: name, MType: protocol.GAUGE, Value: f}
			} else {
				reqError = fmt.Errorf("wrong gauge param: %v", value)
			}
		case protocol.COUNTER:
			if f, err := strconv.ParseInt(value, 10, 64); err == nil {
				metric = models.Metrics{ID: name, MType: protocol.COUNTER, Delta: f}
			} else {
				reqError = fmt.Errorf("wrong counter param: %v", value)
			}
		}
	} else {
		reqError = errors.New("param `name` must be nonempty")
	}

	if reqError == nil {
		h.metricsStorage.Update(metric)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	} else {
		h.logger.Error(reqError)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) GetV1(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if metricType != models.COUNTER && metricType != models.GAUGE {
		h.logger.Errorf("wrong metric type: %v", metricType)
		w.WriteHeader(http.StatusBadRequest)
	} else if metric, ok := h.metricsStorage.Get(metricName, metricType); !ok {
		w.WriteHeader(http.StatusNotFound)
	} else {
		var result string
		if metric.MType == models.GAUGE {
			result = fmt.Sprintf("%v", metric.Value)
		} else {
			result = fmt.Sprintf("%v", metric.Delta)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

func (h *Handler) GetAllAsHTML(w http.ResponseWriter, r *http.Request) {
	storageMetrics := h.metricsStorage.GetAll()
	allmtrcs := make([]string, len(storageMetrics))

	for v, i := "", 0; i < len(storageMetrics); i++ {
		if storageMetrics[i].MType == protocol.GAUGE {
			v = fmt.Sprintf("%v", storageMetrics[i].Value)
		} else {
			v = fmt.Sprintf("%v", storageMetrics[i].Delta)
		}
		allmtrcs[i] = fmt.Sprintf("name: %v value: %v", storageMetrics[i].ID, v)
	}

	tmpl, err := template.New("test").Parse(`
		<html>
			<head>
			<title></title>
			</head>
			<body>
			{{ range $key, $value := . }}
			<li>{{ $value }}</li>
			{{ end }}
			</body>
		</html>`)

	if err != nil {
		h.logger.Error("Error Parsing template: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err1 := tmpl.Execute(w, allmtrcs)
	if err1 != nil {
		h.logger.Error("Error executing template: ", err1)
	}
}
