package monitor

import (
	"math/rand"
)

const mockMaxUsage = 100

type MockMonitor struct{}

func NewMockMonitor() *MockMonitor {
	return &MockMonitor{}
}

func (m *MockMonitor) Start() {}

func (m *MockMonitor) Stop() {}

func (m *MockMonitor) Usage() uint8 {
	return uint8(rand.Intn(mockMaxUsage + 1))
}
