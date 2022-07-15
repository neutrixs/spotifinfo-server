package callback

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/neutrixs/spotifinfo-server/pkg/db/state"
	"github.com/neutrixs/spotifinfo-server/pkg/db/token"
	"github.com/neutrixs/spotifinfo-server/pkg/env"
)

const closeWindow = "<script>window.close();</script>"

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if queries.Get("error") != "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(closeWindow))
		return
	}

	if queries.Get("code") == "" || queries.Get("state") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(closeWindow))
		return
	}

	if state.InitStates.Get(queries.Get("state")) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State not found"))
		return
	}

	clientID, err := env.Get("CLIENT_ID")
	if err != nil {
		log.Fatal(err)
	}

	clientSecret, err := env.Get("CLIENT_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	redirectURI, err := env.Get("REDIRECT_URI")
	if err != nil {
		log.Fatal(err)
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", queries.Get("code"))
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("authorization", authorization)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{}

	res, err := httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	parsedResponse := spotifyAPIResponse{}

	err = json.Unmarshal(bodyData, &parsedResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !scopesMatch(state.InitStates.Get(queries.Get("state")), parsedResponse.Scope) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Mismatched scope"))
		return
	}

	token.InitToken.Add(queries.Get("state"), &token.EachState{
		RefreshToken: parsedResponse.RefreshToken,
	})

	stateCookie := http.Cookie{
		Name: "state",
		Value: queries.Get("state"),
		Expires: time.Now().AddDate(1000, 0, 0),
		Path: "/",
	}

	http.SetCookie(w, &stateCookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(closeWindow))
}

func scopesMatch(scopes1, scopes2 string) bool {
	sortedScopes1 := strings.Split(scopes1, " ")
	sort.Strings(sortedScopes1)

	sortedScopes2 := strings.Split(scopes2, " ")
	sort.Strings(sortedScopes2)

	return strings.Join(sortedScopes1, " ") == strings.Join(sortedScopes2, " ")
}

type spotifyAPIResponse struct {
	AccessToken 	string 	`json:"access_token"`
	TokenType 		string 	`json:"token_type"`
	Scope 			string 	`json:"Scope"`
	Expires 		int		`json:"expires_in"`
	RefreshToken 	string	`json:"refresh_token"`
}