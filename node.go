package gotopo

/*
Node represents a Node in network topology or vertical in a graph
*/
type Node struct {
	// Internal used id to identify the node in topology
	ID          uint8
	ChassisID   string
	Name        string
	Description string
	topo        *Topology
}

// Neighbors returns the neighbor nodes of this node in topology
func (n *Node) Neighbors() []Node {
	return nil
}
