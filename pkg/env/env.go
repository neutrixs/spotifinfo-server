package env

import (
	"errors"
	"os"
	"strings"
)

var NotFound = errors.New("Environment variable not found")

func Get(env string) (string, error) {
	variables := os.Environ()

	for _, variable := range variables {
		split := strings.Split(variable, "=")
		key := split[0]
		val := split[1]

		if key == env {
			return val, nil
		}
	}

	return "", NotFound
}