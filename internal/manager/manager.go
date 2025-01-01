package manager

type Manager struct {
	// ...
}

type Cluster struct {
	ID       uint
	Name     string
	nodes    []*Node
	clusters []*Cluster
}

type Node struct {
	ID   uint
	Name string
	Addr string
}
