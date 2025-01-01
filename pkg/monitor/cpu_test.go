package monitor

import (
	"io"
	"log"
	"strings"
	"testing"
	"time"
)

func TestCPU(t *testing.T) {
	logger := log.Default()

	tests := []struct {
		name     string
		monitor  *CPU
		duration time.Duration
	}{
		{
			"Valid: no monitoring duration",
			NewCPU(1*time.Second, logger),
			0 * time.Second,
		},
		{
			"Valid: short monitoring duration",
			NewCPU(1*time.Second, logger),
			2 * time.Second,
		},
		{
			"Valid: longer monitoring duration",
			NewCPU(1*time.Second, logger),
			4 * time.Second,
		},
		{
			"Valid: short interval",
			&CPU{interval: 250 * time.Millisecond, logger: logger},
			2 * time.Second,
		},
		{
			"Valid: longer interval",
			NewCPU(2*time.Second, logger),
			4 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			test.monitor.Start()
			time.Sleep(test.duration)

			usage := test.monitor.Usage()
			if usage > CPUMaxUsage {
				t.Errorf(
					"CPU{} Usage=%d, MaxUsage=%d",
					usage, CPUMaxUsage,
				)
			}

			test.monitor.Stop()
			if test.monitor.Usage() != CPUMinUsage {
				t.Errorf(
					"CPU{} Usage=%d, ExpectedUsage=%d",
					usage, CPUMinUsage,
				)
			}
		})
	}
}

func TestLoadCpuStats(t *testing.T) {
	tests := []struct {
		name      string
		r         io.Reader
		shouldErr bool
	}{
		{
			"Valid",
			strings.NewReader(
				"cpu  7646126 42 1298362 83649192 2728 235341 99885 0 0 0\n",
			),
			false,
		},
		{
			"Valid: no newline",
			strings.NewReader(
				"cpu  7646126 42 1298362 83649192 2728 235341 99885 0 0 0",
			),
			false,
		},
		{
			"Invalid: missing last entry",
			strings.NewReader(
				"cpu  7646126 42 1298362 83649192 2728 235341 99885 0 0\n",
			),
			true,
		},
		{
			"Invalid: missing first entry",
			strings.NewReader(
				"7646126 42 1298362 83649192 2728 235341 99885 0 0 0\n",
			),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := loadCPUStats(test.r)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"loadCPUStats() Error=%v, ShouldError=%v",
					err, test.shouldErr,
				)
			}
		})
	}
}
