package node

import (
	"log"
	"net/http"
)

func newStatsMux(logger *log.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /summary", getStatsSummary(logger))
	return mux
}
