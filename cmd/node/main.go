package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zvxte/scale/internal/node"
	"github.com/zvxte/scale/pkg/mtls"
)

func main() {
	addr := "localhost:4000"
	caCertFile := "/etc/scale/ca.crt"
	keyFile := "/etc/scale/node.key"
	certFile := "/etc/scale/node.crt"
	cpuInterval := 10 * time.Second
	memInterval := 10 * time.Second

	logger := log.Default()

	tlsConfig, err := mtls.LoadServer(caCertFile, keyFile, certFile)
	if err != nil {
		logger.Println(fmt.Errorf("failed to load TLS config: %w", err))
		return
	}

	server := node.NewServer(addr, tlsConfig, cpuInterval, memInterval, logger)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		logger.Println(fmt.Errorf("failed to run server: %w", err))
		return
	}
}
