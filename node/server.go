package node

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/scale/node/monitor"
)

func NewServer(tlsConfig *tls.Config, logger *log.Logger) *http.Server {
	cpuMonitor := monitor.NewCPUMonitor(5*time.Second, logger)
	cpuMonitor.Start()

	memMonitor := monitor.NewMemMonitor(5*time.Second, logger)
	memMonitor.Start()

	mux := http.NewServeMux()

	statsMux := newStatsMux(cpuMonitor, memMonitor, logger)
	mux.Handle("/stats/", http.StripPrefix("/stats", statsMux))

	return &http.Server{
		Addr:      ":4000",
		Handler:   mux,
		TLSConfig: tlsConfig,
		ErrorLog:  logger,
	}
}
