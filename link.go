package gotopo

/*
Link represents a network link connects two nodes
*/
type Link struct {
	// port id of X axis node
	PortIDX string
	// port id of y axis node
	PortIDY string
	Speed   string
	// Custom properties attached to this link
	Properties interface{}
	topo       *Topology
}

// SetProperties set properties to this link
func (l *Link) SetProperties(p interface{}) error {
	return nil
}
