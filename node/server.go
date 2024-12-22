package node

import (
	"crypto/tls"
	"log"
	"net/http"
)

func NewServer(tlsConfig *tls.Config, logger *log.Logger) *http.Server {
	mux := http.NewServeMux()

	statsMux := newStatsMux(logger)
	mux.Handle("/stats/", http.StripPrefix("/stats", statsMux))

	return &http.Server{
		Addr:      ":4000",
		Handler:   mux,
		TLSConfig: tlsConfig,
		ErrorLog:  logger,
	}
}
