package node

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

func NewServer(
	addr string,
	tlsConfig *tls.Config,
	cpuInterval time.Duration,
	memInterval time.Duration,
	logger *log.Logger,
) *http.Server {
	mux := newMux(cpuInterval, memInterval, logger)
	return &http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Handler:   mux,
		ErrorLog:  logger,
	}
}
