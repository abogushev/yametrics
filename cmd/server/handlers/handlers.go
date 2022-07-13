package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"yametrics/cmd/server/storage"

	"github.com/go-chi/chi"
)

func PostGuage(storage storage.GuageStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			storage.Update(name, v)

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func PostCounter(storage storage.CounterStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		if v, error := strconv.ParseInt(value, 10, 64); error == nil {
			storage.Update(name, v)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func GetMetric(guageStorage storage.GuageStorage, counterStorage storage.CounterStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		var found bool
		var result string

		switch metricType {
		case "gauge":
			var v float64
			v, found = guageStorage.Get(metricName)
			if found {
				result = fmt.Sprintf("%v", v)
			}
		case "counter":
			var v int64
			v, found = counterStorage.Get(metricName)
			if found {
				result = fmt.Sprintf("%v", v)
			}
		}
		if found {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func GetAllAsHTML(guageStorage storage.GuageStorage, counterStorage storage.CounterStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allmtrcs := []string{}
		for n, v := range guageStorage.GetAll() {
			allmtrcs = append(allmtrcs, fmt.Sprintf("name: %v value: %v", n, v))
		}
		for n, v := range counterStorage.GetAll() {
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
			log.Fatal("Error Parsing template: ", err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err1 := tmpl.Execute(w, allmtrcs)
		if err1 != nil {
			log.Fatal("Error executing template: ", err1)
		}
	}
}
