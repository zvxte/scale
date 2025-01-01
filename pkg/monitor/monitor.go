package monitor

type Monitor interface {
	Start()
	Stop()
	Usage() uint8
}
