package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	endpoint := mux.Vars(r)["endpoint"]

	switch endpoint {
	case "login":
		login(w,r)
	}
}