/*
Package lldp provide lldp data structre and endpoint data collection implementations
*/
package lldp

/*
Snapshot represents the data snapshot gatherred from a LLDP endpoint
*/
type Snapshot struct {
	Local LocalSystemData
}

// LocalSystemData is the local lldp data structure
type LocalSystemData struct {
	ChassisID   string
	Name        string
	Description string
	PortTable   []PortTableEntry
}

// PortTableEntry is the porttable entry data structure
type PortTableEntry struct {
	Number      int
	ID          string
	Description string
}
