package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"yametrics/internal/server/storage"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Handler struct {
	GuageStorage   storage.GuageStorage
	CounterStorage storage.CounterStorage
	Logger         *zap.SugaredLogger
}

func (h *Handler) PostGuage(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if v, err := strconv.ParseFloat(value, 64); err == nil && name != "" {
		h.GuageStorage.Update(name, v)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	} else {
		h.Logger.Errorf("wrong params: name=%v, value=%v", name, value)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) PostCounter(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")
	if v, error := strconv.ParseInt(value, 10, 64); error == nil && name != "" {
		h.CounterStorage.Update(name, v)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	} else {
		h.Logger.Errorf("wrong params: name=%v, value=%v", name, value)
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	var found bool
	var result string

	switch metricType {
	case "gauge":
		var v float64
		v, found = h.GuageStorage.Get(metricName)
		if found {
			result = fmt.Sprintf("%v", v)
		}
	case "counter":
		var v int64
		v, found = h.CounterStorage.Get(metricName)
		if found {
			result = fmt.Sprintf("%v", v)
		}
	default:
		h.Logger.Errorf("wrong metric type: %v", metricType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if found {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func (h *Handler) GetAllAsHTML(w http.ResponseWriter, r *http.Request) {
	allmtrcs := []string{}
	for n, v := range h.GuageStorage.GetAll() {
		allmtrcs = append(allmtrcs, fmt.Sprintf("name: %v value: %v", n, v))
	}
	for n, v := range h.CounterStorage.GetAll() {
		allmtrcs = append(allmtrcs, fmt.Sprintf("name: %v value: %v", n, v))
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
		h.Logger.Error("Error Parsing template: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err1 := tmpl.Execute(w, allmtrcs)
	if err1 != nil {
		h.Logger.Error("Error executing template: ", err1)
	}
}
