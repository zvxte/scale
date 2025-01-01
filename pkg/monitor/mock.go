package monitor

import (
	"math/rand"
)

const MockMaxUsage = 100

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Start() {}

func (m *Mock) Stop() {}

func (m *Mock) Usage() uint8 {
	return uint8(rand.Intn(MockMaxUsage + 1))
}
