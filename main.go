package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/neutrixs/spotifinfo-server/pkg/env"
	"github.com/neutrixs/spotifinfo-server/pkg/spa"
)

func main() {
	staticDirPath, err := env.Get("STATIC_DIR_PATH")
	if err != nil {
		log.Fatal("No STATIC_DIR_PATH variable found")
	}

	PORT, err := env.Get("PORT")
	if err != nil {
		PORT = ":8080"
	}
	if !strings.HasPrefix(PORT, ":") {
		PORT = ":" + PORT
	}

	r := mux.NewRouter()

	spa := spa.Handler {
		StaticPath: staticDirPath,
		IndexPath: "index.html",
	}

	r.PathPrefix("/").Handler(spa)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(PORT, nil))
}