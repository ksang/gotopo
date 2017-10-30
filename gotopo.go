/*
Package gotopo provide L2 topology data structure based on LLDP
*/
package gotopo

import "sync"

/*
Topology represents a network
*/
type Topology struct {
	Nodes []Node
	Graph [][]Link
	mu    sync.RWMutex
}

// AddNode add a node to topology
func (t *Topology) AddNode(n Node) error {
	return nil
}

// RemoveNode remove a node from topology
func (t *Topology) RemoveNode(n Node) error {
	return nil
}

// AddLink add a link between two nodes to topology
func (t *Topology) AddLink(idx uint, idy uint, l Link) error {
	return nil
}

// RemoveLink remove a link between two nodes from topology
func (t *Topology) RemoveLink(idx uint, idy uint) error {
	return nil
}
