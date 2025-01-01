package monitor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MemMinUsage    = 0
	MemMaxUsage    = 100
	MemMinInterval = 1 * time.Second

	memStatsFile = "/proc/meminfo"
)

type Mem struct {
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

func NewMem(interval time.Duration, logger *log.Logger) *Mem {
	if interval < MemMinInterval {
		interval = MemMinInterval
	}

	return &Mem{
		interval: interval,
		logger:   logger,
	}
}

func (m *Mem) Start() {
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
				stats, err := readMemStatsFile()
				if err != nil {
					m.logger.Println(
						fmt.Errorf("failed to read mem stats file: %w", err),
					)
					m.Stop()
					return
				}

				// Avoid division by zero
				if stats.total == 0 {
					stats.total = 1
				}

				usage := 100 - ((stats.available * 100) / stats.total)
				if usage > MemMaxUsage {
					usage = MemMaxUsage
				}
				m.setUsage(uint8(usage))

				time.Sleep(m.interval)
			}
		}
	}()
}

func (m *Mem) Stop() {
	if m.cancel != nil {
		m.cancel()
		m.wg.Wait()

		m.ctx = nil
		m.cancel = nil

		m.setUsage(MemMinUsage)
	}
}

func (m *Mem) Usage() uint8 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.usage
}

func (m *Mem) setUsage(usage uint8) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.usage = usage
}

var errInvalidMemStats = errors.New("unexpected mem stats format")

type memStats struct {
	total     uint64
	available uint64
}

func loadMemStats(r io.Reader) (memStats, error) {
	// MemTotal:       16000000 kB\n
	// MemFree:          800000 kB\n
	// MemAvailable:    8000000 kB\n
	// ...
	lines := make([]string, 3)

	scanner := bufio.NewScanner(r)
	for i := 0; i < 3 && scanner.Scan(); i++ {
		lines[i] = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return memStats{}, err
	}

	totalLine := lines[0]
	if !strings.HasPrefix(totalLine, "MemTotal") {
		return memStats{}, errInvalidMemStats
	}
	totalParts := strings.Fields(totalLine)
	if len(totalParts) < 3 {
		return memStats{}, errInvalidMemStats
	}
	total, err := strconv.ParseUint(totalParts[1], 10, 64)
	if err != nil {
		return memStats{}, err
	}

	availableLine := lines[2]
	if !strings.HasPrefix(availableLine, "MemAvailable") {
		return memStats{}, errInvalidMemStats
	}
	availableParts := strings.Fields(availableLine)
	if len(availableParts) < 3 {
		return memStats{}, errInvalidMemStats
	}
	available, err := strconv.ParseUint(availableParts[1], 10, 64)
	if err != nil {
		return memStats{}, err
	}

	return memStats{total: total, available: available}, nil
}

func readMemStatsFile() (memStats, error) {
	file, err := os.Open(memStatsFile)
	if err != nil {
		return memStats{}, err
	}
	defer file.Close()
	return loadMemStats(file)
}
