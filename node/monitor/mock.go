package monitor

import (
	"math/rand"
)

type MockMonitor struct{}

func NewMockMonitor() *MockMonitor {
	return &MockMonitor{}
}

func (m *MockMonitor) Start() {}

func (m *MockMonitor) Stop() {}

func (m *MockMonitor) Usage() uint8 {
	return uint8(rand.Intn(101))
}
