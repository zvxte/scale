package node

import (
	"log"
	"net/http"
	"time"

	"github.com/zvxte/scale/pkg/monitor"
)

func newMux(
	cpuInterval time.Duration,
	memInterval time.Duration,
	logger *log.Logger,
) *http.ServeMux {
	cpu := monitor.NewCPU(cpuInterval, logger)
	cpu.Start()

	mem := monitor.NewMem(memInterval, logger)
	mem.Start()

	mux := http.NewServeMux()
	mux.Handle("GET /stats", getStats(cpu, mem, logger))
	return mux
}
