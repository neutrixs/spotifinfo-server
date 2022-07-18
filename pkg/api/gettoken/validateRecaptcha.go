package gettoken

import (
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

// returns success, errorCodes
func reCAPTCHAIsValid(reCAPTCHAToken string) (bool, []string) {
	reCAPTCHASecret, err := env.Get("RECAPTCHA_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	data := url.Values{}
	data.Set("secret", reCAPTCHASecret)
	data.Set("response", reCAPTCHAToken)

	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	parsedBody := reCAPTCHAVerifyResponse{}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		log.Println(err)
	}

	if !parsedBody.Success {
		return false, parsedBody.ErrorCodes
	}

	if parsedBody.Score < 0.5 {
		return false, []string{"low-score"}
	}

	return true, []string{}
}

// https://developers.google.com/recaptcha/docs/v3#site_verify_response
type reCAPTCHAVerifyResponse struct {
	Success 	bool		`json:"success"`
	Score 		float64 	`json:"score"`
	Action 		string 		`json:"action"`
	ChallengeTS string 		`json:"challenge_ts"`
	Hostname 	string 		`json:"hostname"`
	ErrorCodes	[]string 	`json:"error-codes"`
}