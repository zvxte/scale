package node

import (
	"log"
	"net/http"

	"github.com/zvxte/scale/node/monitor"
)

func newStatsMux(
	cpuMonitor monitor.Monitor,
	memMonitor monitor.Monitor,
	logger *log.Logger,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /summary", getStatsSummary(
		cpuMonitor,
		memMonitor,
		logger,
	))
	return mux
}
