package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"joefang.com/handlers"
)

func main() {
	router := mux.NewRouter()
	router.Path("/logs/files/{file}").HandlerFunc(handlers.GetLogFile).Methods(http.MethodGet)
	router.Path("/logs/files/{file}/lastevents/{lastEvents:[0-9]+}").HandlerFunc(handlers.GetLogEvents).Methods(http.MethodGet)
	router.Path("/logs/files/{file}/lastevents/{lastEvents:[0-9]+}").Queries("filter", "{filter}").HandlerFunc(handlers.GetLogEvents).Methods(http.MethodGet)
	fmt.Println("server listening on port: 8000")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Panic(err)
	}
}
