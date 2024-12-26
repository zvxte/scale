package monitor

import (
	"io"
	"log"
	"strings"
	"testing"
	"time"
)

func TestMem(t *testing.T) {
	logger := log.Default()

	tests := []struct {
		name     string
		mem      *Mem
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

			test.mem.Start()
			time.Sleep(test.duration)

			usage := test.mem.Usage()
			if usage > 100 {
				t.Errorf(
					"Mem{} Usage=%d, MaxUsage=%d",
					usage, MemMaxUsage,
				)
			}

			test.mem.Stop()
			if test.mem.Usage() != 0 {
				t.Errorf(
					"Mem{} Usage=%d, ExpectedUsage=%d",
					usage, test.mem.usage,
				)
			}
		})
	}
}

func TestLoadMemStats(t *testing.T) {
	tests := []struct {
		name      string
		r         io.Reader
		shouldErr bool
	}{
		{
			"Valid",
			strings.NewReader(
				"MemTotal:       16000000 kB\n" +
					"MemFree:         8000000 kB\n" +
					"MemAvailable:   80000000 kB\n",
			),
			false,
		},
		{
			"Valid: zero values",
			strings.NewReader(
				"MemTotal:              0 kB\n" +
					"MemFree:               0 kB\n" +
					"MemAvailable:          0 kB\n",
			),
			false,
		},
		{
			"Valid: empty second line",
			strings.NewReader(
				"MemTotal:       16000000 kB\n" +
					"\n" +
					"MemAvailable:   80000000 kB\n",
			),
			false,
		},
		{
			"Invalid: no second line",
			strings.NewReader(
				"MemTotal:       16000000 kB\n" +
					"MemAvailable:   80000000 kB\n",
			),
			true,
		},
		{
			"Invalid: total value",
			strings.NewReader(
				"MemTotal:                kB\n" +
					"MemFree:         8000000 kB\n" +
					"MemAvailable:   80000000 kB\n",
			),
			true,
		},
		{
			"Invalid: available value",
			strings.NewReader(
				"MemTotal:       16000000 kB\n" +
					"MemFree:         8000000 kB\n" +
					"MemAvailable:       some kB\n",
			),
			true,
		},
		{
			"Invalid: no newline",
			strings.NewReader(
				"MemTotal:       16000000 kB" +
					"MemFree:         8000000 kB\n" +
					"MemAvailable:   80000000 kB\n",
			),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := loadMemStats(test.r)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"readMemStats() Error=%v, ShouldError=%v",
					err, test.shouldErr,
				)
			}
		})
	}
}
