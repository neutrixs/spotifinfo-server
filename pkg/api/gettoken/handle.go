package gettoken

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/neutrixs/spotifinfo-server/pkg/querystring"
)

func Handle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorCode := http.StatusMethodNotAllowed
		w.WriteHeader(errorCode)
		w.Write([]byte(http.StatusText(errorCode)))
		return
	}

	stateCookie := getCookie(r.Cookies(), "state")
	if stateCookie == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing state cookie"))
		return
	}

	var refreshToken string

	sessionsRows, err := db.Query("SELECT refresh_token FROM sessions WHERE state=?", stateCookie)
	if err != nil {
		log.Println(err)
		return
	}
	next := sessionsRows.Next()

	//stateData, ok := token.InitToken.Get(stateCookie)
	if !next {
		response := failedResponseType{
			Success: false,
			ErrorCodes: []string{"state-not-found"},
			Relogback: true,
		}

		json, err := json.Marshal(response)
		if err != nil {
			errorCode := http.StatusInternalServerError
			w.WriteHeader(errorCode)
			w.Write([]byte(http.StatusText(errorCode)))
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(json)
		w.Header().Set("Content-Type", "application/json")
		return
	}
	sessionsRows.Scan(&refreshToken)

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content-Type not allowed"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorCode := http.StatusInternalServerError
		w.WriteHeader(errorCode)
		w.Write([]byte(http.StatusText(errorCode)))
		return
	}

	bodyParams := querystring.Decode(string(body))
	reCAPTCHAToken := bodyParams["reCAPTCHAToken"]
	if reCAPTCHAToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing reCAPTCHAToken body parameter"))
		return
	}

	valid, errorCodes := reCAPTCHAIsValid(reCAPTCHAToken)
	if !valid {
		response := failedResponseType{
			Success: false,
			ErrorCodes: errorCodes,
		}

		json, err := json.Marshal(response)
		if err != nil {
			errorCode := http.StatusInternalServerError
			w.WriteHeader(errorCode)
			w.Write([]byte(http.StatusText(errorCode)))
			return
		}

		w.WriteHeader(http.StatusForbidden)
		w.Write(json)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	tokenData := getToken(refreshToken)
	if tokenData.AccessToken == "" {
		response := failedResponseType{
			Success: false,
			Relogback: true,
		}

		json, err := json.Marshal(response)
		if err != nil {
			errorCode := http.StatusInternalServerError
			w.WriteHeader(errorCode)
			w.Write([]byte(http.StatusText(errorCode)))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	validUntil := time.Now().UnixMilli() + int64(tokenData.ExpiresIn) * 1000

	response := successResponseType{
		Success: true,
		Data: responseDataType{
			Token: tokenData.TokenType + " " + tokenData.AccessToken,
			ValidUntil: int(validUntil),
		},
	}

	json, err := json.Marshal(response)
	if err != nil {
		errorCode := http.StatusInternalServerError
		w.WriteHeader(errorCode)
		w.Write([]byte(http.StatusText(errorCode)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
	w.Header().Set("Content-Type", "application/json")
}

func getCookie(cookies []*http.Cookie, name string) string {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}

	return ""
}