package main

import (
	"fmt"
	"log"

	"github.com/zvxte/scale/node"
)

func main() {
	logger := log.Default()

	tlsConfig, err := node.LoadTLSConfig()
	if err != nil {
		logger.Println(fmt.Errorf("failed to load TLS config: %w", err))
		return
	}

	server := node.NewServer(tlsConfig, logger)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		logger.Println(fmt.Errorf("failed to run server: %w", err))
		return
	}
}
