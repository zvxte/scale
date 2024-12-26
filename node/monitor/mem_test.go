package monitor

import (
	"log"
	"testing"
	"time"
)

func TestMem(t *testing.T) {
	logger := log.Default()

	tests := []struct {
		name     string
		monitor  *Mem
		duration time.Duration
	}{
		{
			"Valid: No monitoring duration",
			NewMem(1*time.Second, logger),
			0 * time.Second,
		},
		{
			"Valid: Short monitoring duration",
			NewMem(1*time.Second, logger),
			2 * time.Second,
		},
		{
			"Valid: longer monitoring duration",
			NewMem(1*time.Second, logger),
			4 * time.Second,
		},
		{
			"Valid: short interval",
			&Mem{interval: 250 * time.Millisecond, logger: logger},
			2 * time.Second,
		},
		{
			"Valid: longer interval",
			NewMem(2*time.Second, logger),
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
					"Mem{} Usage=%d, MaxUsage=%d",
					usage, MemMaxUsage,
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
