package node

import (
	"encoding/json"
	"log"
	"net/http"
)

type statsSummary struct {
	Cpu uint8 `json:"cpu"`
	Mem uint8 `json:"mem"`
}

func getStatsSummary(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(statsSummary{}); err != nil {
			logger.Println(err)
		}
	}
}
