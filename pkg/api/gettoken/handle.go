package gettoken

import (
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorCode := http.StatusMethodNotAllowed
		w.WriteHeader(errorCode)
		w.Write([]byte(http.StatusText(errorCode)))
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