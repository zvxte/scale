package node

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/scale/node/monitor"
)

func NewServer(tlsConfig *tls.Config, logger *log.Logger) *http.Server {
	cpu := monitor.NewCPU(5*time.Second, logger)
	cpu.Start()

	mem := monitor.NewMem(5*time.Second, logger)
	mem.Start()

	mux := http.NewServeMux()

	statsMux := newStatsMux(cpu, mem, logger)
	mux.Handle("/stats/", http.StripPrefix("/stats", statsMux))

	return &http.Server{
		Addr:      ":4000",
		Handler:   mux,
		TLSConfig: tlsConfig,
		ErrorLog:  logger,
	}
}
