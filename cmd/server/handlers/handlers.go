package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"yametrics/cmd/server/storage"
)

func HandleGuage(storage storage.GuageStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		if v, err := strconv.ParseFloat(path[1], 64); err == nil {
			storage.Update(path[0], v)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}
func HandleCounter(storage storage.CounterStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		if v, error := strconv.ParseInt(path[1], 10, 64); error == nil {
			storage.Update(path[0], v)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}
