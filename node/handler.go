package node

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zvxte/scale/node/monitor"
)

type statsSummary struct {
	Cpu uint8 `json:"cpu"`
	Mem uint8 `json:"mem"`
}

func getStatsSummary(
	cpu monitor.Monitor,
	mem monitor.Monitor,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(statsSummary{
			Cpu: cpu.Usage(),
			Mem: mem.Usage(),
		}); err != nil {
			logger.Println(err)
		}
	}
}
