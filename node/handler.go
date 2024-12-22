package node

import (
	"fmt"
	"log"
	"net/http"
)

func getIndex(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s: %s\n", r.Method, r.URL)
		fmt.Fprintf(w, "%s: %s\n", r.Method, r.URL)
	}
}
