package main

import (
	"log"
	"net/http"
	"yametrics/cmd/server/handlers"
	"yametrics/cmd/server/storage"
)

func main() {
	metricsStorage := storage.NewGuageStorage()
	countersStorage := storage.NewCounterStorage()
	guagePath := "/update/gauge/"
	counterPath := "/update/counter/"

	// маршрутизация запросов обработчику
	http.Handle(guagePath, http.StripPrefix(guagePath, handlers.HandleGuage(metricsStorage)))
	http.Handle(counterPath, http.StripPrefix(counterPath, handlers.HandleCounter(countersStorage)))
	http.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) })

	log.Println("server started successfull")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
