package node

import (
	"log"
	"net/http"

	"github.com/zvxte/scale/node/monitor"
)

func newStatsMux(
	cpu monitor.Monitor,
	mem monitor.Monitor,
	logger *log.Logger,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /summary", getStatsSummary(
		cpu,
		mem,
		logger,
	))
	return mux
}
