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

const (
	CPUMinUsage    = 0
	CPUMaxUsage    = 100
	CPUMinInterval = 1 * time.Second
)

type CPU struct {
	// CPU usage percentage
	usage uint8

	// How often should the usage be updated
	interval time.Duration

	logger *log.Logger
	mu     sync.RWMutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewCPU(interval time.Duration, logger *log.Logger) *CPU {
	if interval < CPUMinInterval {
		interval = CPUMinInterval
	}

	return &CPU{
		interval: interval,
		logger:   logger,
	}
}

func (m *CPU) Start() {
	if m.ctx != nil {
		return
	}

	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.wg.Add(1)

	go func() {
		defer m.wg.Done()

		// Initialize stats to make the first usage update accurate.
		stats, err := readCpuStats()
		if err != nil {
			m.logger.Println(
				fmt.Errorf("failed to read CPU stats: %w", err),
			)
			m.Stop()
			return
		}
		time.Sleep(m.interval)

		for {
			select {
			case <-m.ctx.Done():
				return

			default:
				currentStats, err := readCpuStats()
				if err != nil {
					m.logger.Println(
						fmt.Errorf("failed to read CPU stats: %w", err),
					)
					m.Stop()
					return
				}

				totalDiff := currentStats.total - stats.total
				idleDiff := currentStats.idle - stats.idle

				// Avoid division by zero
				if totalDiff == 0 {
					totalDiff = 1
				}

				usage := 100 - ((idleDiff * 100) / totalDiff)
				if usage > CPUMaxUsage {
					usage = CPUMaxUsage
				}
				m.setUsage(uint8(usage))

				stats = currentStats

				time.Sleep(m.interval)
			}
		}
	}()
}

func (m *CPU) Stop() {
	if m.cancel != nil {
		m.cancel()
		m.wg.Wait()

		m.ctx = nil
		m.cancel = nil

		m.setUsage(CPUMinUsage)
	}
}

func (m *CPU) Usage() uint8 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.usage
}

func (m *CPU) setUsage(usage uint8) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.usage = usage
}

const cpuStatsFile = "/proc/stat"

var errInvalidCpuStatsFile = errors.New("unexpected /proc/stat file format")

type cpuStats struct {
	total uint64
	idle  uint64
}

func readCpuStats() (cpuStats, error) {
	file, err := os.Open(cpuStatsFile)
	if err != nil {
		return cpuStats{}, err
	}
	defer file.Close()

	line := ""

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return cpuStats{}, err
	}

	if !strings.HasPrefix(line, "cpu") {
		return cpuStats{}, errInvalidCpuStatsFile
	}

	// "cpu  4287477 2 657908 43117172 1758 128015 49404 0 0 0\n"
	parts := strings.Fields(line)
	if len(parts) < 11 {
		return cpuStats{}, errInvalidCpuStatsFile
	}

	var total, idle uint64
	for i, part := range parts[1:] {
		value, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return cpuStats{}, err
		}

		if i == 3 {
			idle = value
		}
		total += value
	}
	return cpuStats{total: total, idle: idle}, nil
}
