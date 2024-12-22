package node

import (
	"crypto/tls"
	"log"
	"net/http"
)

func NewServer(tlsConfig *tls.Config, logger *log.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", getIndex(logger))

	return &http.Server{
		Handler:   mux,
		TLSConfig: tlsConfig,
		ErrorLog:  logger,
	}
}
