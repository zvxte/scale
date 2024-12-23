package node

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zvxte/scale/node/monitor"
)

func TestGetStatsSummary(t *testing.T) {
	logger := log.Default()
	monitor := monitor.NewMockMonitor()
	req := httptest.NewRequest(
		http.MethodGet, "/stats/summary", nil,
	)

	const (
		maxCpu = 100
		maxMem = 100
	)

	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{"Valid", 200},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler := getStatsSummary(
				monitor, monitor, logger,
			)
			handler(recorder, req)

			resp := recorder.Result()
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf(
					"getStatsSummary() StatusCode=%d, Expected=%d",
					resp.StatusCode, test.expectedStatusCode,
				)
			}

			// Set stats to max values to later verify if they have changed
			stats := statsSummary{Cpu: 255, Mem: 255}

			var rawBody bytes.Buffer
			tee := io.TeeReader(resp.Body, &rawBody)

			decoder := json.NewDecoder(tee)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&stats); err != nil {
				t.Errorf(
					"getStatsSummary() Error=%v",
					err,
				)
			}

			if stats.Cpu == 255 || stats.Mem == 255 {
				t.Errorf(
					"getStatsSummary() Body=%s",
					rawBody.String(),
				)
			}

			if stats.Cpu > maxCpu {
				t.Errorf(
					"getStatsSummary() Cpu=%d, MaxCpu=%d",
					stats.Cpu, maxCpu,
				)
			}

			if stats.Mem > maxMem {
				t.Errorf(
					"getStatsSummary() Mem=%d, MaxMem=%d",
					stats.Mem, maxMem,
				)
			}
		})
	}
}
