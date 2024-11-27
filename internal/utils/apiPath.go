package utils

import (
	"fmt"
	"os"
)

func ApiRoute(route string) string {
	apiV, ok := os.LookupEnv("API_VERSION")
	if !ok {
		apiV = "/api/v1/"
	}

	return fmt.Sprintf("%s%s", apiV, route)
}
