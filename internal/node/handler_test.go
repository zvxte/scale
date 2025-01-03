package node

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zvxte/scale/pkg/monitor"
)

func TestGetStats(t *testing.T) {
	cpu := monitor.NewMock()
	mem := monitor.NewMock()
	logger := log.Default()

	req := httptest.NewRequest(
		http.MethodGet, "/stats", nil,
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
			handler := getStats(cpu, mem, logger)
			handler(recorder, req)

			resp := recorder.Result()
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf(
					"getStats() StatusCode=%d, Expected=%d",
					resp.StatusCode, test.expectedStatusCode,
				)
			}

			stats := stats{
				Cpu: monitor.CPUMaxUsage + 1, Mem: monitor.MemMaxUsage + 1,
			}

			var rawBody bytes.Buffer
			tee := io.TeeReader(resp.Body, &rawBody)

			decoder := json.NewDecoder(tee)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&stats); err != nil {
				t.Errorf(
					"getStats() Error=%v, Body=%s",
					err, rawBody.String(),
				)
			}

			if stats.Cpu > monitor.CPUMaxUsage {
				t.Errorf(
					"getStats() Cpu=%d, MaxCpu=%d, Body=%s",
					stats.Cpu, monitor.CPUMaxUsage, rawBody.String(),
				)
			}

			if stats.Mem > monitor.MemMaxUsage {
				t.Errorf(
					"getStats() Mem=%d, MaxMem=%d, Body=%s",
					stats.Mem, monitor.MemMaxUsage, rawBody.String(),
				)
			}
		})
	}
}
