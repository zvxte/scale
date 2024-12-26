package monitor

import (
	"log"
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
			"Valid: No monitoring duration",
			NewCPU(1*time.Second, logger),
			0 * time.Second,
		},
		{
			"Valid: Short monitoring duration",
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
			defer test.monitor.Stop()

			time.Sleep(test.duration)

			usage := test.monitor.Usage()
			if usage > CPUMaxUsage {
				t.Errorf(
					"CPU{} Usage=%d, MaxUsage=%d",
					usage, CPUMaxUsage,
				)
			}
		})
	}
}

func TestReadCpuStats(t *testing.T) {
	tests := []struct {
		name      string
		shouldErr bool
	}{
		{"Valid", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := readCpuStats()
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"readCpuStats() Error=%v, ShouldError=%v",
					err, test.shouldErr,
				)
			}
		})
	}
}
