package gettoken

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/neutrixs/spotifinfo-server/pkg/db/token"
	"github.com/neutrixs/spotifinfo-server/pkg/querystring"
)

func Handle(w http.ResponseWriter, r *http.Request) {
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

	/*stateData*/_, ok := token.InitToken.Get(stateCookie)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("state not found"))
		return
	}

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
}

func getCookie(cookies []*http.Cookie, name string) string {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}

	return ""
}