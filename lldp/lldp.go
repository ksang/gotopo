/*
Package lldp provide lldp data structre and endpoint data collection implementations
*/
package lldp

/*
Snapshot represents the data snapshot gatherred from a LLDP endpoint
*/
type Snapshot struct {
	Local LLDPLocalSystemData
}

type LLDPLocalSystemData struct {
	ChassisId   string
	Name        string
	Description string
	PortTable   []LLDPPortTableEntry
}

type LLDPPortTableEntry struct {
	Number      int
	Id          string
	Description string
}
