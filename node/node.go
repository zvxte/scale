package node

import (
	"fmt"
	"log"
	"net/http"
)

func Run(logger *log.Logger, config *Config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s: %s\n", r.Method, r.URL)
	})

	server := &http.Server{
		Handler:   mux,
		ErrorLog:  logger,
		TLSConfig: config.TLSConfig,
	}
	return server.ListenAndServeTLS("", "")
}
