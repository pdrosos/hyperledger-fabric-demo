package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register() {
	router := mux.NewRouter()

	// default route
	router.HandleFunc("/", rootHandler).Methods("GET", "HEAD")

	// not found
	//router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.Handle("/", router)
}
