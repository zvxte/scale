package monitor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const memMaxUsage = 100

type MemMonitor struct {
	// Mem usage percentage
	usage uint8

	// How often should the usage be updated
	interval time.Duration

	logger *log.Logger
	mu     sync.RWMutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMemMonitor(interval time.Duration, logger *log.Logger) *MemMonitor {
	return &MemMonitor{
		interval: interval,
		logger:   logger,
	}
}

func (m *MemMonitor) Start() {
	if m.ctx != nil {
		return
	}

	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.wg.Add(1)

	go func() {
		defer m.wg.Done()

		for {
			select {
			case <-m.ctx.Done():
				return

			default:
				stats, err := readMemStats()
				if err != nil {
					m.logger.Println(
						fmt.Errorf("failed to read mem stats: %w", err),
					)
					m.Stop()
					return
				}

				// Avoid division by zero
				if stats.total == 0 {
					stats.total = 1
				}

				usage := 100.0 *
					(1.0 - (float32(stats.available) / float32(stats.total)))
				if usage > memMaxUsage {
					usage = memMaxUsage
				}
				m.setUsage(uint8(usage))

				time.Sleep(m.interval)
			}
		}
	}()
}

func (m *MemMonitor) Stop() {
	if m.cancel != nil {
		m.cancel()
		m.wg.Wait()

		m.ctx = nil
		m.cancel = nil

		m.setUsage(0)
	}
}

func (m *MemMonitor) Usage() uint8 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.usage
}

func (m *MemMonitor) setUsage(usage uint8) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.usage = usage
}

const memStatsFile = "/proc/meminfo"

var errInvalidMemStatsFile = errors.New("unexpected /proc/meminfo file format")

type memStats struct {
	total     uint64
	available uint64
}

func readMemStats() (memStats, error) {
	file, err := os.Open(memStatsFile)
	if err != nil {
		return memStats{}, err
	}
	defer file.Close()

	// MemTotal:       16193664 kB\n
	// MemFree:        11712900 kB\n
	// MemAvailable:   13580972 kB\n
	lines := make([]string, 3)

	scanner := bufio.NewScanner(file)
	for i := 0; i < 3 && scanner.Scan(); i++ {
		lines[i] = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return memStats{}, err
	}

	totalLine := lines[0]
	if !strings.HasPrefix(totalLine, "MemTotal") {
		return memStats{}, errInvalidMemStatsFile
	}
	totalParts := strings.Fields(totalLine)
	if len(totalParts) < 3 {
		return memStats{}, errInvalidMemStatsFile
	}
	total, err := strconv.ParseUint(totalParts[1], 10, 64)
	if err != nil {
		return memStats{}, err
	}

	availableLine := lines[2]
	if !strings.HasPrefix(availableLine, "MemAvailable") {
		return memStats{}, errInvalidMemStatsFile
	}
	availableParts := strings.Fields(availableLine)
	if len(availableParts) < 3 {
		return memStats{}, errInvalidMemStatsFile
	}
	available, err := strconv.ParseUint(availableParts[1], 10, 64)
	if err != nil {
		return memStats{}, err
	}

	return memStats{total: total, available: available}, nil
}
