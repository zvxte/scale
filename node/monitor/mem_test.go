package monitor

import (
	"log"
	"testing"
	"time"
)

func TestMemMonitor(t *testing.T) {
	logger := log.Default()

	tests := []struct {
		name     string
		monitor  *MemMonitor
		duration time.Duration
	}{
		{
			"Valid: No monitoring duration",
			NewMemMonitor(1*time.Second, logger),
			0 * time.Second,
		},
		{
			"Valid: Short monitoring duration",
			NewMemMonitor(1*time.Second, logger),
			2 * time.Second,
		},
		{
			"Valid: longer monitoring duration",
			NewMemMonitor(1*time.Second, logger),
			4 * time.Second,
		},
		{
			"Valid: short interval",
			NewMemMonitor(250*time.Millisecond, logger),
			2 * time.Second,
		},
		{
			"Valid: longer interval",
			NewMemMonitor(2*time.Second, logger),
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
			if usage > 100 {
				t.Errorf(
					"MemMonitor{} Usage=%d, MaxUsage=%d",
					usage, cpuMaxUsage,
				)
			}
		})
	}
}

func TestReadMemStats(t *testing.T) {
	tests := []struct {
		name      string
		shouldErr bool
	}{
		{"Valid", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := readMemStats()
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"readMemStats() Error=%v, ShouldError=%v",
					err, test.shouldErr,
				)
			}
		})
	}
}
