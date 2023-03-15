package api

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/neutrixs/spotifinfo-server/pkg/api/callback"
	"github.com/neutrixs/spotifinfo-server/pkg/api/gettoken"
	"github.com/neutrixs/spotifinfo-server/pkg/api/login"
	"github.com/neutrixs/spotifinfo-server/pkg/env"
)

var db *sql.DB

func Handle(w http.ResponseWriter, r *http.Request) {
	endpoint := mux.Vars(r)["endpoint"]

	switch endpoint {
	case "login":
		login.Handle(w, r, db)
		return
	case "callback":
		callback.Handle(w, r, db)
		return
	case "gettoken":
		gettoken.Handle(w, r, db)
		return
	}

	status := http.StatusNotFound
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}

func init() {
	SQL_LOGIN, err := env.Get("SQL_LOGIN")
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("mysql", SQL_LOGIN)
	if err != nil {
		log.Fatal(err)
	}
}