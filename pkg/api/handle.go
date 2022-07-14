package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neutrixs/spotifinfo-server/pkg/api/login"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	endpoint := mux.Vars(r)["endpoint"]

	switch endpoint {
	case "login":
		login.Handle(w, r)
	}
}