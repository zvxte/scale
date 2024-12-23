package monitor

import (
	"log"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	logger := log.Default()

	tests := []struct {
		name     string
		monitor  *CPUMonitor
		duration time.Duration
	}{
		{
			"Valid: No monitoring duration",
			NewCPUMonitor(1*time.Second, logger),
			0 * time.Second,
		},
		{
			"Valid: Short monitoring duration",
			NewCPUMonitor(1*time.Second, logger),
			2 * time.Second,
		},
		{
			"Valid: long monitoring duration",
			NewCPUMonitor(1*time.Second, logger),
			10 * time.Second,
		},
		{
			"Valid: short interval",
			NewCPUMonitor(250*time.Millisecond, logger),
			2 * time.Second,
		},
		{
			"Valid: long interval",
			NewCPUMonitor(4*time.Second, logger),
			10 * time.Second,
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
					"Start() Usage=%d, MaxUsage=%d",
					usage, cpuMaxUsage,
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
		{
			"Valid",
			false,
		},
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
