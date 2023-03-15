package login

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/neutrixs/spotifinfo-server/pkg/env"
	"github.com/neutrixs/spotifinfo-server/pkg/querystring"
	"golang.org/x/exp/slices"
)

func Handle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cookies := r.Cookies()

	stateCookieIndex := slices.IndexFunc(cookies, func(c *http.Cookie) bool {return c.Name == "state"})
	userLoggedIn := stateCookieIndex != -1
	if userLoggedIn {
		http.Redirect(w, r, "/", 302)
		return
	}

	newState := generateState(21)
	scope, err := env.Get("API_SCOPE")
	if err != nil {
		log.Fatal(err)
	}

	db.Query("INSERT INTO login (state, scopes) VALUES (?, ?)", newState, scope)

	client_id, err := env.Get("CLIENT_ID")
	if err != nil {
		log.Fatal(err)
	}

	redirect_uri, err := env.Get("REDIRECT_URI")
	if err != nil {
		log.Fatal(err)
	}

	params := querystring.Encode(map[string]string {
		"state": newState,
		"response_type": "code",
		"client_id": client_id,
		"redirect_uri": redirect_uri,
		"scope": scope,
		"show_dialog": "true",
	})

	http.Redirect(w, r, "https://accounts.spotify.com/authorize?"+params, 302)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func generateState(length int) string {
	const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789._-"
	var res string

	for i := 0; i < length; i++ {
		randLetter := possible[rand.Intn(len(possible))]
		res += string(randLetter)
	}

	return res
}