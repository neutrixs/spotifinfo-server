package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neutrixs/spotifinfo-server/pkg/api/callback"
	"github.com/neutrixs/spotifinfo-server/pkg/api/gettoken"
	"github.com/neutrixs/spotifinfo-server/pkg/api/login"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	endpoint := mux.Vars(r)["endpoint"]

	switch endpoint {
	case "login":
		login.Handle(w, r)
		return
	case "callback":
		callback.Handle(w, r)
		return
	case "gettoken":
		gettoken.Handle(w, r)
		return
	}

	status := http.StatusNotFound
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}