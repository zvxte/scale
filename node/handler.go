package node

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zvxte/scale/node/monitor"
)

type stats struct {
	Cpu uint8 `json:"cpu"`
	Mem uint8 `json:"mem"`
}

func getStats(
	cpu monitor.Monitor,
	mem monitor.Monitor,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(stats{
			Cpu: cpu.Usage(),
			Mem: mem.Usage(),
		}); err != nil {
			logger.Println(
				fmt.Errorf("failed to encode stats: %w", err),
			)
		}
	}
}
