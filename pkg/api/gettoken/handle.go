package gettoken

import (
	"net/http"

	"github.com/neutrixs/spotifinfo-server/pkg/db/token"
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
}

func getCookie(cookies []*http.Cookie, name string) string {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}

	return ""
}