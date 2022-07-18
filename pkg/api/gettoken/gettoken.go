package gettoken

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/neutrixs/spotifinfo-server/pkg/env"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func getToken(refreshToken string) spotifyGettokenResponse {
	clientID, err := env.Get("CLIENT_ID")
	if err != nil {
		log.Fatal(err)
	}

	clientSecret, err := env.Get("CLIENT_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret)))

	httpClient := http.Client{}
	
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	parsedBody := spotifyGettokenResponse{}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Println(err)
	}

	return parsedBody
}

type spotifyGettokenResponse struct {
	AccessToken 	string 	`json:"access_token"`
	TokenType 		string 	`json:"token_type"`
	Scope 			string 	`json:"scope"`
	ExpiresIn 		int 	`json:"expires_in"`
}