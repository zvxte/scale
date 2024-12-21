package main

import (
	"fmt"
	"log"

	"github.com/zvxte/scale/node"
)

func main() {
	logger := log.Default()

	config, err := node.LoadConfig()
	if err != nil {
		logger.Println(fmt.Errorf("failed to load config: %w", err))
		return
	}

	if err := node.Run(logger, config); err != nil {
		logger.Println(fmt.Errorf("failed to run server: %w", err))
		return
	}
}
