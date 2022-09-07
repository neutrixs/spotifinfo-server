package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neutrixs/spotifinfo-server/pkg/api"
	"github.com/neutrixs/spotifinfo-server/pkg/env"
)

func main() {
	staticDirPath, err := env.Get("STATIC_DIR_PATH")
	if err != nil {
		log.Fatal("No STATIC_DIR_PATH variable found")
	}

	indexHtmlPath, err := env.Get("INDEX_HTML_PATH")
	if err != nil {
		log.Fatal("No INDEX_HTML_PATH variable found")
	}

	PORT, err := env.Get("PORT")
	if err != nil {
		PORT = ":8080"
	}
	
	r := mux.NewRouter()
	r.Path("/api/{endpoint}").Methods("GET", "POST").HandlerFunc(api.Handle)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDirPath))))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexHtmlPath)
	})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(PORT, nil))
}