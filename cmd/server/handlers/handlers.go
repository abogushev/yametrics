package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"yametrics/cmd/server/storage"
)

func HandleGuage(storage storage.GuageStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.ReplaceAll(r.URL.Path, "/update/gauge/", ""), "/")
		if len(path) != 2 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			if v, err := strconv.ParseFloat(path[1], 64); err == nil {
				storage.Update(path[0], v)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}
func HandleCounter(storage storage.CounterStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.ReplaceAll(r.URL.Path, "/update/counter/", ""), "/")
		if len(path) != 2 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			if v, error := strconv.ParseInt(path[1], 10, 64); error == nil {
				storage.Update(path[0], v)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}
