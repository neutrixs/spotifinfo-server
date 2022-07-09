package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	endpoint := mux.Vars(r)["endpoint"]

	w.Write([]byte(endpoint))

	switch endpoint {
	
	}
}