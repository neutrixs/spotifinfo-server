package querystring

import (
	"net/url"
	"strings"
)

func Encode(data map[string]string) string {
	var res []string

	for key := range data {
		val := data[key]

		keyVal := url.QueryEscape(key) + "=" + url.QueryEscape(val)
		res = append(res, keyVal)
	}

	return strings.Join(res, "&")
}

func Decode(data string) map[string]string {
	res := make(map[string]string)

	splittedKeyVals := strings.Split(data, "&")

	for _, keyVal := range splittedKeyVals {
		key, _ := url.QueryUnescape(strings.Split(keyVal, "=")[0])
		val, _ := url.QueryUnescape(strings.Split(keyVal, "=")[1])

		res[key] = val
	}

	return res
}